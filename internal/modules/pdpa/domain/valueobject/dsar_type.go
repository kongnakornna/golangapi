package valueobject

type DSARType int8

const (
	DSARTypeAccess          DSARType = 1
	DSARTypeErasure         DSARType = 3
	DSARTypeWithdrawConsent DSARType = 7
)

func (d DSARType) IsValid() bool {
	switch d {
	case DSARTypeAccess, DSARTypeErasure, DSARTypeWithdrawConsent:
		return true
	}
	return false
}

func (d DSARType) String() string {
	switch d {
	case DSARTypeAccess:
		return "ACCESS"
	case DSARTypeErasure:
		return "ERASURE"
	case DSARTypeWithdrawConsent:
		return "WITHDRAW_CONSENT"
	}
	return "UNKNOWN"
}