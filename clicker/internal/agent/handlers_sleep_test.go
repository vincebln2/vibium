package agent

import (
	"strings"
	"testing"
)

func TestBrowserSleepRejectsOverMax(t *testing.T) {
	h := NewHandlers("", true, "", nil)
	for _, ms := range []float64{30001, 999999} {
		_, err := h.browserSleep(map[string]interface{}{"ms": ms})
		if err == nil {
			t.Fatalf("ms=%v: expected error, got nil", ms)
		}
		if !strings.Contains(err.Error(), "30000") {
			t.Errorf("ms=%v: error should mention the 30000 limit, got: %v", ms, err)
		}
	}
}

func TestBrowserSleepValid(t *testing.T) {
	h := NewHandlers("", true, "", nil)
	result, err := h.browserSleep(map[string]interface{}{"ms": float64(10)})
	if err != nil {
		t.Fatalf("ms=10: unexpected error: %v", err)
	}
	if got := result.Content[0].Text; got != "Slept for 10 ms" {
		t.Errorf("ms=10: got %q", got)
	}
}
