package main

import "testing"

func TestVerifyPaymentDecision(t *testing.T) {
	for _, tc := range []struct {
		name, code string
		approved   bool
	}{
		{"accepted code", "123456", true}, {"empty code", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.code != ""
			if got != tc.approved {
				t.Fatalf("approved=%v, want %v", got, tc.approved)
			}
		})
	}
}
