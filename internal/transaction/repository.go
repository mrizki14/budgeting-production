package transaction

import "gorm.io/gorm"

type Repository interface {
	ListByUser(userID uint, txType string, categoryID uint) ([]Transaction, error)
	FindByID(id uint) (*Transaction, error)
	Create(transaction *Transaction) error
	Save(transaction *Transaction) error
	Delete(transaction *Transaction) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) GormRepository {
	return GormRepository{db: db}
}

func (r GormRepository) ListByUser(userID uint, txType string, categoryID uint) ([]Transaction, error) {
	query := r.db.Preload("Category").Where("user_id = ?", userID).Order("date desc").Order("id desc")
	if txType == "income" || txType == "expense" {
		query = query.Where("type = ?", txType)
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	var transactions []Transaction
	err := query.Find(&transactions).Error
	return transactions, err
}

func (r GormRepository) FindByID(id uint) (*Transaction, error) {
	var transaction Transaction
	if err := r.db.Preload("Category").First(&transaction, id).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r GormRepository) Create(transaction *Transaction) error {
	return r.db.Create(transaction).Error
}

func (r GormRepository) Save(transaction *Transaction) error {
	return r.db.Save(transaction).Error
}

func (r GormRepository) Delete(transaction *Transaction) error {
	return r.db.Delete(transaction).Error
}
