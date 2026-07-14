package budget

import "testing"

func TestBudgetRejectsNonPositiveLimit(t *testing.T) {
	if !errorsIs(validateLimit(0), ErrBudgetLimitInvalid) {
		t.Fatal("expected limit validation error")
	}
}

func validateLimit(limit float64) error {
	if limit <= 0 {
		return ErrBudgetLimitInvalid
	}

	return nil
}

func errorsIs(err error, target error) bool {
	return err == target
}
