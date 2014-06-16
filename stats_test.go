package main

import (
	"testing"
)

// Make sure command are incrementing counters
func TestIncCmdCount(t *testing.T) {
	IncCmdCount("KEYS")
	if cmdCounts["KEYS"] != 1 {
		t.Errorf("Expected to have a count of 1 but it was %d", cmdCounts["KEYS"])
	}
	IncCmdCount("KEYS")
	if cmdCounts["KEYS"] != 2 {
		t.Errorf("Expected to have a count of 2 but it was %d", cmdCounts["KEYS"])
	}

}

// Make sure KEYS is in the list
func TestBrokenCmd(t *testing.T) {

}

// Make sure SMEMBERS is in the list
func TestBadCmd(t *testing.T) {
}
