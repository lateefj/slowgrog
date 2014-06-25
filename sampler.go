package main

import (
	//"fmt"

	"strings"
	"time"
)

const ()

func SampleInfo(cmds DataCmds, status *RedisStatus) (string, error) {
	Logger.Debug("Sampling INFO...")

	info, err := cmds.InfoCmd()

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

func SampleSlowlog(cmds DataCmds, status *RedisStatus) ([]Slowlog, error) {
	Logger.Debug("Sampling slowlog...")
	logs, err := cmds.SlowlogCmd()
	if err != nil {
		Logger.Errorf("DataCmds slowlog failed ewith error %s", err)
		return logs, err
	}
	status.Slowlogs = logs
	return logs, nil
}

func SampleMonitor(cmds DataCmds, stopper chan bool, status *RedisStatus) {
	replies := cmds.MonitorCmd(stopper)
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
}
