-- ============================================================================
-- Views for Read Operations (Tenant Isolation)
-- ============================================================================

-- TenantUserView: Read-only view for users with tenant info
CREATE VIEW "tenant_user_views" AS
SELECT
    u."id",
    u."tenant_id",
    t."name" AS "tenant_name",
    t."slug" AS "tenant_slug",
    u."email",
    u."name",
    u."role",
    u."email_verified",
    u."created_at",
    u."updated_at"
FROM "users" u
INNER JOIN "tenants" t ON u."tenant_id" = t."id";

-- TenantTodoView: Read-only view for todos with user info
CREATE VIEW "tenant_todo_views" AS
SELECT
    td."id",
    td."tenant_id",
    td."user_id",
    u."name" AS "user_name",
    u."email" AS "user_email",
    td."title",
    td."description",
    td."completed",
    td."due_date",
    td."completed_at",
    td."created_at",
    td."updated_at"
FROM "todos" td
INNER JOIN "users" u ON td."user_id" = u."id";

-- ============================================================================
-- Row Level Security (RLS) for Write Operations
-- ============================================================================

-- Enable RLS on users table
ALTER TABLE "users" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "users" FORCE ROW LEVEL SECURITY;

-- Enable RLS on todos table
ALTER TABLE "todos" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "todos" FORCE ROW LEVEL SECURITY;

-- ============================================================================
-- RLS Policies for users table
-- ============================================================================

CREATE POLICY "users_tenant_isolation_select" ON "users"
    FOR SELECT
    USING ("tenant_id" = current_setting('app.current_tenant_id', true));

CREATE POLICY "users_tenant_isolation_insert" ON "users"
    FOR INSERT
    WITH CHECK ("tenant_id" = current_setting('app.current_tenant_id', true));

CREATE POLICY "users_tenant_isolation_update" ON "users"
    FOR UPDATE
    USING ("tenant_id" = current_setting('app.current_tenant_id', true))
    WITH CHECK ("tenant_id" = current_setting('app.current_tenant_id', true));

CREATE POLICY "users_tenant_isolation_delete" ON "users"
    FOR DELETE
    USING ("tenant_id" = current_setting('app.current_tenant_id', true));

-- ============================================================================
-- RLS Policies for todos table
-- ============================================================================

CREATE POLICY "todos_tenant_isolation_select" ON "todos"
    FOR SELECT
    USING ("tenant_id" = current_setting('app.current_tenant_id', true));

CREATE POLICY "todos_tenant_isolation_insert" ON "todos"
    FOR INSERT
    WITH CHECK ("tenant_id" = current_setting('app.current_tenant_id', true));

CREATE POLICY "todos_tenant_isolation_update" ON "todos"
    FOR UPDATE
    USING ("tenant_id" = current_setting('app.current_tenant_id', true))
    WITH CHECK ("tenant_id" = current_setting('app.current_tenant_id', true));

CREATE POLICY "todos_tenant_isolation_delete" ON "todos"
    FOR DELETE
    USING ("tenant_id" = current_setting('app.current_tenant_id', true));
