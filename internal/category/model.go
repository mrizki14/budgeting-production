package category

import "time"

type Category struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Type      string    `gorm:"column:type" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}
