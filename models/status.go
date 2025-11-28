package models

type Status string

const (
	StatusNew        Status = "NEW"
	StatusInProgress Status = "IN_PROGRESS"
	StatusWaiting    Status = "WAITING"
	StatusResolved   Status = "RESOLVED"
	StatusClosed     Status = "CLOSED"
)
