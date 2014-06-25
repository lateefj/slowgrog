package main

import (
	//"fmt"

	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

const ()

func SampleInfo(c redis.Conn, status *RedisStatus) (string, error) {
	Logger.Debug("Sampling INFO...")
	c.Send(INFO)
	c.Flush()
	reply, err := c.Receive()
	if err != nil {
		Logger.Error(err)
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

func SampleSlowlog(c redis.Conn, status *RedisStatus) ([]Slowlog, error) {
	Logger.Debug("Sampling slowlog...")
	entries, err := redis.Values(c.Do(SLOWLOG, "GET", SlowlogSize))
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	logs, err := ParesSlowlogLine(entries, err)
	status.Slowlogs = logs
	return logs, err
}

func SampleMonitor(status *RedisStatus) {
	for {
		c, err := rcon()
		replies := make(chan string, 1000)
		if err != nil {
			Logger.Errorf("Connection failed sleeping and trying again")
			time.Sleep(1 * time.Second)
			continue
		}
		c.Send(MONITOR)
		c.Flush()
		// In background push on the connection
		go func() {
			for {
				reply, err := c.Receive()
				if err != nil {
					Logger.Errorf("Monitor reply error: %s", err)
					close(replies)
					return
				}
				r, err := redis.String(reply, err)
				if err != nil {
					Logger.Errorf("Couldn't convert reply %s", err)
					continue
				}
				replies <- r
				// process pushed message
			}
		}()
		replyIndex := 0
		for {
			reply, ok := <-replies
			if !ok {
				continue
			}
			cmdMon, err := ParseMonitorLine(reply)
			if err != nil {
				Logger.Errorf("Failed to parse line: %s", reply)
			}
			if cmdMon != nil {
				// Append if room else write over them
				if len(status.MonitorSample) <= replyIndex {
					status.MonitorSample = append(status.MonitorSample, cmdMon)
				} else {
					status.MonitorSample[replyIndex] = cmdMon
				}

				// Stats tracking
				status.stats.IncCmdCount(cmdMon.Text)
				// Increment index else reset to 0
				if replyIndex < CmdLimit-1 {
					replyIndex++
				} else {
					// Danger performance danger but this setting has to be overwritten by config so YMMV
					if MonitorSampleLength > 0 {
						time.Sleep(time.Duration(MonitorSampleLength) * time.Microsecond)
					}
					replyIndex = 0
				}
			}
		}
		close(replies)
	}
}
