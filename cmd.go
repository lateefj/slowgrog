package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	MONITOR_BUFFER_SIZE = 1000
	INFO                = "INFO"
	MONITOR             = "MONITOR"
	SLOWLOG             = "SLOWLOG"
)

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

type DataCmds interface {
	SlowlogCmd() ([]Slowlog, error)
	InfoCmd() (string, error)
	MonitorCmd(chan bool) chan string
}

type RedisCmds struct {
	Pool *redis.Pool
}

func NewRedisCmds() *RedisCmds {
	server := fmt.Sprintf("%s:%d", RedisHost, RedisPort)
	p := newPool(server, RedisPassword)
	return &RedisCmds{Pool: p}
}

func (rc *RedisCmds) conn() redis.Conn {
	return rc.Pool.Get()
}

func (rc *RedisCmds) SlowlogCmd() ([]Slowlog, error) {
	c := rc.conn()
	entries, err := redis.Values(c.Do(SLOWLOG, "GET", SlowlogSize))
	if err != nil {
		Logger.Errorf("Redis SLOWLOG GET %s", err)
		return nil, err
	}
	return ParseSlowlogReply(entries, err)
}

func (rc *RedisCmds) InfoCmd() (string, error) {
	c := rc.conn()
	c.Send(INFO)
	c.Flush()
	reply, err := c.Receive()
	if err != nil {
		Logger.Errorf("Failed trying ot get INFO from redis: %s", err)
		return "", err
	}
	return redis.String(reply, err)
}

func (rc *RedisCmds) MonitorCmd(stopper chan bool) chan string {
	replies := make(chan string, MONITOR_BUFFER_SIZE)
	// In background push on the connection
	go func() {
		c := rc.conn()
		c.Send(MONITOR)
		c.Flush()
		timeoutSet := false
		for {
			// Danger performance danger but this setting has to be overwritten by config so YMMV
			// If the length is greater than 0 and the timeout is not already set then
			// this will sleep in the future
			if MonitorSampleLength > 0 && !timeoutSet {
				// Flag the timeout is set so we don't double add sleeps
				timeoutSet = true
				time.AfterFunc(time.Duration(MonitorSampleLength)*time.Microsecond, func() {
					time.Sleep(time.Duration(MonitorSampleLength) * time.Microsecond)
					// Turn timeout off so it will get reset
					timeoutSet = false
				})
			}
			select {
			case <-stopper: // Stops the monitoring!
				break
			default:
				reply, err := c.Receive()
				if err != nil {
					// Try to reconnect on error
					Logger.Errorf("Reconnecting to redis after fail %s", err)
					c.Close()
					c = rc.conn()
					c.Send(MONITOR)
					c.Flush()
					break
				}

				r, err := redis.String(reply, err)
				if err != nil {
					Logger.Errorf("Couldn't convert reply %s", err)
					break
				}
				replies <- r
			}
		}
		c.Close()
		close(replies)
	}()
	return replies
}
