package models

type Status string

const (
	StatusNew        Status = "NEW"
	StatusInProgress Status = "IN_PROGRESS"
	StatusWaiting    Status = "WAITING"
	StatusResolved   Status = "RESOLVED"
	StatusClosed     Status = "CLOSED"
)

var StatusMap = map[string]Status{
	"NEW":         StatusNew,
	"IN_PROGRESS": StatusInProgress,
	"WAITING":     StatusWaiting,
	"RESOLVED":    StatusResolved,
	"CLOSED":      StatusClosed,
}
