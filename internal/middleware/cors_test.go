package middleware

import (
	"testing"
)

func TestIsOriginAllowed(t *testing.T) {
	allowed := "https://ip.ipw.3x.cx"

	tests := []struct {
		origin   string
		expected bool
	}{
		{"https://ip.ipw.3x.cx", true},
		{"https://ipw.3x.cx", true},
		{"https://4.ipw.3x.cx", true},
		{"https://6.ipw.3x.cx", true},
		{"https://test.ipw.3x.cx", true},
		{"http://ipw.3x.cx", true},
		{"https://other.com", false},
		{"https://fakeipw.3x.cx", false},
		{"", false},
	}

	for _, tt := range tests {
		got := isOriginAllowed(tt.origin, allowed)
		if got != tt.expected {
			t.Errorf("isOriginAllowed(%q, %q) = %v; want %v", tt.origin, allowed, got, tt.expected)
		}
	}

	// Test wildcard
	if !isOriginAllowed("https://anything.com", "*") {
		t.Errorf("isOriginAllowed with * should allow any origin")
	}
}
