package entities

import "github.com/google/uuid"

type Task struct {
	Id          uuid.UUID     `json:"id"`
	Title       string        `json:"name"`
	Description string        `json:"description"`
	State       TaskStatus    `json:"state"`
	Priority    Priority      `json:"priority"`
	Owner       User          `json:"owner"`
	IncidentID  uuid.UUID     `json:"incidentid"`
	Comments    []TaskComment `json:"comments"`
}

type TaskStatus string

func NewTask(title, description string, owner User, incident_id uuid.UUID) *Task {
	return &Task{
		Id:          uuid.New(),
		Title:       title,
		Description: description,
		Owner:       owner,
		State:       Backlog,
		Priority:    Low,
		IncidentID:  incident_id,
		Comments:    []TaskComment{},
	}
}

func (t *Task) WithPriority(prio Priority) *Task {
	t.Priority = prio
	return t
}

const (
	Backlog TaskStatus = "Backlog"
	Doing   TaskStatus = "Doing"
	Review  TaskStatus = "Review"
	Done    TaskStatus = "Done"
)

func (t *Task) ScanTo(scan ScanFunc) error {
	return scan(
		&t.Id,
		&t.Title,
		&t.Description,
		&t.State,
		&t.Priority,
		&t.Owner.Id,
		&t.Owner.Firstname,
		&t.Owner.Lastname,
		&t.Owner.Email,
		&t.Owner.Fullname,
		&t.Owner.CreatedAt,
	)
}

type TaskComment struct {
	Title        string
	Content      string
	CreatedAt    uint
	LastModified uint
	Writer       User
	Task         uuid.UUID
}

func (t *TaskComment) ScanTo(scan ScanFunc) error {
	return scan(
		&t.Title,
		&t.Content,
		&t.CreatedAt,
		&t.LastModified,
		&t.Writer.Id,
		&t.Writer.Firstname,
		&t.Writer.Lastname,
		&t.Writer.Email,
		&t.Writer.Fullname,
		&t.Writer.CreatedAt,
	)
}
