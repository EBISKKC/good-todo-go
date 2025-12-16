package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// TenantTodoView holds the schema definition for the TenantTodoView entity.
// This is a read-only view for tenant-scoped todo access.
type TenantTodoView struct {
	ent.Schema
}

// Annotations of the TenantTodoView.
func (TenantTodoView) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Define the view using SQL builder with tenant isolation
		entsql.ViewFor("postgres", func(s *sql.Selector) {
			t := sql.Table("todos").As("t")
			u := sql.Table("users").As("u")
			s.From(t).
				Join(u).On(t.C("user_id"), u.C("id")).
				Where(sql.ExprP("t.tenant_id = current_setting('app.current_tenant_id', true)")).
				Select(
					t.C("id"),
					t.C("tenant_id"),
					t.C("user_id"),
					sql.As(u.C("name"), "user_name"),
					sql.As(u.C("email"), "user_email"),
					t.C("title"),
					t.C("description"),
					t.C("completed"),
					t.C("due_date"),
					t.C("completed_at"),
					t.C("created_at"),
					t.C("updated_at"),
				)
		}),
		// Skip migration for view (view is created via migration file)
		entsql.Skip(),
	}
}

// Fields of the TenantTodoView.
func (TenantTodoView) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("tenant_id"),
		field.String("user_id"),
		field.String("user_name"),
		field.String("user_email"),
		field.String("title"),
		field.Text("description"),
		field.Bool("completed"),
		field.Time("due_date").
			Optional().
			Nillable(),
		field.Time("completed_at").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(func() time.Time {
				return time.Now().UTC()
			}),
		field.Time("updated_at").
			Default(func() time.Time {
				return time.Now().UTC()
			}),
	}
}

// Edges of the TenantTodoView.
func (TenantTodoView) Edges() []ent.Edge {
	return nil
}
