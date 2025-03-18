package report

import (
	"fmt"
	"sort"

	"github.com/xuri/excelize/v2"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

const AccountSummarySheetName = "Account Summary"
const TransactionsHistorySheetName = "Transactions History"

// ExcelStyles holds styling information for Excel sheets
type ExcelStyles struct {
	HeaderStyle   int
	CurrencyStyle int
	DateStyle     int
}

// CreateExcelStyles creates and returns common styles for Excel sheets
func CreateExcelStyles(f *excelize.File) (*ExcelStyles, error) {
	// Header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create header style: %w", err)
	}

	// Currency style
	customCcyStyle := "$#,##0.00;-$#,##0.00"
	currencyStyle, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: &customCcyStyle,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create currency style: %w", err)
	}

	// Date style
	dateStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 14, // mm-dd-yy
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create date style: %w", err)
	}

	return &ExcelStyles{
		HeaderStyle:   headerStyle,
		CurrencyStyle: currencyStyle,
		DateStyle:     dateStyle,
	}, nil
}

// CreateAccountSummarySheet creates the Account Summary sheet in the Excel file
func CreateAccountSummarySheet(f *excelize.File, account *data.Account, transactions []*data.Transaction, styles *ExcelStyles) error {
	sheetName := AccountSummarySheetName
	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		return fmt.Errorf("failed to rename sheet: %w", err)
	}

	// Account information section
	f.SetCellValue(sheetName, "A1", "Account Information")
	titleStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create title style: %w", err)
	}
	f.SetCellStyle(sheetName, "A1", "A1", titleStyle)
	f.MergeCell(sheetName, "A1", "B1")

	// Account Details
	f.SetCellValue(sheetName, "A2", "Account ID:")
	f.SetCellValue(sheetName, "B2", account.ID.String())
	f.SetCellValue(sheetName, "A3", "Account Name:")
	f.SetCellValue(sheetName, "B3", account.Name)
	f.SetCellValue(sheetName, "A4", "Current Balance:")
	f.SetCellValue(sheetName, "B4", float64(account.Balance)/100.0)
	f.SetCellStyle(sheetName, "B4", "B4", styles.CurrencyStyle)
	f.SetCellValue(sheetName, "A5", "Created At:")
	f.SetCellValue(sheetName, "B5", account.CreatedAt.Format("Jan 02, 2006 15:04:05"))
	f.SetCellValue(sheetName, "A6", "Last Updated:")
	f.SetCellValue(sheetName, "B6", account.UpdatedAt.Format("Jan 02, 2006 15:04:05"))

	// Account Statistics section
	f.SetCellValue(sheetName, "A8", "Account Statistics")
	f.SetCellStyle(sheetName, "A8", "A8", titleStyle)
	f.MergeCell(sheetName, "A8", "B8")

	// Calculate transaction statistics
	stats := calculateTransactionStats(account, transactions)

	// Display transaction statistics
	f.SetCellValue(sheetName, "A9", "Total Deposits:")
	f.SetCellValue(sheetName, "B9", float64(stats.TotalDeposits)/100.0)
	f.SetCellStyle(sheetName, "B9", "B9", styles.CurrencyStyle)

	f.SetCellValue(sheetName, "A10", "Total Withdrawals:")
	f.SetCellValue(sheetName, "B10", float64(stats.TotalWithdrawals)/100.0)
	f.SetCellStyle(sheetName, "B10", "B10", styles.CurrencyStyle)

	f.SetCellValue(sheetName, "A11", "Total Transfers In:")
	f.SetCellValue(sheetName, "B11", float64(stats.TotalTransfersIn)/100.0)
	f.SetCellStyle(sheetName, "B11", "B11", styles.CurrencyStyle)

	f.SetCellValue(sheetName, "A12", "Total Transfers Out:")
	f.SetCellValue(sheetName, "B12", float64(stats.TotalTransfersOut)/100.0)
	f.SetCellStyle(sheetName, "B12", "B12", styles.CurrencyStyle)

	f.SetCellValue(sheetName, "A13", "Number of Deposits:")
	f.SetCellValue(sheetName, "B13", stats.NumDeposits)

	f.SetCellValue(sheetName, "A14", "Number of Withdrawals:")
	f.SetCellValue(sheetName, "B14", stats.NumWithdrawals)

	f.SetCellValue(sheetName, "A15", "Number of Transfers In:")
	f.SetCellValue(sheetName, "B15", stats.NumTransfersIn)

	f.SetCellValue(sheetName, "A16", "Number of Transfers Out:")
	f.SetCellValue(sheetName, "B16", stats.NumTransfersOut)

	f.SetCellValue(sheetName, "A17", "Total Number of Transactions:")
	f.SetCellValue(sheetName, "B17", len(transactions))

	// Add account activity summary (if there are transactions)
	if len(transactions) > 0 {
		f.SetCellValue(sheetName, "A19", "Account Activity Summary")
		f.SetCellStyle(sheetName, "A19", "A19", titleStyle)
		f.MergeCell(sheetName, "A19", "B19")

		var oldestTx, newestTx *data.Transaction
		for i, tx := range transactions {
			if i == 0 || tx.CreatedAt.Before(oldestTx.CreatedAt) {
				oldestTx = tx
			}
			if i == 0 || tx.CreatedAt.After(newestTx.CreatedAt) {
				newestTx = tx
			}
		}

		f.SetCellValue(sheetName, "A20", "First Transaction:")
		f.SetCellValue(sheetName, "B20", oldestTx.CreatedAt.Format("Jan 02, 2006"))

		f.SetCellValue(sheetName, "A21", "Latest Transaction:")
		f.SetCellValue(sheetName, "B21", newestTx.CreatedAt.Format("Jan 02, 2006"))

		daysActive := int(newestTx.CreatedAt.Sub(oldestTx.CreatedAt).Hours()/24) + 1
		f.SetCellValue(sheetName, "A22", "Days Account Active:")
		f.SetCellValue(sheetName, "B22", daysActive)

		txPerDay := float64(len(transactions)) / float64(daysActive)
		f.SetCellValue(sheetName, "A23", "Average Transactions Per Day:")
		f.SetCellValue(sheetName, "B23", txPerDay)
	}

	// Auto-fit columns
	for col := 'A'; col <= 'B'; col++ {
		colName := string(col)
		if err := f.SetColWidth(sheetName, colName, colName, 25); err != nil {
			return fmt.Errorf("failed to set column width: %w", err)
		}
	}

	return nil
}

// TransactionStats holds statistical information about account transactions
type TransactionStats struct {
	TotalDeposits     int
	TotalWithdrawals  int
	TotalTransfersIn  int
	TotalTransfersOut int
	NumDeposits       int
	NumWithdrawals    int
	NumTransfersIn    int
	NumTransfersOut   int
}

// calculateTransactionStats calculates transaction statistics for an account
func calculateTransactionStats(account *data.Account, transactions []*data.Transaction) *TransactionStats {
	stats := &TransactionStats{}

	for _, tx := range transactions {
		switch tx.Type {
		case data.DepositTransaction:
			stats.TotalDeposits += int(tx.Amount)
			stats.NumDeposits++
		case data.WithdrawalTransaction:
			stats.TotalWithdrawals += int(tx.Amount)
			stats.NumWithdrawals++
		case data.TransferTransaction:
			if tx.Recipient == account.ID {
				stats.TotalTransfersIn += int(tx.Amount)
				stats.NumTransfersIn++
			} else {
				stats.TotalTransfersOut += int(tx.Amount)
				stats.NumTransfersOut++
			}
		}
	}

	return stats
}

// CreateTransactionHistorySheet creates the Transaction History sheet in the Excel file
func CreateTransactionHistorySheet(f *excelize.File, account *data.Account, transactions []*data.Transaction, styles *ExcelStyles) error {
	sheetName := TransactionsHistorySheetName
	f.NewSheet(sheetName)

	// Set transaction headers
	headers := []string{"Date", "Type", "Amount", "Balance After", "Details"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, styles.HeaderStyle)
	}

	// Sort transactions chronologically (newest first)
	sortedTx := make([]*data.Transaction, len(transactions))
	copy(sortedTx, transactions)
	sort.Slice(sortedTx, func(i, j int) bool {
		return sortedTx[i].CreatedAt.After(sortedTx[j].CreatedAt)
	})

	// Calculate historical balances
	txWithBalance := calculateHistoricalBalances(account, sortedTx)

	// Add transaction data
	for i, txInfo := range txWithBalance {
		row := i + 2 // Start from row 2 (after header)
		tx := txInfo.Transaction

		// Format transaction type and details
		txType, details := formatTransactionTypeAndDetails(account, tx)

		// Set amount with sign
		var amount float64
		if tx.Type == data.DepositTransaction || (tx.Type == data.TransferTransaction && tx.Recipient == account.ID) {
			amount = float64(tx.Amount) / 100.0
		} else {
			amount = -float64(tx.Amount) / 100.0
		}

		// Add transaction row
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), tx.CreatedAt.Format("Jan 02, 2006 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), txType)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), amount)
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styles.CurrencyStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), float64(txInfo.BalanceAfter)/100.0)
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), styles.CurrencyStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), details)
	}

	// Auto-fit columns
	for col := 'A'; col <= 'E'; col++ {
		colName := string(col)
		width := 20.0
		if col == 'E' {
			width = 55.0 // Make details column wider
		}
		if err := f.SetColWidth(sheetName, colName, colName, width); err != nil {
			return fmt.Errorf("failed to set column width: %w", err)
		}
	}

	return nil
}

// TransactionWithBalance pairs a transaction with its calculated balance after the transaction
type TransactionWithBalance struct {
	Transaction  *data.Transaction
	BalanceAfter int
}

// calculateHistoricalBalances calculates the balance after each transaction
func calculateHistoricalBalances(account *data.Account, transactions []*data.Transaction) []TransactionWithBalance {
	result := make([]TransactionWithBalance, len(transactions))

	// Start with current balance and work backwards through sorted transactions (newest first)
	runningBalance := account.Balance

	for i, tx := range transactions {
		// Calculate balance before this transaction
		var balanceBeforeTx int
		switch tx.Type {
		case data.DepositTransaction:
			balanceBeforeTx = runningBalance - int(tx.Amount)
		case data.WithdrawalTransaction:
			balanceBeforeTx = runningBalance + int(tx.Amount)
		case data.TransferTransaction:
			if tx.Recipient == account.ID {
				balanceBeforeTx = runningBalance - int(tx.Amount)
			} else {
				balanceBeforeTx = runningBalance + int(tx.Amount)
			}
		}

		// Store the result
		result[i] = TransactionWithBalance{
			Transaction:  tx,
			BalanceAfter: runningBalance,
		}

		// Update running balance for next iteration
		runningBalance = balanceBeforeTx
	}

	return result
}

// formatTransactionTypeAndDetails returns formatted type and details for a transaction
func formatTransactionTypeAndDetails(account *data.Account, tx *data.Transaction) (string, string) {
	var txType, details string

	switch tx.Type {
	case data.DepositTransaction:
		txType = "Deposit"
		details = "Deposit to account (ATM transaction)"
	case data.WithdrawalTransaction:
		txType = "Withdrawal"
		details = "Withdrawal from account"
	case data.TransferTransaction:
		if tx.Recipient == account.ID {
			txType = "Transfer In"
			details = fmt.Sprintf("From account: %s", tx.Sender)
		} else {
			txType = "Transfer Out"
			details = fmt.Sprintf("To account: %s", tx.Recipient)
		}
	}

	return txType, details
}
