package report

import "time"

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) Service {
	return Service{repo: repo, now: time.Now}
}

func (s Service) Summary(userID uint, month int, year int) (Summary, error) {
	month, year = normalizePeriod(month, year, s.now())

	totalIncome, err := s.repo.SumByType(userID, month, year, "income")
	if err != nil {
		return Summary{}, err
	}
	totalExpenses, err := s.repo.SumByType(userID, month, year, "expense")
	if err != nil {
		return Summary{}, err
	}
	incomeByCategory, err := s.repo.CategoryBreakdown(userID, month, year, "income")
	if err != nil {
		return Summary{}, err
	}
	expenseByCategory, err := s.repo.CategoryBreakdown(userID, month, year, "expense")
	if err != nil {
		return Summary{}, err
	}

	return Summary{
		TotalIncome:       totalIncome,
		TotalExpenses:     totalExpenses,
		NetSavings:        totalIncome - totalExpenses,
		IncomeByCategory:  incomeByCategory,
		ExpenseByCategory: expenseByCategory,
		Month:             month,
		Year:              year,
	}, nil
}

func normalizePeriod(month int, year int, now time.Time) (int, int) {
	if month < 1 || month > 12 {
		month = int(now.Month())
	}

	validYears := map[int]bool{
		now.Year() - 1: true,
		now.Year():     true,
		now.Year() + 1: true,
	}
	if !validYears[year] {
		year = now.Year()
	}

	return month, year
}
