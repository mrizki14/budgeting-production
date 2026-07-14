package transaction

import (
	"errors"
	"testing"
)

type transactionServiceRepository struct {
	record *Transaction
	err    error
}

func (r transactionServiceRepository) ListByUser(uint, string, uint) ([]Transaction, error) {
	return nil, nil
}
func (r transactionServiceRepository) FindByID(uint) (*Transaction, error) { return r.record, r.err }
func (r transactionServiceRepository) Create(*Transaction) error           { return nil }
func (r transactionServiceRepository) Save(*Transaction) error             { return nil }
func (r transactionServiceRepository) Delete(*Transaction) error           { return nil }

func TestGetReturnsOwnedTransaction(t *testing.T) {
	service := NewService(transactionServiceRepository{record: &Transaction{ID: 7, UserID: 3}}, nil)

	got, err := service.Get(3, 7)

	if err != nil || got.ID != 7 {
		t.Fatalf("expected owned transaction, got %#v, %v", got, err)
	}
}

func TestGetRejectsTransactionOwnedByAnotherUser(t *testing.T) {
	service := NewService(transactionServiceRepository{record: &Transaction{ID: 7, UserID: 9}}, nil)

	_, err := service.Get(3, 7)

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}
