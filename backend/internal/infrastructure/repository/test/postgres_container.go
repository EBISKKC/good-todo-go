package test

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

type PostgresContainer struct {
	Container *postgres.PostgresContainer
	DB        *sql.DB
	DSN       string
}

func NewPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	dbName := "test_db"
	dbUser := "test_user"
	dbPassword := "test_password"

	container, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresContainer{
		Container: container,
		DB:        db,
		DSN:       connStr,
	}, nil
}

func (pc *PostgresContainer) Close(ctx context.Context) error {
	if pc.DB != nil {
		pc.DB.Close()
	}
	if pc.Container != nil {
		return pc.Container.Terminate(ctx)
	}
	return nil
}

// SetupRLS sets up RLS policies and app user for tenant isolation testing
func (pc *PostgresContainer) SetupRLS(ctx context.Context) error {
	queries := []string{
		// Enable RLS on users table
		`ALTER TABLE users ENABLE ROW LEVEL SECURITY`,
		`ALTER TABLE users FORCE ROW LEVEL SECURITY`,

		// Enable RLS on todos table
		`ALTER TABLE todos ENABLE ROW LEVEL SECURITY`,
		`ALTER TABLE todos FORCE ROW LEVEL SECURITY`,

		// Drop existing policies if they exist
		`DROP POLICY IF EXISTS users_tenant_isolation ON users`,
		`DROP POLICY IF EXISTS todos_tenant_isolation ON todos`,

		// Create RLS policies
		`CREATE POLICY users_tenant_isolation ON users
			FOR ALL
			USING (tenant_id = current_setting('app.current_tenant_id', true))
			WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true))`,

		`CREATE POLICY todos_tenant_isolation ON todos
			FOR ALL
			USING (tenant_id = current_setting('app.current_tenant_id', true))
			WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true))`,

		// Create app user for RLS testing
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'goodtodo_app') THEN
				CREATE USER goodtodo_app WITH PASSWORD 'app_secret';
			END IF;
		END $$`,

		// Grant privileges to app user
		`GRANT USAGE ON SCHEMA public TO goodtodo_app`,
		`GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO goodtodo_app`,
		`GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO goodtodo_app`,
		`ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO goodtodo_app`,
		`ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO goodtodo_app`,
		`ALTER USER goodtodo_app NOBYPASSRLS`,
	}

	for _, q := range queries {
		if _, err := pc.DB.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", q, err)
		}
	}

	return nil
}

// GetAppUserDSN returns DSN for the app user (RLS enforced)
func (pc *PostgresContainer) GetAppUserDSN(ctx context.Context) (string, error) {
	host, err := pc.Container.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := pc.Container.MappedPort(ctx, "5432")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("postgres://goodtodo_app:app_secret@%s:%s/test_db?sslmode=disable", host, port.Port()), nil
}
