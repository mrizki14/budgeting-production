package budget

import (
	"errors"
	"testing"
)

type budgetServiceRepository struct {
	record *Budget
	err    error
}

func (r budgetServiceRepository) ListByUser(uint) ([]Budget, error) { return nil, nil }
func (r budgetServiceRepository) FindByID(uint) (*Budget, error)    { return r.record, r.err }
func (r budgetServiceRepository) Create(*Budget) error              { return nil }
func (r budgetServiceRepository) Save(*Budget) error                { return nil }
func (r budgetServiceRepository) Delete(*Budget) error              { return nil }
func (r budgetServiceRepository) ExistsForPeriod(uint, uint, int, int, uint) (bool, error) {
	return false, nil
}

func TestGetReturnsOwnedBudget(t *testing.T) {
	service := NewService(budgetServiceRepository{record: &Budget{ID: 7, UserID: 3}}, nil)

	got, err := service.Get(3, 7)

	if err != nil || got.ID != 7 {
		t.Fatalf("expected owned budget, got %#v, %v", got, err)
	}
}

func TestGetRejectsBudgetOwnedByAnotherUser(t *testing.T) {
	service := NewService(budgetServiceRepository{record: &Budget{ID: 7, UserID: 9}}, nil)

	_, err := service.Get(3, 7)

	if !errors.Is(err, ErrBudgetForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}
