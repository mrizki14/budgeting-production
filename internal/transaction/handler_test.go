package transaction

import "testing"

func TestTransactionValidateTypeMismatch(t *testing.T) {
	service := Service{}
	if err := service.validateCategory(1, 1, "other"); err == nil {
		t.Fatal("expected validation error")
	}
}
