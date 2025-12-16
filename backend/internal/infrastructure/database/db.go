package database

import (
	"context"
	"database/sql"
	"fmt"

	"good-todo-go/internal/ent"
	"good-todo-go/internal/infrastructure/environment"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

// NewEntClient creates an ent.Client for database operations
// Note: All database access should use TenantScopedTx for RLS enforcement
func NewEntClient(cfg *environment.Config) (*ent.Client, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	// Open sql.DB
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	// Create ent client from sql.DB
	drv := entsql.OpenDB("postgres", db)
	client := ent.NewClient(ent.Driver(drv))

	// デバッグモードを有効化（開発時のみ）
	if cfg.AppEnv == "local" {
		client = client.Debug()
	}

	return client, nil
}

func CloseEntClient(client *ent.Client) error {
	if client != nil {
		return client.Close()
	}
	return nil
}

func RunMigrations(ctx context.Context, client *ent.Client) error {
	return client.Schema.Create(ctx)
}
