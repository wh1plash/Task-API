package store

import (
	"context"
	"database/sql"
	"fmt"
	"task/types"
	"time"

	_ "github.com/lib/pq"
)

type TaskStore interface {
	InsertTask(context.Context, *types.Task) (*types.Task, error)
	GetTasks(ctx context.Context, offset, limit int, status string) ([]*types.Task, error)
	GetTask(context.Context, string) (*types.Task, error)
	UpdateTask(context.Context, string, *types.TaskParams) (*types.Task, error)
	DeleteTask(ctx context.Context, id string) (int, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (p *PostgresStore) Init() error {
	return p.createTasksTable()
}

func (p *PostgresStore) createTasksTable() error {
	query := `create table if not exists tasks (
		id serial primary key,
		title varchar(50),
		description varchar(50),
		status varchar(50),
		created_at BIGINT,
		updated_at BIGINT
	)`

	_, err := p.db.Exec(query)
	if err != nil {
		return err
	}

	return err
}

func (p *PostgresStore) GetTask(ctx context.Context, id string) (*types.Task, error) {
	rows, err := p.db.QueryContext(ctx, "select * from tasks where id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	task := &types.Task{}
	if err := rows.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return task, nil
}

func (p *PostgresStore) GetTasks(ctx context.Context, offset, limit int, status string) ([]*types.Task, error) {
	rows, err := p.db.QueryContext(ctx, `
	SELECT id, title, description, status, created_at, updated_at
	FROM tasks
	WHERE ($1 = '' OR status = $1)
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*types.Task{}
	for rows.Next() {
		task := new(types.Task)
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return nil, sql.ErrNoRows
	}

	return tasks, nil

}

func (p *PostgresStore) DeleteTask(ctx context.Context, id string) (int, error) {
	query := `Delete from tasks where id=$1 RETURNING id`

	var deletedID int
	err := p.db.QueryRowContext(ctx, query, id).Scan(&deletedID)
	if err != nil {
		return 0, sql.ErrNoRows
	}

	return deletedID, nil
}

func (p *PostgresStore) UpdateTask(ctx context.Context, id string, params *types.TaskParams) (*types.Task, error) {
	query := fmt.Sprintf(`Update tasks
			SET %s
			WHERE id=$1
			returning id, title, description, status, created_at, updated_at
	`, fmt.Sprintf("title='%s', description='%s', status='%s', updated_at=%d",
		params.Title,
		params.Description,
		params.Status,
		uint64(time.Now().UnixNano())))

	updTask := types.Task{}
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&updTask.ID,
		&updTask.Title,
		&updTask.Description,
		&updTask.Status,
		&updTask.CreatedAt,
		&updTask.UpdatedAt,
	)

	if err != nil {
		return nil, sql.ErrNoRows
	}

	return &updTask, nil
}

func (p *PostgresStore) InsertTask(ctx context.Context, task *types.Task) (*types.Task, error) {
	query := `insert into tasks
	(title, description, status, created_at, updated_at)
	values($1, $2, $3, $4, $5)
	returning id, title, description, status, created_at, updated_at`

	insTask := &types.Task{}
	err := p.db.QueryRowContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(
		&insTask.ID,
		&insTask.Title,
		&insTask.Description,
		&insTask.Status,
		&insTask.CreatedAt,
		&insTask.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return insTask, nil
}
