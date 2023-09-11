package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"exit  ", []string{"exit"}},
		{"    help     ", []string{"help"}},
		{"    help me     ", []string{"help", "me"}},
		{"", []string{""}},
		{"    ", []string{""}},
		{"HELP", []string{"help"}},
	}

	for _, c := range cases {
		got := cleanInput(c.input)
		if len(got) != len(c.expected) {
			t.Errorf("Lengths differ. Expected %d, got %d", len(c.expected), len(got))
		}
		for i := range got {
			if got[i] != c.expected[i] {
				t.Errorf("At index %d: expected %s, got %s", i, c.expected[i], got[i])
			}
		}
	}
}
