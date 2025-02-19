package models

import "github.com/omegatymbjiep/ilab1/internal/data"

type Auth struct {
	db data.MainQ
}

func NewAuth(db data.MainQ) *Auth {
	return &Auth{db: db}
}
