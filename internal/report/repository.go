package report

import "gorm.io/gorm"

type CategoryTotal struct {
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
}

type Summary struct {
	TotalIncome       float64         `json:"total_income"`
	TotalExpenses     float64         `json:"total_expenses"`
	NetSavings        float64         `json:"net_savings"`
	IncomeByCategory  []CategoryTotal `json:"income_by_category"`
	ExpenseByCategory []CategoryTotal `json:"expense_by_category"`
	Month             int             `json:"month"`
	Year              int             `json:"year"`
}

type Repository interface {
	SumByType(userID uint, month int, year int, txType string) (float64, error)
	CategoryBreakdown(userID uint, month int, year int, txType string) ([]CategoryTotal, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) GormRepository {
	return GormRepository{db: db}
}

func (r GormRepository) SumByType(userID uint, month int, year int, txType string) (float64, error) {
	var total float64
	err := r.db.Table("transactions").
		Where("user_id = ? AND type = ? AND MONTH(date) = ? AND YEAR(date) = ?", userID, txType, month, year).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r GormRepository) CategoryBreakdown(userID uint, month int, year int, txType string) ([]CategoryTotal, error) {
	var result []CategoryTotal
	err := r.db.Table("transactions").
		Select("categories.name as category_name, COALESCE(SUM(transactions.amount), 0) as total_amount").
		Joins("join categories on categories.id = transactions.category_id").
		Where("transactions.user_id = ? AND transactions.type = ? AND MONTH(transactions.date) = ? AND YEAR(transactions.date) = ?", userID, txType, month, year).
		Group("categories.name").
		Order("total_amount desc").
		Scan(&result).Error
	return result, err
}
