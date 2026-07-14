package category

import "gorm.io/gorm"

type Repository interface {
	ListByUser(userID uint) ([]Category, error)
	FindByID(id uint) (*Category, error)
	Create(category *Category) error
	Save(category *Category) error
	Delete(category *Category) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) GormRepository {
	return GormRepository{db: db}
}

func (r GormRepository) ListByUser(userID uint) ([]Category, error) {
	var categories []Category
	err := r.db.Where("user_id = ?", userID).Order("type asc").Order("name asc").Find(&categories).Error
	return categories, err
}

func (r GormRepository) FindByID(id uint) (*Category, error) {
	var category Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r GormRepository) Create(category *Category) error {
	return r.db.Create(category).Error
}

func (r GormRepository) Save(category *Category) error {
	return r.db.Save(category).Error
}

func (r GormRepository) Delete(category *Category) error {
	return r.db.Delete(category).Error
}
