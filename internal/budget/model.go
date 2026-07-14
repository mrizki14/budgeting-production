package budget

import (
	"time"

	"budgeting-app/golang/internal/category"
)

type Budget struct {
	ID          uint              `gorm:"column:id;primaryKey" json:"id"`
	UserID      uint              `gorm:"column:user_id" json:"user_id"`
	CategoryID  uint              `gorm:"column:category_id" json:"category_id"`
	Month       int               `gorm:"column:month" json:"month"`
	Year        int               `gorm:"column:year" json:"year"`
	LimitAmount float64           `gorm:"column:limit_amount" json:"limit_amount"`
	CreatedAt   time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"column:updated_at" json:"updated_at"`
	Category    category.Category `gorm:"foreignKey:CategoryID" json:"category"`
}

func (Budget) TableName() string {
	return "budgets"
}
