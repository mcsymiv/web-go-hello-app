package main

import "testing"

func Test_Run(t *testing.T) {
	err := run()
	if err != nil {
		t.Error("Failed 'run()'")
	}
}
