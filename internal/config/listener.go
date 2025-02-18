package config

import (
	"fmt"
	"net"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type listenerConfig struct {
	Addr string `fig:"addr,required"`
}

func (c *config) Listener() net.Listener {
	return c.listener.Do(func() interface{} {
		var cfg listenerConfig

		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "listener")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out listener: %w", err))
		}

		listener, err := net.Listen("tcp", cfg.Addr)
		if err != nil {
			panic(fmt.Errorf("failed to init http listener: %w", err))
		}

		return listener
	}).(net.Listener)
}
