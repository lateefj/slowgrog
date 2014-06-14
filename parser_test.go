package main

import (
	"testing"
)

func TestMonitorLine(t *testing.T) {
	l1 := "1402710620.616109 [0 127.0.0.1:64643] \"PING\""
	cmd, err := ParseMonitorLine(l1)
	if err != nil {
		t.Errorf("Failed to parse line %s with error %s", l1, err)
	}
	if cmd.Text != "PING" {
		t.Errorf("Expected cmd.Text to be PING but it was %s", cmd.Text)
	}
	l2 := "1402728075.671283 [0 127.0.0.1:50488] \"set\" \"foo\" \"bar\""
	cmd, err = ParseMonitorLine(l2)
	if err != nil {
		t.Errorf("Failed to parse line %s with error %s", l2, err)
	}
	if cmd.Text != "SET" {
		t.Errorf("Expected cmd.Text to be SET but it was %s", cmd.Text)
		if cmd.Params[0] != "foo" {
			t.Errorf("Expected first param to be foo but was %s", cmd.Params[0])
		}
		if cmd.Params[1] != "bar" {
			t.Errorf("Expected first param to be bar but was %s", cmd.Params)
		}
	}
	l3 := "1402728079.287446 [0 127.0.0.1:50488] \"keys\" \"foo\""
	cmd, err = ParseMonitorLine(l3)
	if err != nil {
		t.Errorf("Failed to parse line %s with error %s", l3, err)
	}
	if cmd.Text != "KEYS" {
		t.Errorf("Expected cmd.Text to be KEYS but it was %s", cmd.Text)
	}
	if len(cmd.Params) != 1 {
		t.Errorf("Expected 1 param bug got %d, as %s", len(cmd.Params), cmd.Params)
	}
	if cmd.Params[0] != "foo" {
		t.Errorf("Exepcted 'foo' as first param but was %s", cmd.Params[0])
	}
}
