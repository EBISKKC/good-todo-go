-- Drop views (RLS handles tenant isolation, views are redundant)
DROP VIEW IF EXISTS "tenant_todo_views";
DROP VIEW IF EXISTS "tenant_user_views";
