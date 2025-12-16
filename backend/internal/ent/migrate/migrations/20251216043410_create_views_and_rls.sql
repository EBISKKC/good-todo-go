-- Create "tenant_user_views" view with tenant isolation
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
JOIN "tenants" t ON t."id" = u."tenant_id"
WHERE u."tenant_id" = current_setting('app.current_tenant_id', true);

-- Create "tenant_todo_views" view with tenant isolation
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
JOIN "users" u ON u."id" = td."user_id"
WHERE td."tenant_id" = current_setting('app.current_tenant_id', true);

-- Enable RLS on users table
ALTER TABLE "users" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "users" FORCE ROW LEVEL SECURITY;

-- Enable RLS on todos table
ALTER TABLE "todos" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "todos" FORCE ROW LEVEL SECURITY;

-- RLS Policy for users (ALL operations)
CREATE POLICY "users_tenant_isolation" ON "users"
    FOR ALL
    USING ("tenant_id" = current_setting('app.current_tenant_id', true))
    WITH CHECK ("tenant_id" = current_setting('app.current_tenant_id', true));

-- RLS Policy for todos (ALL operations)
CREATE POLICY "todos_tenant_isolation" ON "todos"
    FOR ALL
    USING ("tenant_id" = current_setting('app.current_tenant_id', true))
    WITH CHECK ("tenant_id" = current_setting('app.current_tenant_id', true));
