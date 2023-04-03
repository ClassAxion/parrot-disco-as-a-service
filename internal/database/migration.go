package database

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/mlvzk/gopgs/migrate"
)

//go:embed migrations/*
var migrationsFs embed.FS

func RunMigrations(ctx context.Context, databaseUrl string) error {
	conn, err := pgx.Connect(ctx, databaseUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to db for migrations: %w", err)
	}

	migrations, err := migrate.EmbedFsToMigrations(migrationsFs, "migrations")
	if err != nil {
		return fmt.Errorf("failed to convert migration files to migrations: %w", err)
	}

	return migrate.Migrate(ctx, conn, "public", migrations)
}
