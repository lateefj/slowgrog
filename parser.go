package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
)

type MonitorCmd struct {
	Timestamp float64  `json:"timestamp"`
	Text      string   `json:"text"`
	Params    []string `json:"params"`
}

// Parses a line from the monitor redis command.
// Do we care about Host, Port? If so need to implement..
// TODO: This parse sucks but MVP and all....
func ParseMonitorLine(l string) (*MonitorCmd, error) {
	if l == "OK" {
		return nil, nil
	}
	m := &MonitorCmd{}
	si := strings.Index(l, "[")
	ei := strings.Index(l, "]")
	t, err := strconv.ParseFloat(l[0:si-1], 10)
	if err != nil {
		log.Printf("Could not convert timestamp from string to float: %s", t)
	}
	m.Timestamp = t
	cmdPart := strings.Split(l[ei+2:], " ")
	// Upper case for consistency the command and trim and extra "
	m.Text = strings.ToUpper(strings.Trim(cmdPart[0], "\""))
	parts := cmdPart[1:]
	m.Params = make([]string, len(parts))

	// Trim off " from params
	for i, p := range parts {
		m.Params[i] = strings.Trim(p, "\"")
	}
	return m, nil
}

type Slowlog struct {
	ID        int64
	Timestamp int64
	Duration  int64
	Command   []string
}

// Parse the slowlog
// XXX: Not working yet need to figure out how to conver slowlog to a struct
func ParesSlowlogLine(entries []interface{}, err error) ([]Slowlog, error) {
	logs := make([]Slowlog, 0)
	//Trace.Printf("Slowlog data is: %v\n", reply)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		e, ok := entry.([]interface{})
		if !ok {
			Error.Println("Bad Slowlog entry")
			continue
		}
		l := Slowlog{}
		_, err = redis.Scan(e, &l.ID, &l.Timestamp, &l.Duration, &l.Command)
		if err != nil {
			Error.Printf("Error trying to scan slowlog is %s", err)
			continue
		}
		Trace.Printf("Ok log is: %v\n", l)
		logs = append(logs, l)
	}
	return logs, nil
}
