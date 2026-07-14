package category

import "errors"

var (
	ErrInvalidCategoryType = errors.New("type must be income or expense")
	ErrCategoryForbidden   = errors.New("forbidden")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) List(userID uint) ([]Category, error) {
	return s.repo.ListByUser(userID)
}

func (s Service) Get(userID uint, id uint) (*Category, error) {
	categoryRecord, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if categoryRecord.UserID != userID {
		return nil, ErrCategoryForbidden
	}

	return categoryRecord, nil
}

func (s Service) Create(userID uint, name string, categoryType string) (*Category, error) {
	if err := s.validateType(categoryType); err != nil {
		return nil, err
	}

	category := &Category{
		UserID: userID,
		Name:   name,
		Type:   categoryType,
	}

	return category, s.repo.Create(category)
}

func (s Service) Update(userID uint, id uint, name string, categoryType string) (*Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category.UserID != userID {
		return nil, ErrCategoryForbidden
	}
	if err := s.validateType(categoryType); err != nil {
		return nil, err
	}

	category.Name = name
	category.Type = categoryType
	return category, s.repo.Save(category)
}

func (s Service) Delete(userID uint, id uint) error {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if category.UserID != userID {
		return ErrCategoryForbidden
	}

	return s.repo.Delete(category)
}

func (s Service) validateType(categoryType string) error {
	if categoryType != "income" && categoryType != "expense" {
		return ErrInvalidCategoryType
	}

	return nil
}
