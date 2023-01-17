package main

import "testing"

func TestRun(t *testing.T) {
	var e string = "dev"
	_, err := run(&a)
	if err != nil {
		t.Error("Failed 'run()'")
	}
}
