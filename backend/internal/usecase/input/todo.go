package input

import "time"

type CreateTodoInput struct {
	UserID      string
	TenantID    string
	Title       string
	Description string
	IsPublic    bool
	DueDate     *time.Time
}

type UpdateTodoInput struct {
	TodoID      string
	UserID      string
	Title       *string
	Description *string
	Completed   *bool
	IsPublic    *bool
	DueDate     *time.Time
}

type GetTodosInput struct {
	UserID string
	Limit  int
	Offset int
}

type GetPublicTodosInput struct {
	Limit  int
	Offset int
}
