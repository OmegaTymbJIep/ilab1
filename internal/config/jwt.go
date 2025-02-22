package config

import (
	"fmt"
	"os"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type JWT struct {
	SigningKey jwk.Key `fig:"signing_key,required"`
}

type jwt struct {
	SigningKeyPath string `fig:"signing_key_path,required"`
}

func (c *config) JWT() *JWT {
	return c.jwt.Do(func() interface{} {
		var cfg jwt

		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "jwt")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out jwt: %w", err))
		}

		signingKeyBytes, err := os.ReadFile(cfg.SigningKeyPath)
		if err != nil {
			panic(fmt.Errorf("failed to read JWT signing key: %w", err))
		}

		signingKey, err := jwk.ParseKey(signingKeyBytes, jwk.WithPEM(true))
		if err != nil {
			panic(fmt.Errorf("failed to parse JWT signing key: %w", err))
		}

		return &JWT{
			SigningKey: signingKey,
		}
	}).(*JWT)
}
