package store

import (
	"Todo"
	"context"
	"strings"
)

type TaskService struct {
	db *DB
}

func NewTaskService(db *DB) *TaskService {
	return &TaskService{db}
}

func (t *TaskService) GetTasks(ctx context.Context, filter Todo.TaskFilter) ([]*Todo.Task, error) {
	tx, err := t.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	tasks, err := getTasks(ctx, tx, filter)

	if err != nil {
		return nil, err
	}

	return tasks, nil

}

func (t *TaskService) GetTaskById(ctx context.Context, id int) (*Todo.Task, error) {

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	task, err := getTaskById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *TaskService) CreateTask(ctx context.Context, task *Todo.Task) error {

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := createTask(ctx, tx, task); err != nil {
		return err
	}
	return tx.Commit()
}

func (t *TaskService) UpdateTask(ctx context.Context, id int, upd Todo.TaskUpdate) (*Todo.Task, error) {

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	task, err := updateTask(ctx, tx, id, upd)

	if err != nil {
		return task, err
	}

	return task, tx.Commit()
}

func (t *TaskService) DeleteTask(ctx context.Context, id int) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteTask(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func deleteTask(ctx context.Context, tx *Tx, id int) error {

	if _, err := getTaskById(ctx, tx, id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx,`DELETE FROM tasks WHERE id = ?`, id); err != nil{
		return err
	}
	return nil
}

func updateTask(ctx context.Context, tx *Tx, id int, upd Todo.TaskUpdate) (*Todo.Task, error) {

	task, err := getTaskById(ctx, tx, id)

	if err != nil {
		return nil, err
	}

	if v := upd.Name; v != nil {
		task.Name = *v
	}

	task.UpdatedAt = tx.now

	if err := task.Validate(); err != nil {
		return task, err
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE tasks
		SET name = ?,
			updated_at = ?
		WHERE  id =?
	`,
		task.Name,
		task.UpdatedAt,
		id,
	); err != nil {
		return task, err
	}

	return task, nil
}

func getTasks(ctx context.Context, tx *Tx, filter Todo.TaskFilter) (_ []*Todo.Task, err error) {

	where := []string{"1 = 1"}
	var args []interface{}

	if v := filter.TaskId; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
		    name,
		    created_at,
		    updated_at
		FROM tasks
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Todo.Task, 0)

	for rows.Next() {
		var task Todo.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func getTaskById(ctx context.Context, tx *Tx, id int) (*Todo.Task, error) {
	tasks, err := getTasks(ctx, tx, Todo.TaskFilter{TaskId: &id})
	if err != nil {
		return nil, err
	} else if len(tasks) == 0 {
		return nil, &Todo.Error{Code: Todo.ENOTFOUND, Message: "Task not found."}
	}
	return tasks[0], nil
}

func createTask(ctx context.Context, tx *Tx, task *Todo.Task) error {

	task.CreatedAt = tx.now
	task.UpdatedAt = task.CreatedAt

	if err := task.Validate(); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
			INSERT INTO tasks (name, created_at, updated_at) VALUES (?,?,?)`,
		task.Name, task.CreatedAt, task.UpdatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	task.ID = int(id)
	return nil
}
