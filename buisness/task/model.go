package task

const (
	StatusDone = "done"
)

// Info task information
type Info struct {
	ID          int    `db:"id" json:"id" param:"task_id"`
	UserID      int    `db:"user_id" json:"-"`
	Title       string `db:"title" json:"title" validate:"required"`
	Description string `db:"description" json:"description" validate:"required"`
	Status      string `db:"status" json:"status" validate:"oneof=todo in-progress done"`
}
