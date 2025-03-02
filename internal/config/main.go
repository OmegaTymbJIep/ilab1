package config

import (
	"net"

	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser

	MVC() *MVC
	JWT() *JWT
	ATM() *ATM
	Listener() net.Listener
}

type config struct {
	comfig.Logger
	pgdb.Databaser

	listener comfig.Once
	mvc      comfig.Once
	jwt      comfig.Once
	atm      comfig.Once

	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:    getter,
		Logger:    comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Databaser: pgdb.NewDatabaser(getter),
	}
}
