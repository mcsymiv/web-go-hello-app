package main

import "testing"

func TestRun(t *testing.T) {
	var a = []string{"dev"}
	_, err := run(a)
	if err != nil {
		t.Error("Failed 'run()'")
	}
}
