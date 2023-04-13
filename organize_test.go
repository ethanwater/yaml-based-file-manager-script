package main

import "testing"

func TestLogStatus(t *testing.T) {
	if CheckLogStatus() == false {
<<<<<<< Updated upstream
		t.Errorf("can't find log")
=======
		t.Errorf("backup.log not detected")
>>>>>>> Stashed changes
	}
}
