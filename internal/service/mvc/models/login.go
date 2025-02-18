package models

import "github.com/omegatymbjiep/ilab1/internal/data"

type Login struct {
	db data.MainQ
}

func NewLogin(db data.MainQ) *Login {
	return &Login{db: db}
}
