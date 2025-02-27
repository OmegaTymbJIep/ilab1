package views

import "github.com/omegatymbjiep/ilab1/internal/data"

type AccountsList struct {
	Accounts []*data.Account
}

type Account struct {
	Account  *data.Account
	Transactions []*data.Transaction
}
