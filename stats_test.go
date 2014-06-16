package main

import (
	"testing"
)

// Make sure command are incrementing counters
func TestIncCmdCount(t *testing.T) {
	s := NewStats()
	s.IncCmdCount("KEYS")
	if s.Counts()["KEYS"] != 1 {
		t.Errorf("Expected to have a count of 1 but it was %d", s.cmdCounts["KEYS"])
	}
	s.IncCmdCount("KEYS")
	if s.Counts()["KEYS"] != 2 {
		t.Errorf("Expected to have a count of 2 but it was %d", s.Counts()["KEYS"])
	}

}

// Make sure SMEMBERS is in the list
func TestBadCmd(t *testing.T) {
	s := NewStats()
	for _, k := range BadCmdList {
		s.IncCmdCount(k)
	}
	bc := s.BadCmds()
	for _, k := range BadCmdList {
		if bc[k] != 1 {
			t.Errorf("Expected key %s in BadCmds response to be 1 but is %d", k, bc[k])
		}

	}
}
