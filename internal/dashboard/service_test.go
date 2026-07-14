package dashboard

import (
	"testing"
	"time"
)

type dashboardRepository struct {
	income    float64
	expenses  float64
	breakdown []ExpenseBreakdown
	budgets   []BudgetRow
}

func (r dashboardRepository) SumTransactions(_ uint, transactionType string) (float64, error) {
	if transactionType == "income" {
		return r.income, nil
	}
	return r.expenses, nil
}

func (r dashboardRepository) ExpenseBreakdown(uint, int, int) ([]ExpenseBreakdown, error) {
	return r.breakdown, nil
}

func (r dashboardRepository) BudgetRows(uint, int, int) ([]BudgetRow, error) {
	return r.budgets, nil
}

func TestSummaryCalculatesBalanceAndBudgetStatus(t *testing.T) {
	repo := dashboardRepository{
		income:   8_000_000,
		expenses: 3_000_000,
		budgets: []BudgetRow{
			{Category: "Food", Spent: 760_000, Limit: 1_000_000},
			{Category: "Travel", Spent: 950_000, Limit: 1_000_000},
		},
	}
	service := NewService(repo)
	service.now = func() time.Time {
		return time.Date(2026, time.July, 14, 0, 0, 0, 0, time.UTC)
	}

	got, err := service.Summary(4)

	if err != nil {
		t.Fatal(err)
	}
	if got.TotalBalance != 5_000_000 {
		t.Fatalf("expected balance 5000000, got %v", got.TotalBalance)
	}
	if got.BudgetUsage[0].Percentage != 76 || got.BudgetUsage[0].Status != "warn" {
		t.Fatalf("expected warn at 76%%, got %#v", got.BudgetUsage[0])
	}
	if got.BudgetUsage[1].Percentage != 95 || got.BudgetUsage[1].Status != "danger" {
		t.Fatalf("expected danger at 95%%, got %#v", got.BudgetUsage[1])
	}
	if got.CurrentMonthLabel != "July 2026" {
		t.Fatalf("expected July 2026, got %q", got.CurrentMonthLabel)
	}
}

func TestSummaryCapsBudgetUsageAtOneHundredPercent(t *testing.T) {
	service := NewService(dashboardRepository{budgets: []BudgetRow{
		{Category: "Food", Spent: 1_500_000, Limit: 1_000_000},
		{Category: "No limit", Spent: 50_000, Limit: 0},
	}})

	got, err := service.Summary(4)

	if err != nil {
		t.Fatal(err)
	}
	if got.BudgetUsage[0].Percentage != 100 || got.BudgetUsage[1].Percentage != 0 {
		t.Fatalf("unexpected percentages: %#v", got.BudgetUsage)
	}
}

func TestSummaryReturnsEmptyCollectionsInsteadOfNull(t *testing.T) {
	service := NewService(dashboardRepository{})

	got, err := service.Summary(4)

	if err != nil {
		t.Fatal(err)
	}
	if got.ExpenseByCategory == nil || got.BudgetUsage == nil {
		t.Fatalf("dashboard collections must encode as JSON arrays: %#v", got)
	}
}
