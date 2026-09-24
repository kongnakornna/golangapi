package users

import (
	"context"
	"strings"
	"time"

	"icmongolang/internal"
	"icmongolang/internal/models"

	"github.com/google/uuid"
)

type UserPgRepository interface {
	internal.PgRepository[models.SdUser]
	GetByEmail(ctx context.Context, email string) (*models.SdUser, error)
	GetByEmailStrict(ctx context.Context, email string) (*models.SdUser, error)
	GetByUsernameStrict(ctx context.Context, username string) (*models.SdUser, error)
	UpdatePassword(ctx context.Context, exp *models.SdUser, newPassword string) (*models.SdUser, error)
	UpdateVerificationCode(ctx context.Context, exp *models.SdUser, newVerificationCode string) (*models.SdUser, error)
	UpdateVerification(ctx context.Context, exp *models.SdUser, newVerificationCode string, newVerified bool) (*models.SdUser, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (*models.SdUser, error)
	UpdatePasswordReset(ctx context.Context, exp *models.SdUser, passwordResetToken string, passwordResetAt time.Time) (*models.SdUser, error)
	GetByResetTokenResetAt(ctx context.Context, resetToken string, resetAt time.Time) (*models.SdUser, error)
	UpdatePasswordResetToken(ctx context.Context, exp *models.SdUser, newPassword string, resetToken string) (*models.SdUser, error)

	// เพิ่ม method ใหม่สำหรับ raw query
	GetMultiRaw(ctx context.Context, limit, offset int) ([]*models.SdUser, error) // GetMultiWithTotal คืนค่ารายชื่อผู้ใช้ + จำนวนทั้งหมด (support filter, sort)
	GetMultiWithTotal(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*models.SdUser, int64, error)

	// RBAC profile (join chain: sd_user → sd_user_role → sd_user_roles_access → sd_user_roles_permision)
	GetProfileWithPermissions(ctx context.Context, id uuid.UUID) (*UserProfileDetail, error)
	// Admin listing: keyword/status/active_status/sort/page + total count
	ListUsers(ctx context.Context, f UserFilter) ([]*models.SdUser, int64, error)
	// Statistics (CASE aggregation)
	GetStatistics(ctx context.Context) ([]UserStatusCount, error)
	// Target lists สำหรับ fan-out ตาม notification preference
	GetActiveByNotification(ctx context.Context, ch NotificationChannel) ([]*models.SdUser, error)
	UpdateActiveStatus(ctx context.Context, id uuid.UUID, status int16) error
	UpdateAvatar(ctx context.Context, id uuid.UUID, avatarPath, avatar string) error
}

// PermissionFlag – 1 แถวจาก sd_user_roles_permision (join ผ่าน access)
type PermissionFlag struct {
	Name     string  `gorm:"column:name" json:"name"`
	Detail   *string `gorm:"column:detail" json:"detail,omitempty"`
	Insert   int     `gorm:"column:insert" json:"insert"`
	Update   int     `gorm:"column:update" json:"update"`
	Delete   int     `gorm:"column:delete" json:"delete"`
	Select   int     `gorm:"column:select" json:"select"`
	Log      int     `gorm:"column:log" json:"log"`
	Config   int     `gorm:"column:config" json:"config"`
	Truncate int     `gorm:"column:truncate" json:"truncate"`
}

// UserProfileDetail – user + role title + permission flags
type UserProfileDetail struct {
	User        *models.SdUser   `json:"user"`
	RoleTitle   string           `json:"role_title"`
	Permissions []PermissionFlag `json:"permissions"`
}

// UserFilter – admin list filters (whitelist-driven)
type UserFilter struct {
	Keyword        string // ILIKE บน username OR email
	Statuses       []int16
	ActiveStatuses []int16
	SortField      string // whitelist: username|email|createddate|status|active_status
	SortOrder      string // ASC|DESC (default DESC)
	Limit          int
	Offset         int
}

func (f UserFilter) SortOrderEffective() string {
	if strings.EqualFold(strings.TrimSpace(f.SortOrder), "asc") {
		return "ASC"
	}
	return "DESC"
}

var allowedUserSortColumns = map[string]string{
	"username":      "username",
	"email":         "email",
	"createddate":   "createddate",
	"created_at":    "createddate",
	"status":        "status",
	"active_status": "active_status",
}

// ResolveUserSortColumn – pure function whitelist sort column (unit-testable)
func ResolveUserSortColumn(field string) string {
	if col, ok := allowedUserSortColumns[strings.ToLower(strings.TrimSpace(field))]; ok {
		return col
	}
	return "createddate"
}

// NotificationChannel – channel enum สำหรับ target list
type NotificationChannel string

const (
	NotifyEmail NotificationChannel = "email"
	NotifySms   NotificationChannel = "sms"
	NotifyLine  NotificationChannel = "line"
)

// UserStatusCount – ย้ายมาจาก repository/pg_repository.go (เดิมประกาศฝั่ง impl)
type UserStatusCount struct {
	Status     string `gorm:"column:status_label"`
	Count      int64  `gorm:"column:count"`
	RoleName   string `gorm:"column:role_name"`
	SuperCount int64  `gorm:"column:super_count"`
}
