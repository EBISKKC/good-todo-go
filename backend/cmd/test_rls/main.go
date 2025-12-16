package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"good-todo-go/internal/infrastructure/database"

	_ "github.com/lib/pq"
)

func main() {
	// Use app user (RLS enforced) for testing
	connStr := "postgres://goodtodo_app:app_secret@localhost:5434/goodtodo_dev?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	// Create ent client
	client := database.NewEntClientWithDB(db)
	defer client.Close()

	ctx := context.Background()

	fmt.Println("=== RLS Test ===")
	fmt.Println()

	// 1. Get all tenants and users first (without RLS context)
	fmt.Println("1. Fetching all tenants (admin view)...")
	tenants, err := client.Tenant.Query().All(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch tenants: %v", err)
	}

	if len(tenants) < 2 {
		fmt.Println("ERROR: Need at least 2 tenants to test RLS. Please create more test data.")
		fmt.Printf("Current tenants: %d\n", len(tenants))
		return
	}

	fmt.Printf("Found %d tenants:\n", len(tenants))
	for _, t := range tenants {
		fmt.Printf("  - %s (slug: %s)\n", t.ID, t.Slug)
	}
	fmt.Println()

	tenant1 := tenants[0]
	tenant2 := tenants[1]

	// 2. Test without tenant context (should return empty or all depending on RLS config)
	fmt.Println("2. Testing User query WITHOUT tenant context...")
	users, err := client.User.Query().All(ctx)
	if err != nil {
		fmt.Printf("   Query failed (expected if RLS is strict): %v\n", err)
	} else {
		fmt.Printf("   WARNING: Got %d users WITHOUT setting tenant context!\n", len(users))
		if len(users) > 0 {
			fmt.Println("   RLS might NOT be working correctly!")
			for _, u := range users {
				fmt.Printf("     - %s (tenant: %s)\n", u.Email, u.TenantID)
			}
		}
	}
	fmt.Println()

	// 3. Test with tenant1 context
	fmt.Printf("3. Testing User query WITH tenant1 context (%s)...\n", tenant1.ID)
	tx1, err := database.WithTenantScope(ctx, client, tenant1.ID)
	if err != nil {
		log.Fatalf("Failed to create tenant scope: %v", err)
	}

	users1, err := tx1.User.Query().All(ctx)
	if err != nil {
		fmt.Printf("   Query failed: %v\n", err)
	} else {
		fmt.Printf("   Found %d users in tenant1\n", len(users1))
		for _, u := range users1 {
			if u.TenantID != tenant1.ID {
				fmt.Printf("   ERROR: User %s belongs to tenant %s, not %s!\n", u.Email, u.TenantID, tenant1.ID)
			} else {
				fmt.Printf("   OK: %s (tenant: %s)\n", u.Email, u.TenantID)
			}
		}
	}
	tx1.Rollback()
	fmt.Println()

	// 4. Test with tenant2 context
	fmt.Printf("4. Testing User query WITH tenant2 context (%s)...\n", tenant2.ID)
	tx2, err := database.WithTenantScope(ctx, client, tenant2.ID)
	if err != nil {
		log.Fatalf("Failed to create tenant scope: %v", err)
	}

	users2, err := tx2.User.Query().All(ctx)
	if err != nil {
		fmt.Printf("   Query failed: %v\n", err)
	} else {
		fmt.Printf("   Found %d users in tenant2\n", len(users2))
		for _, u := range users2 {
			if u.TenantID != tenant2.ID {
				fmt.Printf("   ERROR: User %s belongs to tenant %s, not %s!\n", u.Email, u.TenantID, tenant2.ID)
			} else {
				fmt.Printf("   OK: %s (tenant: %s)\n", u.Email, u.TenantID)
			}
		}
	}
	tx2.Rollback()
	fmt.Println()

	// 5. Cross-tenant access test - try to get tenant1's user with tenant2's context
	fmt.Println("5. Cross-tenant access test...")
	if len(users1) > 0 {
		targetUserID := users1[0].ID
		fmt.Printf("   Trying to access user %s (tenant1) with tenant2 context...\n", targetUserID)

		tx3, err := database.WithTenantScope(ctx, client, tenant2.ID)
		if err != nil {
			log.Fatalf("Failed to create tenant scope: %v", err)
		}

		crossUser, err := tx3.User.Get(ctx, targetUserID)
		if err != nil {
			fmt.Printf("   OK: Access denied (error: %v)\n", err)
		} else {
			fmt.Printf("   ERROR: RLS BYPASS! Got user: %s (tenant: %s)\n", crossUser.Email, crossUser.TenantID)
		}
		tx3.Rollback()
	}
	fmt.Println()

	// 6. Test Todos similarly
	fmt.Println("6. Testing Todo query with tenant contexts...")
	tx4, err := database.WithTenantScope(ctx, client, tenant1.ID)
	if err != nil {
		log.Fatalf("Failed to create tenant scope: %v", err)
	}

	todos1, err := tx4.Todo.Query().All(ctx)
	if err != nil {
		fmt.Printf("   Query failed: %v\n", err)
	} else {
		fmt.Printf("   Found %d todos in tenant1\n", len(todos1))
		for _, t := range todos1 {
			if t.TenantID != tenant1.ID {
				fmt.Printf("   ERROR: Todo %s belongs to tenant %s, not %s!\n", t.ID, t.TenantID, tenant1.ID)
			}
		}
	}
	tx4.Rollback()

	fmt.Println()
	fmt.Println("=== RLS Test Complete ===")
}
