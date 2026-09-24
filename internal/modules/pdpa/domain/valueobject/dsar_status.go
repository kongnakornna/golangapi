package valueobject

type DSARStatus int8

const (
	DSARStatusPending    DSARStatus = 1
	DSARStatusProcessing DSARStatus = 3
	DSARStatusCompleted  DSARStatus = 4
	DSARStatusRejected   DSARStatus = 5
)

func (s DSARStatus) IsValid() bool {
	switch s {
	case DSARStatusPending, DSARStatusProcessing, DSARStatusCompleted, DSARStatusRejected:
		return true
	}
	return false
}

func (s DSARStatus) String() string {
	switch s {
	case DSARStatusPending:
		return "PENDING"
	case DSARStatusProcessing:
		return "PROCESSING"
	case DSARStatusCompleted:
		return "COMPLETED"
	case DSARStatusRejected:
		return "REJECTED"
	}
	return "UNKNOWN"
}