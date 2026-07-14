package budget

import (
	"errors"

	"budgeting-app/golang/backend/category"
)

var (
	ErrBudgetForbidden       = errors.New("forbidden")
	ErrBudgetCategoryOwner   = errors.New("category must belong to current user")
	ErrBudgetExpenseOnly     = errors.New("budget category must be expense")
	ErrBudgetLimitInvalid    = errors.New("limit amount must be greater than zero")
	ErrBudgetDuplicatePeriod = errors.New("budget already exists for category and period")
)

type Service struct {
	repo         Repository
	categoryRepo category.Repository
}

func NewService(repo Repository, categoryRepo category.Repository) Service {
	return Service{repo: repo, categoryRepo: categoryRepo}
}

func (s Service) List(userID uint) ([]Budget, error) {
	return s.repo.ListByUser(userID)
}

func (s Service) Get(userID uint, id uint) (*Budget, error) {
	budgetRecord, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if budgetRecord.UserID != userID {
		return nil, ErrBudgetForbidden
	}

	return budgetRecord, nil
}

func (s Service) Create(userID uint, categoryID uint, month int, year int, limitAmount float64) (*Budget, error) {
	if err := s.validateCreate(userID, categoryID, month, year, limitAmount, 0); err != nil {
		return nil, err
	}

	budget := &Budget{
		UserID:      userID,
		CategoryID:  categoryID,
		Month:       month,
		Year:        year,
		LimitAmount: limitAmount,
	}

	return budget, s.repo.Create(budget)
}

func (s Service) Update(userID uint, id uint, categoryID uint, month int, year int, limitAmount float64) (*Budget, error) {
	budget, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if budget.UserID != userID {
		return nil, ErrBudgetForbidden
	}
	if err := s.validateCreate(userID, categoryID, month, year, limitAmount, id); err != nil {
		return nil, err
	}

	budget.CategoryID = categoryID
	budget.Month = month
	budget.Year = year
	budget.LimitAmount = limitAmount

	return budget, s.repo.Save(budget)
}

func (s Service) Delete(userID uint, id uint) error {
	budget, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if budget.UserID != userID {
		return ErrBudgetForbidden
	}

	return s.repo.Delete(budget)
}

func (s Service) validateCreate(userID uint, categoryID uint, month int, year int, limitAmount float64, ignoreID uint) error {
	categoryRecord, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return err
	}
	if categoryRecord.UserID != userID {
		return ErrBudgetCategoryOwner
	}
	if categoryRecord.Type != "expense" {
		return ErrBudgetExpenseOnly
	}
	if limitAmount <= 0 {
		return ErrBudgetLimitInvalid
	}

	exists, err := s.repo.ExistsForPeriod(userID, categoryID, month, year, ignoreID)
	if err != nil {
		return err
	}
	if exists {
		return ErrBudgetDuplicatePeriod
	}

	return nil
}
