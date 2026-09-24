package valueobject

type ConsentStatus int8

const (
	ConsentRevoked ConsentStatus = 0
	ConsentGranted ConsentStatus = 1
	ConsentExpired ConsentStatus = 2
	ConsentDeleted ConsentStatus = 3
)

func (c ConsentStatus) IsValid() bool {
	switch c {
	case ConsentGranted, ConsentRevoked, ConsentExpired, ConsentDeleted:
		return true
	}
	return false
}

func (c ConsentStatus) String() string {
	switch c {
	case ConsentGranted:
		return "GRANTED"
	case ConsentRevoked:
		return "REVOKED"
	case ConsentExpired:
		return "EXPIRED"
	case ConsentDeleted:
		return "DELETED"
	}
	return "UNKNOWN"
}