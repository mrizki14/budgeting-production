package dashboard

type ExpenseBreakdown struct {
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
}

type BudgetRow struct {
	Category string
	Spent    float64
	Limit    float64
}

type BudgetUsage struct {
	Category   string  `json:"category"`
	Spent      float64 `json:"spent"`
	Limit      float64 `json:"limit"`
	Percentage int     `json:"percentage"`
	Status     string  `json:"status"`
}

type Summary struct {
	TotalBalance      float64            `json:"total_balance"`
	TotalIncome       float64            `json:"total_income"`
	TotalExpenses     float64            `json:"total_expenses"`
	ExpenseByCategory []ExpenseBreakdown `json:"expense_by_category"`
	BudgetUsage       []BudgetUsage      `json:"budget_usage"`
	CurrentMonthLabel string             `json:"current_month_label"`
}
