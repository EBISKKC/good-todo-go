-- Create application user for RLS enforcement
-- This user does NOT have BYPASSRLS privilege, so RLS policies will be enforced

-- Create app user (if not exists)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'goodtodo_app') THEN
        CREATE USER goodtodo_app WITH PASSWORD 'app_secret';
    END IF;
END
$$;

-- Grant necessary privileges to app user
GRANT CONNECT ON DATABASE goodtodo_dev TO goodtodo_app;
GRANT USAGE ON SCHEMA public TO goodtodo_app;

-- Grant table privileges
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO goodtodo_app;

-- Grant sequence privileges (for auto-generated IDs if any)
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO goodtodo_app;

-- Ensure future tables/sequences also get the grants
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO goodtodo_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO goodtodo_app;

-- Explicitly ensure the user does NOT bypass RLS (should be default, but being explicit)
ALTER USER goodtodo_app NOBYPASSRLS;
