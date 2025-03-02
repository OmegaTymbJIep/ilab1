package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type ATM struct {
	PublicKey *ecdsa.PublicKey
}

type atm struct {
	PublicKeyPath string `fig:"public_key_path,required"`
}

func (c *config) ATM() *ATM {
	return c.atm.Do(func() interface{} {
		var cfg atm

		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "atm")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out jwt: %w", err))
		}

		pemBytes, err := os.ReadFile(cfg.PublicKeyPath)
		if err != nil {
			panic(fmt.Errorf("failed to read PEM file: %w", err))
		}

		block, _ := pem.Decode(pemBytes)
		if block == nil {
			panic(fmt.Errorf("failed to decode PEM block"))
		}

		genericPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			panic(fmt.Errorf("failed to parse public key: %w", err))
		}

		publicKey, ok := genericPublicKey.(*ecdsa.PublicKey)
		if !ok {
			panic(fmt.Errorf("public key is not ECDSA"))
		}

		return &ATM{
			PublicKey: publicKey,
		}
	}).(*ATM)
}
