package dashboard

import (
	"math"
	"time"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) Service {
	return Service{repo: repo, now: time.Now}
}

func (s Service) Summary(userID uint) (Summary, error) {
	now := s.now()
	totalIncome, err := s.repo.SumTransactions(userID, "income")
	if err != nil {
		return Summary{}, err
	}
	totalExpenses, err := s.repo.SumTransactions(userID, "expense")
	if err != nil {
		return Summary{}, err
	}
	expenseByCategory, err := s.repo.ExpenseBreakdown(userID, int(now.Month()), now.Year())
	if err != nil {
		return Summary{}, err
	}
	if expenseByCategory == nil {
		expenseByCategory = make([]ExpenseBreakdown, 0)
	}
	rows, err := s.repo.BudgetRows(userID, int(now.Month()), now.Year())
	if err != nil {
		return Summary{}, err
	}

	usage := make([]BudgetUsage, 0, len(rows))
	for _, row := range rows {
		percentage := 0
		if row.Limit > 0 {
			percentage = min(100, int(math.Round((row.Spent/row.Limit)*100)))
		}
		status := "good"
		if percentage >= 90 {
			status = "danger"
		} else if percentage >= 75 {
			status = "warn"
		}
		usage = append(usage, BudgetUsage{
			Category: row.Category, Spent: row.Spent, Limit: row.Limit,
			Percentage: percentage, Status: status,
		})
	}

	return Summary{
		TotalBalance:      totalIncome - totalExpenses,
		TotalIncome:       totalIncome,
		TotalExpenses:     totalExpenses,
		ExpenseByCategory: expenseByCategory,
		BudgetUsage:       usage,
		CurrentMonthLabel: now.Format("January 2006"),
	}, nil
}
