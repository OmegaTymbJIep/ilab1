package config

import (
	"fmt"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type MVC struct {
	TemplatesDir string `fig:"templates_dir,required"`
}

func (c *config) MVC() *MVC {
	return c.mvc.Do(func() interface{} {
		var cfg MVC

		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "mvc")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out mvc: %w", err))
		}

		return &cfg
	}).(*MVC)
}
