package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"exit  ", "exit"},
		{"    help     ", "help"},
		{"    help me     ", "help"},
		{"", ""},
		{"    ", ""},
		{"HELP", "help"},
	}

	for _, c := range cases {
		got := cleanInput(c.input)
		if len(got) != len(c.expected) {
			t.Errorf("Lengths differ. Expected %d, got %d", len(c.expected), len(got))
		}
		if got != c.expected {
			t.Errorf("Expected %s, got %s", c.expected, got)
		}
	}
}
