package repository

import (
	"context"

	"good-todo-go/internal/domain/model"
	"good-todo-go/internal/domain/repository"
	"good-todo-go/internal/ent"
	"good-todo-go/internal/ent/user"
	"good-todo-go/internal/infrastructure/database"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) repository.IUserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*model.User, error) {
	tx, err := database.TenantScopedTx(ctx, r.client)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	u, err := tx.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return toUserModel(u), nil
}

func (r *UserRepository) FindByIDs(ctx context.Context, userIDs []string) ([]*model.User, error) {
	if len(userIDs) == 0 {
		return []*model.User{}, nil
	}

	tx, err := database.TenantScopedTx(ctx, r.client)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	users, err := tx.User.Query().
		Where(user.IDIn(userIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	result := make([]*model.User, len(users))
	for i, u := range users {
		result[i] = toUserModel(u)
	}
	return result, nil
}

func (r *UserRepository) Update(ctx context.Context, u *model.User) (*model.User, error) {
	tx, err := database.TenantScopedTx(ctx, r.client)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	updated, err := tx.User.UpdateOneID(u.ID).
		SetName(u.Name).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return toUserModel(updated), nil
}
