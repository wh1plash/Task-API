package types

import "time"

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   uint64 `json:"createdAt"`
	UpdatedAt   uint64 `json:"updatedAt"`
}

type TaskParams struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
}

func NewTaskFromParams(params TaskParams) (*Task, error) {
	return &Task{
		Title:       params.Title,
		Description: params.Description,
		Status:      params.Status,
		CreatedAt:   uint64(time.Now().UnixNano()),
		UpdatedAt:   uint64(time.Now().UnixNano()),
	}, nil
}
