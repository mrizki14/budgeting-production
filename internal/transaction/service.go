package transaction

import (
	"errors"
	"time"

	"budgeting-app/golang/internal/category"
)

var (
	ErrForbidden         = errors.New("forbidden")
	ErrCategoryOwnership = errors.New("category must belong to current user")
	ErrTransactionType   = errors.New("transaction type must match category type")
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
	ErrInvalidType       = errors.New("type must be income or expense")
)

type Service struct {
	repo         Repository
	categoryRepo category.Repository
}

func NewService(repo Repository, categoryRepo category.Repository) Service {
	return Service{repo: repo, categoryRepo: categoryRepo}
}

func (s Service) List(userID uint, txType string, categoryID uint) ([]Transaction, error) {
	return s.repo.ListByUser(userID, txType, categoryID)
}

func (s Service) Get(userID uint, id uint) (*Transaction, error) {
	transactionRecord, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if transactionRecord.UserID != userID {
		return nil, ErrForbidden
	}

	return transactionRecord, nil
}

func (s Service) Create(userID uint, categoryID uint, amount float64, txType string, date time.Time, description string) (*Transaction, error) {
	if err := s.validateCategory(userID, categoryID, txType); err != nil {
		return nil, err
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	transaction := &Transaction{
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      amount,
		Type:        txType,
		Date:        date,
		Description: description,
	}

	return transaction, s.repo.Create(transaction)
}

func (s Service) Update(userID uint, id uint, categoryID uint, amount float64, txType string, date time.Time, description string) (*Transaction, error) {
	transaction, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if transaction.UserID != userID {
		return nil, ErrForbidden
	}
	if err := s.validateCategory(userID, categoryID, txType); err != nil {
		return nil, err
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	transaction.CategoryID = categoryID
	transaction.Amount = amount
	transaction.Type = txType
	transaction.Date = date
	transaction.Description = description

	return transaction, s.repo.Save(transaction)
}

func (s Service) Delete(userID uint, id uint) error {
	transaction, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if transaction.UserID != userID {
		return ErrForbidden
	}

	return s.repo.Delete(transaction)
}

func (s Service) validateCategory(userID uint, categoryID uint, txType string) error {
	if txType != "income" && txType != "expense" {
		return ErrInvalidType
	}

	categoryRecord, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return err
	}
	if categoryRecord.UserID != userID {
		return ErrCategoryOwnership
	}
	if categoryRecord.Type != txType {
		return ErrTransactionType
	}

	return nil
}
