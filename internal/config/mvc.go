package config

import (
	"fmt"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type MVCConfig struct {
	TemplatesDir string `fig:"templates_dir,required"`
}

func (c *config) MVC() MVCConfig {
	return c.mvc.Do(func() interface{} {
		var cfg MVCConfig

		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "mvc")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out mvc: %w", err))
		}

		return cfg
	}).(MVCConfig)
}
