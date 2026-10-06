package domain

type Role string

const (
	RoleUser    Role = "user"
	RoleManager Role = "manager"
	RoleAdmin   Role = "admin"
)

// IsValid ตรวจว่า role string ตรงกับ enum
func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleManager, RoleAdmin:
		return true
	}
	return false
}

// rank ใช้เทียบระดับสิทธิ์ (สูง = สิทธิ์มาก)
func (r Role) rank() int {
	switch r {
	case RoleAdmin:
		return 3
	case RoleManager:
		return 2
	case RoleUser:
		return 1
	}
	return 0
}

// AtLeast ใช้ใน RequireRole
func (r Role) AtLeast(min Role) bool {
	return r.rank() >= min.rank()
}

// RoleFromID map role_id จาก DB → Role
// ⚠️ ปรับตัวเลขให้ตรงกับ schema จริงของคุณ
func RoleFromID(roleID int) Role {
	switch roleID {
	case 1:
		return RoleAdmin
	case 2:
		return RoleManager
	default:
		return RoleUser
	}
}