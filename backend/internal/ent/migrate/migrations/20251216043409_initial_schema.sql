-- Create "tenants" table
CREATE TABLE "tenants" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "slug" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "tenants_slug_key" to table: "tenants"
CREATE UNIQUE INDEX "tenants_slug_key" ON "tenants" ("slug");
-- Create "users" table
CREATE TABLE "users" (
  "id" character varying NOT NULL,
  "email" character varying NOT NULL,
  "password_hash" character varying NOT NULL,
  "name" character varying NOT NULL DEFAULT '',
  "role" character varying NOT NULL DEFAULT 'member',
  "email_verified" boolean NOT NULL DEFAULT false,
  "verification_token" character varying NULL,
  "verification_token_expires_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "tenant_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_tenants_users" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "user_tenant_id" to table: "users"
CREATE INDEX "user_tenant_id" ON "users" ("tenant_id");
-- Create index "user_tenant_id_email" to table: "users"
CREATE UNIQUE INDEX "user_tenant_id_email" ON "users" ("tenant_id", "email");
-- Create "todos" table
CREATE TABLE "todos" (
  "id" character varying NOT NULL,
  "tenant_id" character varying NOT NULL,
  "title" character varying NOT NULL,
  "description" text NULL DEFAULT '',
  "completed" boolean NOT NULL DEFAULT false,
  "due_date" timestamptz NULL,
  "completed_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "todos_users_todos" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "todo_tenant_id" to table: "todos"
CREATE INDEX "todo_tenant_id" ON "todos" ("tenant_id");
-- Create index "todo_tenant_id_user_id" to table: "todos"
CREATE INDEX "todo_tenant_id_user_id" ON "todos" ("tenant_id", "user_id");
-- Create index "todo_user_id" to table: "todos"
CREATE INDEX "todo_user_id" ON "todos" ("user_id");
