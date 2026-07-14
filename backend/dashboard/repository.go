package dashboard

import "gorm.io/gorm"

type Repository interface {
	SumTransactions(userID uint, transactionType string) (float64, error)
	ExpenseBreakdown(userID uint, month int, year int) ([]ExpenseBreakdown, error)
	BudgetRows(userID uint, month int, year int) ([]BudgetRow, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) GormRepository {
	return GormRepository{db: db}
}

func (r GormRepository) SumTransactions(userID uint, transactionType string) (float64, error) {
	var total float64
	err := r.db.Table("transactions").
		Where("user_id = ? AND type = ?", userID, transactionType).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r GormRepository) ExpenseBreakdown(userID uint, month int, year int) ([]ExpenseBreakdown, error) {
	var rows []ExpenseBreakdown
	err := r.db.Table("transactions").
		Select("categories.name AS category_name, SUM(transactions.amount) AS total_amount").
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.user_id = ?", userID).
		Where("transactions.type = ?", "expense").
		Where("MONTH(transactions.date) = ? AND YEAR(transactions.date) = ?", month, year).
		Group("categories.name").
		Order("total_amount DESC").
		Scan(&rows).Error
	return rows, err
}

func (r GormRepository) BudgetRows(userID uint, month int, year int) ([]BudgetRow, error) {
	var rows []BudgetRow
	err := r.db.Table("budgets").
		Select("categories.name AS category, COALESCE(SUM(transactions.amount), 0) AS spent, budgets.limit_amount AS `limit`").
		Joins("JOIN categories ON categories.id = budgets.category_id").
		Joins("LEFT JOIN transactions ON transactions.category_id = budgets.category_id AND transactions.user_id = budgets.user_id AND transactions.type = 'expense' AND MONTH(transactions.date) = ? AND YEAR(transactions.date) = ?", month, year).
		Where("budgets.user_id = ? AND budgets.month = ? AND budgets.year = ?", userID, month, year).
		Group("budgets.id, categories.name, budgets.limit_amount").
		Order("budgets.category_id ASC").
		Scan(&rows).Error
	return rows, err
}
