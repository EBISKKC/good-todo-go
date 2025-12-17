package common

import (
	"context"
	"testing"
	"time"

	"good-todo-go/internal/ent"
	"good-todo-go/internal/ent/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// CreateTenant creates a tenant for testing using ent builder
func CreateTenant(t *testing.T, client *ent.Client, builder *ent.TenantCreate) *ent.Tenant {
	t.Helper()
	ctx := context.Background()

	tenant, err := builder.Save(ctx)
	require.NoError(t, err)
	return tenant
}

// DefaultTenantBuilder returns default tenant builder for testing
func DefaultTenantBuilder(client *ent.Client, id string) *ent.TenantCreate {
	if id == "" {
		id = uuid.New().String()
	}
	return client.Tenant.Create().
		SetID(id).
		SetName("Test Tenant").
		SetSlug("test-tenant-" + id[:8])
}

// CreateUser creates a user for testing using ent builder
func CreateUser(t *testing.T, client *ent.Client, builder *ent.UserCreate) *ent.User {
	t.Helper()
	ctx := context.Background()

	u, err := builder.Save(ctx)
	require.NoError(t, err)
	return u
}

// DefaultUserBuilder returns default user builder for testing
func DefaultUserBuilder(client *ent.Client, id, tenantID string) *ent.UserCreate {
	if id == "" {
		id = uuid.New().String()
	}
	return client.User.Create().
		SetID(id).
		SetTenantID(tenantID).
		SetEmail(id[:8] + "@example.com").
		SetPasswordHash("$2a$10$dummy_hash_for_testing").
		SetName("Test User").
		SetRole(user.RoleMember).
		SetEmailVerified(true)
}

// CreateTodo creates a todo for testing using ent builder
func CreateTodo(t *testing.T, client *ent.Client, builder *ent.TodoCreate) *ent.Todo {
	t.Helper()
	ctx := context.Background()

	todo, err := builder.Save(ctx)
	require.NoError(t, err)
	return todo
}

// DefaultTodoBuilder returns default todo builder for testing
func DefaultTodoBuilder(client *ent.Client, id, tenantID, userID string) *ent.TodoCreate {
	if id == "" {
		id = uuid.New().String()
	}
	return client.Todo.Create().
		SetID(id).
		SetTenantID(tenantID).
		SetUserID(userID).
		SetTitle("Test Todo").
		SetDescription("Test Description").
		SetCompleted(false).
		SetIsPublic(false)
}

// PublicTodoBuilder returns todo builder for public todo
func PublicTodoBuilder(client *ent.Client, id, tenantID, userID string) *ent.TodoCreate {
	return DefaultTodoBuilder(client, id, tenantID, userID).
		SetIsPublic(true)
}

// CompletedTodoBuilder returns todo builder for completed todo
func CompletedTodoBuilder(client *ent.Client, id, tenantID, userID string) *ent.TodoCreate {
	now := time.Now()
	return DefaultTodoBuilder(client, id, tenantID, userID).
		SetCompleted(true).
		SetCompletedAt(now)
}

// TestDataSet represents a complete set of test data for multi-tenant testing
type TestDataSet struct {
	Tenant1 *ent.Tenant
	Tenant2 *ent.Tenant
	User1   *ent.User // belongs to Tenant1
	User2   *ent.User // belongs to Tenant1
	User3   *ent.User // belongs to Tenant2
	Todo1   *ent.Todo // belongs to User1 (Tenant1), private
	Todo2   *ent.Todo // belongs to User1 (Tenant1), public
	Todo3   *ent.Todo // belongs to User2 (Tenant1), private
	Todo4   *ent.Todo // belongs to User3 (Tenant2), private
	Todo5   *ent.Todo // belongs to User3 (Tenant2), public
}

// CreateTestDataSet creates a complete set of test data for RLS testing
func CreateTestDataSet(t *testing.T, client *ent.Client) *TestDataSet {
	t.Helper()

	// Create tenants
	tenant1 := CreateTenant(t, client,
		DefaultTenantBuilder(client, "").
			SetName("Tenant One").
			SetSlug("tenant-one"))

	tenant2 := CreateTenant(t, client,
		DefaultTenantBuilder(client, "").
			SetName("Tenant Two").
			SetSlug("tenant-two"))

	// Create users
	user1 := CreateUser(t, client,
		DefaultUserBuilder(client, "", tenant1.ID).
			SetEmail("user1@tenant1.com").
			SetName("User One"))

	user2 := CreateUser(t, client,
		DefaultUserBuilder(client, "", tenant1.ID).
			SetEmail("user2@tenant1.com").
			SetName("User Two"))

	user3 := CreateUser(t, client,
		DefaultUserBuilder(client, "", tenant2.ID).
			SetEmail("user3@tenant2.com").
			SetName("User Three"))

	// Create todos
	todo1 := CreateTodo(t, client,
		DefaultTodoBuilder(client, "", tenant1.ID, user1.ID).
			SetTitle("Todo 1 (Tenant1, User1, Private)").
			SetDescription("Private todo for user1").
			SetIsPublic(false))

	todo2 := CreateTodo(t, client,
		DefaultTodoBuilder(client, "", tenant1.ID, user1.ID).
			SetTitle("Todo 2 (Tenant1, User1, Public)").
			SetDescription("Public todo for user1").
			SetIsPublic(true))

	todo3 := CreateTodo(t, client,
		DefaultTodoBuilder(client, "", tenant1.ID, user2.ID).
			SetTitle("Todo 3 (Tenant1, User2, Private)").
			SetDescription("Private todo for user2").
			SetIsPublic(false))

	todo4 := CreateTodo(t, client,
		DefaultTodoBuilder(client, "", tenant2.ID, user3.ID).
			SetTitle("Todo 4 (Tenant2, User3, Private)").
			SetDescription("Private todo for user3").
			SetIsPublic(false))

	todo5 := CreateTodo(t, client,
		DefaultTodoBuilder(client, "", tenant2.ID, user3.ID).
			SetTitle("Todo 5 (Tenant2, User3, Public)").
			SetDescription("Public todo for user3").
			SetIsPublic(true))

	return &TestDataSet{
		Tenant1: tenant1,
		Tenant2: tenant2,
		User1:   user1,
		User2:   user2,
		User3:   user3,
		Todo1:   todo1,
		Todo2:   todo2,
		Todo3:   todo3,
		Todo4:   todo4,
		Todo5:   todo5,
	}
}
