package task

import (
	"context"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Task hold resources for Task to executing buissness operation
type Task struct {
	db  *sqlx.DB
	log zerolog.Logger
}

// New init Task
func New(db *sqlx.DB, log zerolog.Logger) Task {
	return Task{
		db:  db,
		log: log,
	}
}

// Create creating task in database
func (t Task) Create(ctx context.Context, info Info) (Info, error) {
	q := `INSERT INTO tasks (title, description, status, user_id) VALUES (:title, :description, :status, :user_id) RETURNING id`
	stmt, err := t.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return info, errors.Wrap(err, "Create: preparing name statment")
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &info.ID, info)
	if err != nil {
		return info, errors.Wrap(err, "Create: creating task and querying last inserted id")
	}

	return info, nil
}

// Query querying task from database for user
func (t Task) Query(ctx context.Context, userID int) ([]Info, error) {
	var tasks []Info
	q := "SELECT id, title, description, status FROM tasks WHERE user_id = $1"
	err := t.db.SelectContext(ctx, &tasks, q, userID)
	if err != nil {
		return nil, errors.Wrap(err, "Query: querying task")
	}

	return tasks, nil
}

// QueryByID querying task from database for user with id filter
func (t Task) QueryByID(ctx context.Context, id, userID int) (Info, error) {
	var task Info
	q := "SELECT id, title, description, status FROM tasks WHERE id = $1 AND user_id = $2"
	err := t.db.GetContext(ctx, &task, q, id, userID)
	return task, err
}

// Update updating task in database for user
func (t Task) Update(ctx context.Context, info Info) error {
	q := `UPDATE tasks
				SET title = :title,
				    description = :description,
				    status = :status,
					updated_at = now()
				WHERE id = :id AND user_id = :user_id`
	_, err := t.db.NamedExecContext(ctx, q, info)
	if err != nil {
		return errors.Wrap(err, "Update: executing update query")
	}

	return nil
}

// Delete deleting task for user
func (t Task) Delete(ctx context.Context, id, userID int) error {
	q := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	_, err := t.db.ExecContext(ctx, q, id, userID)
	if err != nil {
		return errors.Wrap(err, "Delete: executing delete task query")
	}

	return nil
}

func (t Task) UpdateStatus(ctx context.Context, ids []int, msgChan chan string, userID int) {
	q := `UPDATE tasks
				SET status = :status,
					updated_at = now()
				WHERE id = :id AND user_id = :user_id`

	var wg sync.WaitGroup
	for _, id := range ids {
		id := id
		wg.Add(1)
		go func() {
			_, err := t.db.NamedExecContext(ctx, q, Info{ID: id, UserID: userID, Status: StatusDone})
			if err != nil {
				msgChan <- fmt.Sprintf("unable to set task number %d as completed.", id)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	close(msgChan)
}
