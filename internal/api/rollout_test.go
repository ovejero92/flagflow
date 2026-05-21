package api

import "testing"

func TestEvaluateRollout_Consistent(t *testing.T) {
	a := evaluateRollout("checkout-v2", "user-42", 50)
	b := evaluateRollout("checkout-v2", "user-42", 50)
	if a != b {
		t.Fatal("rollout must be consistent for same user and flag")
	}
}

func TestEvaluateRollout_Boundaries(t *testing.T) {
	if !evaluateRollout("f", "u", 100) {
		t.Fatal("100% rollout should always enable")
	}
	if evaluateRollout("f", "u", 0) {
		t.Fatal("0% rollout should never enable")
	}
}
