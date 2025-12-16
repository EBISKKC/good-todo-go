package repository

import (
	"context"

	"good-todo-go/internal/domain/model"
	"good-todo-go/internal/domain/repository"
	"good-todo-go/internal/ent"
	"good-todo-go/internal/ent/todo"
	"good-todo-go/internal/infrastructure/database"
)

type TodoRepository struct {
	client *ent.Client
}

func NewTodoRepository(client *ent.Client) repository.ITodoRepository {
	return &TodoRepository{client: client}
}

// FindByID reads a single todo (RLS handles tenant isolation)
func (r *TodoRepository) FindByID(ctx context.Context, todoID string) (*model.Todo, error) {
	tx, err := database.TenantScopedTx(ctx, r.client)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	t, err := tx.Todo.Get(ctx, todoID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return toTodoModel(t), nil
}

// FindByUserID reads todos (RLS handles tenant isolation)
func (r *TodoRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Todo, error) {
	tx, err := database.TenantScopedTx(ctx, r.client)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	todos, err := tx.Todo.Query().
		Where(todo.UserIDEQ(userID)).
		Order(ent.Desc(todo.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	result := make([]*model.Todo, len(todos))
	for i, t := range todos {
		result[i] = toTodoModel(t)
	}
	return result, nil
}

// CountByUserID counts todos (RLS handles tenant isolation)
func (r *TodoRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	tx, err := database.TenantScopedTx(ctx, r.client)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	count, err := tx.Todo.Query().
		Where(todo.UserIDEQ(userID)).
		Count(ctx)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return count, nil
}

// FindPublic reads public todos from the same tenant (RLS handles tenant isolation)
func (r *TodoRepository) FindPublic(ctx context.Context, limit, offset int) ([]*model.Todo, error) {
	tx, err := database.TenantScopedTx(ctx, r.client)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	todos, err := tx.Todo.Query().
		Where(todo.IsPublicEQ(true)).
		Order(ent.Desc(todo.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	result := make([]*model.Todo, len(todos))
	for i, t := range todos {
		result[i] = toTodoModel(t)
	}
	return result, nil
}

// CountPublic counts public todos in the same tenant (RLS handles tenant isolation)
func (r *TodoRepository) CountPublic(ctx context.Context) (int, error) {
	tx, err := database.TenantScopedTx(ctx, r.client)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	count, err := tx.Todo.Query().
		Where(todo.IsPublicEQ(true)).
		Count(ctx)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return count, nil
}

// Create writes directly to todos table (RLS protected)
func (r *TodoRepository) Create(ctx context.Context, t *model.Todo) (*model.Todo, error) {
	tx, err := database.TenantScopedTx(ctx, r.client)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	builder := tx.Todo.Create().
		SetID(t.ID).
		SetTenantID(t.TenantID).
		SetUserID(t.UserID).
		SetTitle(t.Title).
		SetDescription(t.Description).
		SetCompleted(t.Completed).
		SetIsPublic(t.IsPublic)

	if t.DueDate != nil {
		builder.SetDueDate(*t.DueDate)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return toTodoModel(created), nil
}

// Update writes directly to todos table (RLS protected)
func (r *TodoRepository) Update(ctx context.Context, t *model.Todo) (*model.Todo, error) {
	tx, err := database.TenantScopedTx(ctx, r.client)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	builder := tx.Todo.UpdateOneID(t.ID).
		SetTitle(t.Title).
		SetDescription(t.Description).
		SetCompleted(t.Completed).
		SetIsPublic(t.IsPublic)

	if t.DueDate != nil {
		builder.SetDueDate(*t.DueDate)
	} else {
		builder.ClearDueDate()
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return toTodoModel(updated), nil
}

// Delete writes directly to todos table (RLS protected)
func (r *TodoRepository) Delete(ctx context.Context, todoID string) error {
	tx, err := database.TenantScopedTx(ctx, r.client)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.Todo.DeleteOneID(todoID).Exec(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

// toTodoModel converts ent.Todo to model.Todo
func toTodoModel(t *ent.Todo) *model.Todo {
	return &model.Todo{
		ID:          t.ID,
		UserID:      t.UserID,
		TenantID:    t.TenantID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		IsPublic:    t.IsPublic,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
