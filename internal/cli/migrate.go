package cli

import (
	"fmt"

	migrate "github.com/rubenv/sql-migrate"

	"github.com/omegatymbjiep/ilab1/internal/assets"
	"github.com/omegatymbjiep/ilab1/internal/config"
)

var migrations = &migrate.EmbedFileSystemMigrationSource{
	FileSystem: assets.Migrations,
	Root:       "migrations",
}

func MigrateUp(cfg config.Config) error {
	applied, err := migrate.Exec(cfg.DB().RawDB(), "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	cfg.Log().WithField("applied", applied).Info("migrations applied")
	return nil
}

func MigrateDown(cfg config.Config) error {
	applied, err := migrate.Exec(cfg.DB().RawDB(), "postgres", migrations, migrate.Down)
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	cfg.Log().WithField("applied", applied).Info("migrations applied")
	return nil
}
