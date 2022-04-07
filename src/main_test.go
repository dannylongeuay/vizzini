package main

import "testing"

func TestMain(t *testing.T) {
	if 2+2 != 4 {
		t.Error("Math Error")
	}
}
