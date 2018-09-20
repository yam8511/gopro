package main

import (
	"testing"
)

func TestGetDemo(t *testing.T) {
	if GetDemo() != "demo" {
		t.Error("Doesn't demo")
		return
	}
	t.Log("Demo OK")
}
