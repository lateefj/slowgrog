package main

import (
	"log"
	"strconv"
	"strings"
)

type MonitorCmd struct {
	Timestamp float64  `json:"timestamp"`
	Text      string   `json:"text"`
	Params    []string `json:"params"`
}

// Parses a line from the monitor redis command.
// Do we care about Host, Port? If so need to implement..
func ParseMonitorLine(l string) (*MonitorCmd, error) {
	Info.Println("\nParsing monitor line: " + l)
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
