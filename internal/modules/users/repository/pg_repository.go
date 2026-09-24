package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"icmongolang/internal/models"
	"icmongolang/internal/repository"
	"icmongolang/internal/modules/users"
	"icmongolang/pkg/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPgRepo struct {
	repository.PgRepo[models.SdUser]
}

func CreateUserPgRepository(db *gorm.DB) users.UserPgRepository {
	return &UserPgRepo{
		PgRepo: repository.CreatePgRepo[models.SdUser](db),
	}
}

func (r *UserPgRepo) GetByEmail(ctx context.Context, identifier string) (*models.SdUser, error) {
	normalized := strings.ToLower(strings.TrimSpace(identifier))

	var user models.SdUser

	// 1. Try email
	result := r.DB.WithContext(ctx).First(&user, "email = ?", normalized)
	if result.Error == nil {
		return &user, nil
	}
	// If it's a real DB error (not "not found"), return it immediately
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 2. Email not found → try username
	result = r.DB.WithContext(ctx).First(&user, "username = ?", normalized)
	if result.Error != nil {
		// This could be "not found" or any other error – return as is
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmailStrict looks up a user by the email column only (no username fallback).
func (r *UserPgRepo) GetByEmailStrict(ctx context.Context, email string) (*models.SdUser, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))

	var user models.SdUser
	if err := r.DB.WithContext(ctx).First(&user, "email = ?", normalized).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsernameStrict looks up a user by the username column only (no email fallback).
func (r *UserPgRepo) GetByUsernameStrict(ctx context.Context, username string) (*models.SdUser, error) {
	normalized := strings.ToLower(strings.TrimSpace(username))

	var user models.SdUser
	if err := r.DB.WithContext(ctx).First(&user, "username = ?", normalized).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserPgRepo) UpdatePassword(ctx context.Context, exp *models.SdUser, newPassword string) (*models.SdUser, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("password").
		Updates(map[string]interface{}{"password": newPassword}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *UserPgRepo) UpdateVerificationCode(ctx context.Context, exp *models.SdUser, newVerificationCode string) (*models.SdUser, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("verification_code").
		Updates(map[string]interface{}{"verification_code": newVerificationCode}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *UserPgRepo) UpdateVerification(ctx context.Context, exp *models.SdUser, newVerificationCode string, newVerified bool) (*models.SdUser, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("verification_code", "verified").
		Updates(map[string]interface{}{
			"verification_code": newVerificationCode,
			"verified":          newVerified,
		}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *UserPgRepo) GetByVerificationCode(ctx context.Context, verificationCode string) (*models.SdUser, error) {
	var obj *models.SdUser
	if result := r.DB.WithContext(ctx).First(&obj, "verification_code = ?", verificationCode); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *UserPgRepo) UpdatePasswordReset(ctx context.Context, exp *models.SdUser, passwordResetToken string, passwordResetAt time.Time) (*models.SdUser, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("password_reset_token", "password_reset_at").
		Updates(map[string]interface{}{
			"password_reset_token": passwordResetToken,
			"password_reset_at":    passwordResetAt,
		}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *UserPgRepo) GetByResetTokenResetAt(ctx context.Context, resetToken string, resetAt time.Time) (*models.SdUser, error) {
	var obj *models.SdUser
	if result := r.DB.WithContext(ctx).First(&obj, "password_reset_token = ? AND password_reset_at > ?", resetToken, resetAt); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *UserPgRepo) UpdatePasswordResetToken(ctx context.Context, exp *models.SdUser, newPassword string, resetToken string) (*models.SdUser, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("password", "password_reset_token").
		Updates(map[string]interface{}{"password": newPassword, "password_reset_token": resetToken}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

// GetProfileWithPermissions join chain: sd_user → sd_user_role → access → permision
func (r *UserPgRepo) GetProfileWithPermissions(ctx context.Context, id uuid.UUID) (*users.UserProfileDetail, error) {
	var user models.SdUser
	if err := r.DB.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	detail := &users.UserProfileDetail{User: &user}

	var roleTitle string
	if err := r.DB.WithContext(ctx).Table("sd_user_role").
		Select("title").
		Where("role_id = ?", user.RoleID).
		Limit(1).
		Scan(&roleTitle).Error; err != nil {
		return nil, err
	}
	detail.RoleTitle = roleTitle

	var perms []users.PermissionFlag
	// quote reserved-word columns ด้วย double quotes
	err := r.DB.WithContext(ctx).
		Table("sd_user_roles_permision p").
		Select(`p.name, p.detail, p."insert", p."update", p."delete", p."select", p."log", p."config", p."truncate"`).
		Joins("JOIN sd_user_roles_access a ON a.role_type_id = p.role_type_id").
		Where("a.role_id = ?", user.RoleID).
		Scan(&perms).Error
	if err != nil {
		return nil, err
	}
	if perms == nil {
		perms = []users.PermissionFlag{}
	}
	detail.Permissions = perms
	return detail, nil
}

// ListUsers – dynamic filter + sort (whitelist) + paging + total
func (r *UserPgRepo) ListUsers(ctx context.Context, f users.UserFilter) ([]*models.SdUser, int64, error) {
	q := r.DB.WithContext(ctx).Model(&models.SdUser{})

	if kw := strings.TrimSpace(f.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("username ILIKE ? OR email ILIKE ?", like, like)
	}
	if len(f.Statuses) > 0 {
		q = q.Where("status IN ?", f.Statuses)
	}
	if len(f.ActiveStatuses) > 0 {
		q = q.Where("active_status IN ?", f.ActiveStatuses)
	}

	q = q.Order(users.ResolveUserSortColumn(f.SortField) + " " + f.SortOrderEffective())

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	limit := f.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}
	var out []*models.SdUser
	if err := q.Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

// GetStatistics – CASE aggregation
func (r *UserPgRepo) GetStatistics(ctx context.Context) ([]users.UserStatusCount, error) {
	var stats []users.UserStatusCount
	err := r.DB.WithContext(ctx).Model(&models.SdUser{}).
		Select(`
			CASE WHEN status = 1 THEN 'Active' WHEN status = 0 THEN 'Inactive' ELSE 'Unknown' END as status_label,
			CASE role_id WHEN 1 THEN 'Super Admin' WHEN 2 THEN 'Normal User' ELSE 'Other' END as role_name,
			COUNT(*) as count,
			SUM(CASE WHEN is_superuser = true THEN 1 ELSE 0 END) as super_count
		`).
		Group("status_label, role_name").
		Order("status_label, role_name").
		Scan(&stats).Error
	return stats, err
}

// GetActiveByNotification – target list ตาม channel (column whitelist กัน injection)
func (r *UserPgRepo) GetActiveByNotification(ctx context.Context, ch users.NotificationChannel) ([]*models.SdUser, error) {
	col := map[users.NotificationChannel]string{
		users.NotifyEmail: "email_notification",
		users.NotifySms:   "sms_notification",
		users.NotifyLine:  "line_notification",
	}[ch]
	if col == "" {
		return nil, fmt.Errorf("invalid notification channel: %s", ch)
	}
	var out []*models.SdUser
	err := r.DB.WithContext(ctx).
		Where("status = ? AND "+col+" = ?", 1, 1).
		Order("createddate DESC").
		Find(&out).Error
	return out, err
}

// UpdateActiveStatus – toggle active_status (0/1)
func (r *UserPgRepo) UpdateActiveStatus(ctx context.Context, id uuid.UUID, status int16) error {
	res := r.DB.WithContext(ctx).Model(&models.SdUser{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"active_status": status, "updateddate": time.Now().In(helpers.GetTimeLocation())})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateAvatar – เก็บ path/filename avatar ลง sd_user
func (r *UserPgRepo) UpdateAvatar(ctx context.Context, id uuid.UUID, avatarPath, avatar string) error {
	res := r.DB.WithContext(ctx).Model(&models.SdUser{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"avatarpath": avatarPath, "avatar": avatar, "updateddate": time.Now().In(helpers.GetTimeLocation())})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetMultiRaw ดึงรายชื่อผู้ใช้แบบแบ่งหน้าโดยใช้ Raw SQL Query
// เรียงลำดับตามวันที่สร้างล่าสุด (createddate DESC)
func (r *UserPgRepo) GetMultiRaw(ctx context.Context, limit, offset int) ([]*models.SdUser, error) {
	var users []*models.SdUser
	query := `
        SELECT id, createddate, updateddate, deletedate, role_id, email, username,
               firstname, lastname, fullname, nickname, idcard, lastsignindate,
               status, active_status, network_id, remark, infomation_agree_status,
               gender, birthday, online_status, message, network_type_id,
               public_status, type_id, avatarpath, avatar, refresh_token, loginfailed,
               public_notification, sms_notification, email_notification, line_notification,
               mobile_number, phone_number, lineid, system_id, location_id, verified,
               verification_code, password_reset_token, password_reset_at, is_superuser
        FROM sd_user
        ORDER BY createddate DESC
    `
	err := r.DB.WithContext(ctx).Raw(query).Limit(limit).Offset(offset).Scan(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetMultiWithTotal คืนค่ารายชื่อผู้ใช้และจำนวนทั้งหมด (ใช้ GORM แทน raw query)
// GetMultiWithTotal returns users with total count, supports dynamic filters and sorting
func (r *UserPgRepo) GetMultiWithTotal(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*models.SdUser, int64, error) {
	var users []*models.SdUser
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.SdUser{})

	// Apply filters
	if email, ok := filters["email"].(string); ok && email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}
	if username, ok := filters["username"].(string); ok && username != "" {
		query = query.Where("username ILIKE ?", "%"+username+"%")
	}
	if fullname, ok := filters["fullname"].(string); ok && fullname != "" {
		query = query.Where("fullname ILIKE ?", "%"+fullname+"%")
	}
	if status, ok := filters["status"].(int16); ok {
		query = query.Where("status = ?", status)
	}
	if roleID, ok := filters["role_id"].(int); ok && roleID > 0 {
		query = query.Where("role_id = ?", roleID)
	}
	if verified, ok := filters["verified"].(bool); ok {
		query = query.Where("verified = ?", verified)
	}

	// Sorting
	sortBy := "createddate"
	sortOrder := "DESC"
	if sb, ok := filters["sort_by"].(string); ok && sb != "" {
		// Map allowed fields
		switch sb {
		case "email", "username", "fullname", "status", "role_id", "verified", "createddate", "updateddate":
			sortBy = sb
		case "created_at":
			sortBy = "createddate"
		case "updated_at":
			sortBy = "updateddate"
		}
	}
	if so, ok := filters["sort_order"].(string); ok && so == "asc" {
		sortOrder = "ASC"
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	err := query.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}
