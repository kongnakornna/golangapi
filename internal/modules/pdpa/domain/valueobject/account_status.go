package valueobject

type AccountStatus int8

const (
	AccountActive     AccountStatus = 1
	AccountSuspended  AccountStatus = 2
	AccountTerminated AccountStatus = 4
)

func (a AccountStatus) IsValid() bool {
	switch a {
	case AccountActive, AccountSuspended, AccountTerminated:
		return true
	}
	return false
}

func (a AccountStatus) String() string {
	switch a {
	case AccountActive:
		return "ACTIVE"
	case AccountSuspended:
		return "SUSPENDED"
	case AccountTerminated:
		return "TERMINATED"
	}
	return "UNKNOWN"
}