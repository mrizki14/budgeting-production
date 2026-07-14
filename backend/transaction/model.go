package transaction

import (
	"time"

	"budgeting-app/golang/backend/category"
)

type Transaction struct {
	ID          uint              `gorm:"column:id;primaryKey" json:"id"`
	UserID      uint              `gorm:"column:user_id" json:"user_id"`
	CategoryID  uint              `gorm:"column:category_id" json:"category_id"`
	Amount      float64           `gorm:"column:amount" json:"amount"`
	Type        string            `gorm:"column:type" json:"type"`
	Date        time.Time         `gorm:"column:date" json:"date"`
	Description string            `gorm:"column:description" json:"description"`
	CreatedAt   time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"column:updated_at" json:"updated_at"`
	Category    category.Category `gorm:"foreignKey:CategoryID" json:"category"`
}

func (Transaction) TableName() string {
	return "transactions"
}
