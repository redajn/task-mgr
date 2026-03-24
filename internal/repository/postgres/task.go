package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/redajn/task-mgr/internal/domain"
)

type TaskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) GetByID(ctx context.Context, id int64) (domain.Task, error) {
	const query = `
		SELECT id, title, done, created_at, updated_at
		FROM task
		WHERE id = $1 AND deleted_at IS NULL
	`

	var t domain.Task
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, domain.ErrTaskNotFound
		}
		return domain.Task{}, fmt.Errorf("get task by id: %w", err)
	}
	return t, nil
}

func (r *TaskRepo) List(ctx context.Context, filter domain.TaskFilter) ([]domain.Task, error) {
	query := `
		SELECT id, title, done, created_at, updated_at
		FROM tasks
		WHERE deleted_at IS NULL
	`

	args := []any{}
	if filter.Done != nil {
		args = append(args, *filter.Done)
		query += fmt.Sprintf(" AND done = $%d", len(args))
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("lost tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepo) Create(ctx context.Context, input domain.CreateTaskInput) (domain.Task, error) {
	const query = `
		INSERT INTO tasks (title)
		VALUES ($1)
		RETURNING id, title, done, created_at, updated_at
	`

	var t domain.Task
	err := r.db.QueryRow(ctx, query, input.Title).Scan(
		&t.ID, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}
	return t, nil
}

func (r *TaskRepo) Update(ctx context.Context, id int64, input domain.UpdateTaskInput) (domain.Task, error) {
	setClauses := []string{}
	args := []any{}

	if input.Title != nil {
		args = append(args, *input.Title)
		setClauses = append(setClauses, fmt.Sprintf("title - $%d", len(args)))
	}
	if input.Done != nil {
		args = append(args, *input.Done)
		setClauses = append(setClauses, fmt.Sprintf("done = $%d", len(args)))
	}
	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE tasks
		SET %s, update_at = NOW()
		WHERE id = $%d AND deleted_at IS NULL
		RETURNING id, title, done, created_at, updated_at
	`,
		strings.Join(setClauses, ", "),
		len(args),
	)
	var t domain.Task
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&t.ID, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, domain.ErrTaskNotFound
		}
		return domain.Task{}, fmt.Errorf("update task: %w", err)
	}
	return t, nil
}

func (r *TaskRepo) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE tasks SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}
