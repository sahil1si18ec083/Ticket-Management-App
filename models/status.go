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

func IsValidTransition(from, to Status) bool {
	switch from {
	case StatusNew:
		return to == StatusInProgress || to == StatusWaiting
	case StatusInProgress:
		return to == StatusWaiting || to == StatusResolved
	case StatusWaiting:
		return to == StatusResolved || to == StatusInProgress
	case StatusResolved:
		return to == StatusClosed
	default:
		return false

	}

}
