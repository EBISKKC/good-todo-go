-- Add is_public column to todos table
ALTER TABLE "todos" ADD COLUMN "is_public" boolean NOT NULL DEFAULT false;

-- Create index for tenant + is_public queries
CREATE INDEX "todo_tenant_id_is_public" ON "todos" ("tenant_id", "is_public");

-- Drop and recreate the view to include is_public
DROP VIEW IF EXISTS "tenant_todo_views";

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
    td."is_public",
    td."due_date",
    td."completed_at",
    td."created_at",
    td."updated_at"
FROM "todos" td
JOIN "users" u ON u."id" = td."user_id"
WHERE td."tenant_id" = current_setting('app.current_tenant_id', true);
