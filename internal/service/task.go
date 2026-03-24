package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/redajn/task-mgr/internal/domain"
)

type TaskService struct {
	repo domain.TaskRepository
}

func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (domain.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) ListTask(ctx context.Context, filter domain.TaskFilter) ([]domain.Task, error) {
	return s.repo.List(ctx, filter)
}

func (s *TaskService) CreateTask(ctx context.Context, input domain.CreateTaskInput) (domain.Task, error) {
	if input.Title == "" {
		return domain.Task{}, domain.ErrTaskTitleEmpty
	}

	if len(input.Title) > 255 {
		return domain.Task{}, fmt.Errorf("title roo long: max 255 chars")
	}

	return s.repo.Create(ctx, input)
}

func (s *TaskService) UpdateTask(ctx context.Context, id int64, input domain.UpdateTaskInput) (domain.Task, error) {
	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		return domain.Task{}, domain.ErrTaskTitleEmpty
	}
	return s.repo.Update(ctx, id, input)
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

var _ domain.TaskService = (*TaskService)(nil)
