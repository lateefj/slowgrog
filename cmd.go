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

func rcon() (redis.Conn, error) {
	Logger.Debugf("Connecting to redis host %s on port %d", RedisHost, RedisPort)
	c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", RedisHost, RedisPort))
	if RedisPassword != "" {
		if _, err := c.Do("AUTH", RedisPassword); err != nil {
			c.Close()
			return c, err
		}
	}
	return c, err
}

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
	MonitorCmd() chan string
	InfoCmd() (string, error)
	SlowlogCmd() []Slowlog
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

func (rc *RedisCmds) MonitorCmd() chan string {
	replies := make(chan string, MONITOR_BUFFER_SIZE)
	// In background push on the connection
	go func() {
		c := rc.conn()
		c.Send(MONITOR)
		c.Flush()
		for {
			reply, err := c.Receive()
			if err != nil {
				// Try to reconnect on error
				Logger.Errorf("Reconnecting to redis after fail %s", err)
				c.Close()
				c := rc.conn()
				c.Send(MONITOR)
				c.Flush()
				continue
			}
			r, err := redis.String(reply, err)
			if err != nil {
				Logger.Errorf("Couldn't convert reply %s", err)
				continue
			}
			replies <- r
			// process pushed message
		}
		c.Close()
	}()
	return replies
}
