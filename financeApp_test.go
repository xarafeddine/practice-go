// main_test.go
package main

import (
	"testing"
	"time"
)

func TestIncomeManagement(t *testing.T) {
	fm := NewFinanceManager()

	// Test adding valid income
	err := fm.AddIncome(time.Now(), "Salary", 5000.0)
	if err != nil {
		t.Errorf("Failed to add valid income: %v", err)
	}

	// Test adding invalid income
	err = fm.AddIncome(time.Now(), "Bonus", -100.0)
	if err == nil {
		t.Error("Expected error for negative income amount")
	}

	// Test total income calculation
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now().AddDate(0, 0, 1)
	total := fm.GetTotalIncome(startDate, endDate)
	if total != 5000.0 {
		t.Errorf("Expected total income of 5000.0, got %f", total)
	}
}

func TestExpenseManagement(t *testing.T) {
	fm := NewFinanceManager()

	// Test adding valid expense
	err := fm.AddExpense(time.Now(), "Rent", 1000.0)
	if err != nil {
		t.Errorf("Failed to add valid expense: %v", err)
	}

	// Test adding invalid expense
	err = fm.AddExpense(time.Now(), "Food", -50.0)
	if err == nil {
		t.Error("Expected error for negative expense amount")
	}

	// Test total expenses calculation
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now().AddDate(0, 0, 1)
	total := fm.GetTotalExpenses(startDate, endDate)
	if total != 1000.0 {
		t.Errorf("Expected total expenses of 1000.0, got %f", total)
	}
}

func TestInvestmentManagement(t *testing.T) {
	fm := NewFinanceManager()

	// Test adding valid investment
	err := fm.AddInvestment(time.Now(), "Stocks", 2000.0)
	if err != nil {
		t.Errorf("Failed to add valid investment: %v", err)
	}

	// Test adding invalid investment
	err = fm.AddInvestment(time.Now(), "Bonds", -500.0)
	if err == nil {
		t.Error("Expected error for negative investment value")
	}

	// Test total investments calculation
	total := fm.GetTotalInvestments()
	if total != 2000.0 {
		t.Errorf("Expected total investments of 2000.0, got %f", total)
	}
}

func TestMonthlyReport(t *testing.T) {
	fm := NewFinanceManager()

	// Add some test data
	currentTime := time.Now()
	fm.AddIncome(currentTime, "Salary", 5000.0)
	fm.AddExpense(currentTime, "Rent", 1000.0)
	fm.AddInvestment(currentTime, "Stocks", 2000.0)

	// Generate and check report
	report := fm.GenerateMonthlyReport(currentTime.Year(), currentTime.Month())
	if report == "" {
		t.Error("Expected non-empty report")
	}
}
