// # Go Finance Application

// Create a simple finance application in Go that allows users to manage their income, expenses, and investments. The application should have the following features:

// ## Data Structures
// 1. `Income` struct with fields for `Date`, `Source`, and `Amount`.
// 2. `Expense` struct with fields for `Date`, `Category`, and `Amount`.
// 3. `Investment` struct with fields for `Date`, `Asset`, and `Value`.

// ## Core Functionality
// 1. **Income Management**:
//    - Add a new income entry
//    - View a list of all income entries
//    - Calculate total income for a given period (e.g., month, year)

// 2. **Expense Management**:
//    - Add a new expense entry
//    - View a list of all expense entries
//    - Calculate total expenses for a given period (e.g., month, year)
//    - Categorize expenses (e.g., rent, groceries, utilities)

// 3. **Investment Management**:
//    - Add a new investment entry
//    - View a list of all investment entries
//    - Calculate total investment value
//    - Provide basic analytics (e.g., average investment value, highest/lowest investment value)

// 4. **Financial Reports**:
//    - Generate a monthly/yearly financial report that includes:
//      - Total income
//      - Total expenses
//      - Total investments
//      - Net profit/loss

// ## Requirements
// 1. Use Go's built-in `time` package for handling dates.
// 2. Implement all CRUD (Create, Read, Update, Delete) operations for the data structures.
// 3. Provide command-line interface (CLI) for users to interact with the application.
// 4. Use Go's `flag` package to handle command-line arguments.
// 5. Implement basic error handling and input validation.
// 6. Write unit tests for the core functionality.

// ## Bonus Features (Optional)
// 1. Persistent storage using a file or a simple database (e.g., SQLite).
// 2. Visualize financial data using a charting library (e.g., `gonum/plot`).
// 3. Implement basic budgeting features (e.g., set monthly/yearly budgets, track budget progress).
// 4. Allow users to export financial reports in various formats (e.g., CSV, PDF).

// This exercise covers the following Go concepts:
// - Structs and methods
// - Slices and arrays
// - Command-line interface (CLI) with the `flag` package
// - Error handling
// - Unit testing
// - Basic file I/O or database integration (optional)
// - Data visualization (optional)

// Let me know if you have any questions or need further clarification on this exercise!

// main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Income struct {
	Date   time.Time
	Source string
	Amount float64
}

type Expense struct {
	Date     time.Time
	Category string
	Amount   float64
}

type Investment struct {
	Date  time.Time
	Asset string
	Value float64
}

type FinanceManager struct {
	incomes     []Income
	expenses    []Expense
	investments []Investment
}

func NewFinanceManager() *FinanceManager {
	return &FinanceManager{
		incomes:     make([]Income, 0),
		expenses:    make([]Expense, 0),
		investments: make([]Investment, 0),
	}
}

// Income methods
func (fm *FinanceManager) AddIncome(date time.Time, source string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	income := Income{
		Date:   date,
		Source: source,
		Amount: amount,
	}
	fm.incomes = append(fm.incomes, income)
	return nil
}

func (fm *FinanceManager) GetTotalIncome(startDate, endDate time.Time) float64 {
	var total float64
	for _, income := range fm.incomes {
		if (income.Date.After(startDate) || income.Date.Equal(startDate)) &&
			(income.Date.Before(endDate) || income.Date.Equal(endDate)) {
			total += income.Amount
		}
	}
	return total
}

// Expense methods
func (fm *FinanceManager) AddExpense(date time.Time, category string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	expense := Expense{
		Date:     date,
		Category: category,
		Amount:   amount,
	}
	fm.expenses = append(fm.expenses, expense)
	return nil
}

func (fm *FinanceManager) GetTotalExpenses(startDate, endDate time.Time) float64 {
	var total float64
	for _, expense := range fm.expenses {
		if (expense.Date.After(startDate) || expense.Date.Equal(startDate)) &&
			(expense.Date.Before(endDate) || expense.Date.Equal(endDate)) {
			total += expense.Amount
		}
	}
	return total
}

// Investment methods
func (fm *FinanceManager) AddInvestment(date time.Time, asset string, value float64) error {
	if value <= 0 {
		return fmt.Errorf("value must be positive")
	}
	investment := Investment{
		Date:  date,
		Asset: asset,
		Value: value,
	}
	fm.investments = append(fm.investments, investment)
	return nil
}

func (fm *FinanceManager) GetTotalInvestments() float64 {
	var total float64
	for _, investment := range fm.investments {
		total += investment.Value
	}
	return total
}

// Report generation
func (fm *FinanceManager) GenerateMonthlyReport(year int, month time.Month) string {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	totalIncome := fm.GetTotalIncome(startDate, endDate)
	totalExpenses := fm.GetTotalExpenses(startDate, endDate)
	totalInvestments := fm.GetTotalInvestments()
	netProfit := totalIncome - totalExpenses

	report := fmt.Sprintf(`
Financial Report for %s %d
-------------------------
Total Income:     $%.2f
Total Expenses:   $%.2f
Total Investments: $%.2f
Net Profit/Loss:  $%.2f
`, month.String(), year, totalIncome, totalExpenses, totalInvestments, netProfit)

	return report
}

func financeApp() {
	fm := NewFinanceManager()
	reader := bufio.NewReader(os.Stdin)

	// Command line flags
	mode := flag.String("mode", "interactive", "Mode of operation: interactive or report")
	flag.Parse()

	if *mode == "interactive" {
		for {
			fmt.Println("\nFinance Manager")
			fmt.Println("1. Add Income")
			fmt.Println("2. Add Expense")
			fmt.Println("3. Add Investment")
			fmt.Println("4. Generate Monthly Report")
			fmt.Println("5. Exit")
			fmt.Print("Choose an option: ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "1":
				handleAddIncome(fm, reader)
			case "2":
				handleAddExpense(fm, reader)
			case "3":
				handleAddInvestment(fm, reader)
			case "4":
				handleGenerateReport(fm, reader)
			case "5":
				fmt.Println("Goodbye!")
				return
			default:
				fmt.Println("Invalid option")
			}
		}
	}
}

func handleAddIncome(fm *FinanceManager, reader *bufio.Reader) {
	fmt.Print("Enter source: ")
	source, _ := reader.ReadString('\n')
	source = strings.TrimSpace(source)

	fmt.Print("Enter amount: ")
	amountStr, _ := reader.ReadString('\n')
	amount, err := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
	if err != nil {
		fmt.Println("Invalid amount")
		return
	}

	err = fm.AddIncome(time.Now(), source, amount)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Income added successfully")
}

func handleAddExpense(fm *FinanceManager, reader *bufio.Reader) {
	fmt.Print("Enter category: ")
	category, _ := reader.ReadString('\n')
	category = strings.TrimSpace(category)

	fmt.Print("Enter amount: ")
	amountStr, _ := reader.ReadString('\n')
	amount, err := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
	if err != nil {
		fmt.Println("Invalid amount")
		return
	}

	err = fm.AddExpense(time.Now(), category, amount)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Expense added successfully")
}

func handleAddInvestment(fm *FinanceManager, reader *bufio.Reader) {
	fmt.Print("Enter asset name: ")
	asset, _ := reader.ReadString('\n')
	asset = strings.TrimSpace(asset)

	fmt.Print("Enter value: ")
	valueStr, _ := reader.ReadString('\n')
	value, err := strconv.ParseFloat(strings.TrimSpace(valueStr), 64)
	if err != nil {
		fmt.Println("Invalid value")
		return
	}

	err = fm.AddInvestment(time.Now(), asset, value)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Investment added successfully")
}

func handleGenerateReport(fm *FinanceManager, reader *bufio.Reader) {
	currentTime := time.Now()
	report := fm.GenerateMonthlyReport(currentTime.Year(), currentTime.Month())
	fmt.Println(report)
}
