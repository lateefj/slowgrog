package main

import (
	"strings"
	"testing"
	"time"
)

func TestMonitorCmd(t *testing.T) {
	rc := NewRedisCmds()
	replies := rc.MonitorCmd()
	timeout := time.AfterFunc(1*time.Second, func() {
		t.Fatalf("Failed timout execed reply")
	})
	rc.InfoCmd()
	for r := range replies {
		f := strings.Index(r, "INFO")
		if f > -1 {
			break
		}
	}
	timeout.Stop()
}

func TestInfoCmd(t *testing.T) {
	rc := NewRedisCmds()
	s, err := rc.InfoCmd()
	if err != nil {
		t.Fatalf("Error running InfoCmd %s", err)
	}
	if strings.Index(s, "redis_version") < 0 {
		t.Fatalf("Expected to find 'redis_version' in the InfoCmd")
	}
}
