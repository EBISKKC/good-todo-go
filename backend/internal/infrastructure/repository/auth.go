package repository

import (
	"context"
	"fmt"

	"good-todo-go/internal/domain/model"
	"good-todo-go/internal/domain/repository"
	"good-todo-go/internal/ent"
	"good-todo-go/internal/ent/tenant"
	"good-todo-go/internal/ent/user"
)

type AuthRepository struct {
	client *ent.Client
}

func NewAuthRepository(client *ent.Client) repository.IAuthRepository {
	return &AuthRepository{client: client}
}

func (r *AuthRepository) FindTenantBySlug(ctx context.Context, slug string) (*model.Tenant, error) {
	t, err := r.client.Tenant.Query().
		Where(tenant.SlugEQ(slug)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return toTenantModel(t), nil
}

func (r *AuthRepository) CreateTenant(ctx context.Context, t *model.Tenant) (*model.Tenant, error) {
	created, err := r.client.Tenant.Create().
		SetID(t.ID).
		SetName(t.Name).
		SetSlug(t.Slug).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toTenantModel(created), nil
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, tenantID, email string) (*model.User, error) {
	// Start transaction with tenant context for RLS
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Set tenant context for RLS
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET LOCAL app.current_tenant_id = '%s'", tenantID)); err != nil {
		return nil, fmt.Errorf("failed to set tenant context: %w", err)
	}

	u, err := tx.User.Query().
		Where(
			user.TenantIDEQ(tenantID),
			user.EmailEQ(email),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return toUserModel(u), nil
}

func (r *AuthRepository) FindUserByID(ctx context.Context, tenantID, userID string) (*model.User, error) {
	// Start transaction with tenant context for RLS
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Set tenant context for RLS
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET LOCAL app.current_tenant_id = '%s'", tenantID)); err != nil {
		return nil, fmt.Errorf("failed to set tenant context: %w", err)
	}

	u, err := tx.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return toUserModel(u), nil
}

func (r *AuthRepository) FindUserByVerificationToken(ctx context.Context, token string) (*model.User, error) {
	u, err := r.client.User.Query().
		Where(user.VerificationTokenEQ(token)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return toUserModel(u), nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	// Start transaction with tenant context for RLS
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Set tenant context for RLS
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET LOCAL app.current_tenant_id = '%s'", u.TenantID)); err != nil {
		return nil, fmt.Errorf("failed to set tenant context: %w", err)
	}

	builder := tx.User.Create().
		SetID(u.ID).
		SetTenantID(u.TenantID).
		SetEmail(u.Email).
		SetPasswordHash(u.PasswordHash).
		SetName(u.Name).
		SetRole(user.Role(u.Role)).
		SetEmailVerified(u.EmailVerified)

	if u.VerificationToken != nil {
		builder.SetVerificationToken(*u.VerificationToken)
	}
	if u.VerificationTokenExpiresAt != nil {
		builder.SetVerificationTokenExpiresAt(*u.VerificationTokenExpiresAt)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return toUserModel(created), nil
}

func (r *AuthRepository) UpdateUser(ctx context.Context, u *model.User) (*model.User, error) {
	// Start transaction with tenant context for RLS
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Set tenant context for RLS
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET LOCAL app.current_tenant_id = '%s'", u.TenantID)); err != nil {
		return nil, fmt.Errorf("failed to set tenant context: %w", err)
	}

	builder := tx.User.UpdateOneID(u.ID).
		SetName(u.Name).
		SetEmailVerified(u.EmailVerified)

	if u.VerificationToken != nil {
		builder.SetVerificationToken(*u.VerificationToken)
	} else {
		builder.ClearVerificationToken()
	}

	if u.VerificationTokenExpiresAt != nil {
		builder.SetVerificationTokenExpiresAt(*u.VerificationTokenExpiresAt)
	} else {
		builder.ClearVerificationTokenExpiresAt()
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return toUserModel(updated), nil
}

func toTenantModel(t *ent.Tenant) *model.Tenant {
	return &model.Tenant{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func toUserModel(u *ent.User) *model.User {
	return &model.User{
		ID:                         u.ID,
		TenantID:                   u.TenantID,
		Email:                      u.Email,
		PasswordHash:               u.PasswordHash,
		Name:                       u.Name,
		Role:                       string(u.Role),
		EmailVerified:              u.EmailVerified,
		VerificationToken:          u.VerificationToken,
		VerificationTokenExpiresAt: u.VerificationTokenExpiresAt,
		CreatedAt:                  u.CreatedAt,
		UpdatedAt:                  u.UpdatedAt,
	}
}
