package category

import (
	"errors"
	"testing"
)

type categoryServiceRepository struct {
	record *Category
	err    error
}

func (r categoryServiceRepository) ListByUser(uint) ([]Category, error) { return nil, nil }
func (r categoryServiceRepository) FindByID(uint) (*Category, error)    { return r.record, r.err }
func (r categoryServiceRepository) Create(*Category) error              { return nil }
func (r categoryServiceRepository) Save(*Category) error                { return nil }
func (r categoryServiceRepository) Delete(*Category) error              { return nil }

func TestGetReturnsOwnedCategory(t *testing.T) {
	service := NewService(categoryServiceRepository{record: &Category{ID: 7, UserID: 3}})

	got, err := service.Get(3, 7)

	if err != nil || got.ID != 7 {
		t.Fatalf("expected owned category, got %#v, %v", got, err)
	}
}

func TestGetRejectsCategoryOwnedByAnotherUser(t *testing.T) {
	service := NewService(categoryServiceRepository{record: &Category{ID: 7, UserID: 9}})

	_, err := service.Get(3, 7)

	if !errors.Is(err, ErrCategoryForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}
