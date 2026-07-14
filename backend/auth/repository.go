package auth

import "gorm.io/gorm"

type Repository interface {
	FindByEmail(email string) (*User, error)
	FindByEmailExceptID(email string, id uint) (*User, error)
	FindByID(id uint) (*User, error)
	Create(user *User) error
	Save(user *User) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) GormRepository {
	return GormRepository{db: db}
}

func (r GormRepository) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r GormRepository) FindByEmailExceptID(email string, id uint) (*User, error) {
	var user User
	if err := r.db.Where("email = ? AND id <> ?", email, id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r GormRepository) FindByID(id uint) (*User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r GormRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r GormRepository) Save(user *User) error {
	return r.db.Save(user).Error
}
