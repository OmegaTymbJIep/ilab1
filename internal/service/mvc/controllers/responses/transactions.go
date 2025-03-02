package responses

type TransactionResult struct {
	NewBalance int `json:"new_balance"`
}

func NewTransactionResult(balance int) *TransactionResult {
	return &TransactionResult{
		NewBalance: balance,
	}
}
