package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrTaskTitleEmpty = errors.New("task title cannot be empty")
	ErrTaskForbidden  = errors.New("access denied")
)

type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTaskInput struct {
	UserID int64  `json:"-"`
	Title  string `json:"title"`
}

type UpdateTaskInput struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

type TaskFilter struct {
	UserID int64
	Done   *bool
}
type TaskRepository interface {
	GetByID(ctx context.Context, id int64) (Task, error)
	List(ctx context.Context, filter TaskFilter) ([]Task, error)
	Create(ctx context.Context, input CreateTaskInput) (Task, error)
	Update(ctx context.Context, id int64, input UpdateTaskInput) (Task, error)
	Delete(ctx context.Context, id int64) error
}

type TaskService interface {
	GetTask(ctx context.Context, userID, taskID int64) (Task, error)
	ListTask(ctx context.Context, filter TaskFilter) ([]Task, error)
	CreateTask(ctx context.Context, input CreateTaskInput) (Task, error)
	UpdateTask(ctx context.Context, userID, taskID int64, input UpdateTaskInput) (Task, error)
	DeleteTask(ctx context.Context, userID, taskID int64) error
}
