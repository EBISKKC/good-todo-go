//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
package repository

import (
	"context"

	"good-todo-go/internal/domain/model"
)

type IUserRepository interface {
	FindByID(ctx context.Context, userID string) (*model.User, error)
	FindByIDs(ctx context.Context, userIDs []string) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
}
