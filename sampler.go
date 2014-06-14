package main

import (
	//"fmt"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	INFO    = "INFO"
	MONITOR = "MONITOR"
	SLOWLOG = "SLOWLOG"
)

func SampleInfo(c redis.Conn, status *RedisStatus) (string, error) {
	Trace.Println("Sampling slowlog...")
	c.Send(INFO)
	c.Flush()
	reply, err := c.Receive()
	if err != nil {
		Error.Println(err)
	}
	info, err := redis.String(reply, err)

	lines := strings.Split(info, "\n")
	for _, l := range lines {
		if strings.IndexAny(l, ":") > -1 {
			l = strings.Trim(l, "\n")
			l = strings.Trim(l, "\r")
			kv := strings.Split(l, ":")
			if len(kv) == 2 {
				status.Info[kv[0]] = kv[1]
			}
		}
	}
	return info, err
}
func SampleSlowlog(c redis.Conn, status *RedisStatus) (string, error) {
	Trace.Println("Sampling slowlog...")
	c.Send(SLOWLOG, "get", "10")
	c.Flush()
	reply, err := c.Receive()
	if err != nil {
		Error.Println(err)
	}
	return redis.String(reply, err)
}

func SampleMonitor(c redis.Conn, status *RedisStatus) {
	replies := make(chan string, 100)
	tout := make(chan bool, 2)
	c.Send(MONITOR)
	c.Flush()
	go func() {
		time.Sleep(time.Second * 10)
		tout <- true
		tout <- true
	}()
	go func() {
		for {
			select {
			case <-tout:
				println("Ok timeout happen now quiting..")
				return
			default:
				reply, err := c.Receive()
				if err != nil {
					Error.Println(err)

				}
				r, err := redis.String(reply, err)
				replies <- r
				// process pushed message
			}
		}
	}()
	Trace.Println("Sampling monitor")
	for {
		select {
		case reply := <-replies:
			cmdMon, err := ParseMonitorLine(reply)
			if err != nil {
				Error.Printf("Failed to parse line: %s", reply)
			}
			if cmdMon != nil {
				status.MonitorSample = append(status.MonitorSample, cmdMon)
			}
		case <-tout:
			println("timeout")
			return
		}
	}
}
