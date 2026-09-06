package rollout

import (
	"strconv"
	"testing"
)

func TestEvaluateDeterministic(t *testing.T) {
	const key = "feature"
	const user = "user-42"
	const percent = 50

	first := Evaluate(key, user, percent)
	for i := 0; i < 1000; i++ {
		if got := Evaluate(key, user, percent); got != first {
			t.Fatalf("Evaluate not deterministic: got %v then %v", first, got)
		}
	}
}

func TestEvaluateZeroPercentAlwaysFalse(t *testing.T) {
	for _, user := range []string{"a", "b", "user-42", "", "xyz"} {
		if Evaluate("feature", user, 0) {
			t.Fatalf("expected false for percent 0 and user %q", user)
		}
	}
}

func TestEvaluateHundredPercentAlwaysTrue(t *testing.T) {
	for _, user := range []string{"a", "b", "user-42", "", "xyz"} {
		if !Evaluate("feature", user, 100) {
			t.Fatalf("expected true for percent 100 and user %q", user)
		}
	}
}

func TestEvaluateDifferentUsersMayDiffer(t *testing.T) {
	// The distribution is deterministic per user; a sample of users should
	// produce both outcomes for a mid-range percentage, proving the hash
	// actually varies across users rather than being constant.
	sawTrue := false
	sawFalse := false
	for i := 0; i < 1000 && !(sawTrue && sawFalse); i++ {
		if Evaluate("feature", "user-"+strconv.Itoa(i), 50) {
			sawTrue = true
		} else {
			sawFalse = true
		}
	}
	if !sawTrue || !sawFalse {
		t.Fatalf("expected a mix of outcomes; saw true=%v false=%v", sawTrue, sawFalse)
	}
}
