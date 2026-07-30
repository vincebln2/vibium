package main

import "testing"

func TestWSSchemeSuggestion(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"https://example.com", "wss://example.com"},
		{"http://example.com", "ws://example.com"},
		{"http://example.com:9222/session", "ws://example.com:9222/session"},
		{"HTTPS://example.com/Path", "wss://example.com/Path"},
		{"Http://example.com", "ws://example.com"},
		{"wss://example.com", ""},
		{"ws://example.com", ""},
		{"example.com", ""},
		{"", ""},
	}
	for _, c := range cases {
		if got := wsSchemeSuggestion(c.in); got != c.want {
			t.Errorf("wsSchemeSuggestion(%q): got %q, want %q", c.in, got, c.want)
		}
	}
}
