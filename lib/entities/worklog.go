package entities

import "github.com/google/uuid"

type Worklog struct {
	Id          uuid.UUID
	Content     string
	Type        WorklogType
	Incident_Id uuid.UUID
	Writer      User
	CreatedAt   int64
}

type WorklogType string

const (
	WT_General WorklogType = "General"
	WT_Comment WorklogType = "Comment"
)

func (w *Worklog) ScanTo(scan ScanFunc) error {
	return scan(
		&w.Id,
		&w.Content,
		&w.Type,
		&w.Incident_Id,
		&w.Writer,
		&w.CreatedAt,
	)
}
