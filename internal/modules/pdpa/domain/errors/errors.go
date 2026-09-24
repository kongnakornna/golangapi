package domainerrors

import "errors"

var (
	ErrConsentNotFound              = errors.New("consent record not found")
	ErrConsentAlreadyRevoked        = errors.New("consent already revoked")
	ErrConsentExpired               = errors.New("consent expired")
	ErrDSARNotFound                 = errors.New("DSAR request not found")
	ErrDSARAlreadyProcessed         = errors.New("DSAR request already processed")
	ErrInvalidPurpose               = errors.New("invalid consent purpose")
	ErrOTPExpired                   = errors.New("OTP expired or invalid")
	ErrRateLimitExceeded            = errors.New("too many DSAR requests, please try again later")
	ErrRevokeNotAllowed             = errors.New("cannot revoke non-granted consent")
	ErrAccountNotFound              = errors.New("user account status not found")
	ErrAccountNotSuspended          = errors.New("account is not suspended")
	ErrAccountNotTerminated         = errors.New("account is not terminated")
	ErrDeletionNotConfirmed         = errors.New("deletion not confirmed for this account")
	ErrImmediateDeletionNotAllowed  = errors.New("immediate deletion is not allowed for this account")
	ErrAccountNotActive             = errors.New("account is not active, cannot perform this action")
)