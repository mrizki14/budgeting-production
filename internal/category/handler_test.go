package category

import "testing"

func TestCategoryTypeValidation(t *testing.T) {
	service := Service{}

	if err := service.validateType("other"); err == nil {
		t.Fatal("expected validation error for invalid category type")
	}
}
