package Todo

import (
	"context"
	"time"
)

type Task struct {
	ID		int    `json:"id"`
	Name	string `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t *Task) Validate() error {
	if t.Name == "" {
		return Errorf(EINVALID, "Task name required")
	}

	return nil
}

type TaskService interface {
	//Retrieve Single Task by Id
	GetTaskById(ctx context.Context, id int) (*Task, error)
	//Retrieve all tasks
	GetTasks(ctx context.Context, filter TaskFilter) ([]*Task, error)
	//Create New Task
	CreateTask(ctx context.Context, task *Task) error
	//Update Task
	UpdateTask(ctx context.Context, id int, upd TaskUpdate) (*Task, error)
}

type TaskFilter struct {
	TaskId *int `json:"taskId"`
}

type TaskUpdate struct {
	Name *string `json:"name"`
}