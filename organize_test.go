package main

import "testing"

func TestLogStatus(t *testing.T) {
	if CheckLogStatus() == false {
		t.Errorf("can't find log")
	}
}
