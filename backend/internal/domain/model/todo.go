package model

import "time"

type Todo struct {
	ID          string
	UserID      string
	TenantID    string
	Title       string
	Description string
	Completed   bool
	IsPublic    bool
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
