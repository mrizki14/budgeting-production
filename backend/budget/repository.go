package budget

import "gorm.io/gorm"

type Repository interface {
	ListByUser(userID uint) ([]Budget, error)
	FindByID(id uint) (*Budget, error)
	Create(budget *Budget) error
	Save(budget *Budget) error
	Delete(budget *Budget) error
	ExistsForPeriod(userID uint, categoryID uint, month int, year int, ignoreID uint) (bool, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) GormRepository {
	return GormRepository{db: db}
}

func (r GormRepository) ListByUser(userID uint) ([]Budget, error) {
	var budgets []Budget
	err := r.db.Preload("Category").Where("user_id = ?", userID).Order("year desc").Order("month desc").Order("category_id asc").Find(&budgets).Error
	return budgets, err
}

func (r GormRepository) FindByID(id uint) (*Budget, error) {
	var budget Budget
	if err := r.db.Preload("Category").First(&budget, id).Error; err != nil {
		return nil, err
	}

	return &budget, nil
}

func (r GormRepository) Create(budget *Budget) error {
	return r.db.Create(budget).Error
}

func (r GormRepository) Save(budget *Budget) error {
	return r.db.Save(budget).Error
}

func (r GormRepository) Delete(budget *Budget) error {
	return r.db.Delete(budget).Error
}

func (r GormRepository) ExistsForPeriod(userID uint, categoryID uint, month int, year int, ignoreID uint) (bool, error) {
	query := r.db.Model(&Budget{}).
		Where("user_id = ? AND category_id = ? AND month = ? AND year = ?", userID, categoryID, month, year)
	if ignoreID > 0 {
		query = query.Where("id <> ?", ignoreID)
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}
