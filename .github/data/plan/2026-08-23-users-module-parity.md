# Plan: Users Module Parity — เติมส่วนที่ขาดจาก gistdaapi + แก้ security/dead-code

Repo: `C:\github\icmongolang` (Go 1.26, chi v5, GORM + Postgres, go-redis, asynq)

## ภาพรวมการแก้ไข

| # | Layer | ไฟล์ | สิ่งที่แก้ |
|---|-------|------|-----------|
| 1 | model | `internal/models/sd_user_roles_access.go` | สร้างใหม่ — entity ตาราง `sd_user_roles_access` |
| 2 | model | `internal/models/sd_user_roles_permission.go` | สร้างใหม่ — entity ตาราง `sd_user_roles_permision` |
| 3 | types/interface | `internal/modules/users/pg_repository.go` | เพิ่ม read-model types + 6 method ใน interface, ลบ `GetByResetToken` |
| 4 | infra | `internal/redis_repository.go`, `internal/repository/redis.go` | เพิ่ม `Expire` (TTL ให้ refresh-token set) |
| 5 | config | `config/config.go`, `config/config.default.yml` | เพิ่ม `upload` config (dir, max size) |
| 6 | repo impl | `internal/modules/users/repository/pg_repository.go` | impl 6 method ใหม่ + ลบ dead code (~12 methods ผิด entity) |
| 7 | redis impl | `internal/modules/users/redis_repository.go`, `repository/redis_repository.go` | เพิ่ม `SetDetail`/`GetDetail` สำหรับ profile cache |
| 8 | interface | `internal/modules/users/usecase.go` | เพิ่ม 7 use case method + profile key builder |
| 9 | usecase impl | `internal/modules/users/usecase/usecase.go` | impl method ใหม่ + fix enumeration/TTL/plaintext/hardcoded URL/cache |
| 10 | presenter | `internal/modules/users/presenter/presenters.go` | ตัด field sensitive + DTO ใหม่ |
| 11 | handler | `internal/modules/users/delivery/http/handlers.go` | handler ใหม่ 7 ตัว + ลบ PasswordTemp mapping |
| 12 | route | `internal/modules/users/delivery/http/routes.go` | ลบ route unauth + เพิ่ม route ใหม่ |
| 13 | server | `internal/server/handlers.go` | mount static `/uploads/avatar/*` |
| 14 | test | `*_test.go` (3 ไฟล์) | unit tests repo/usecase/handler |

Unit test decision: **เขียน** (ดู section `## Unit Tests` ท้ายไฟล์)

---

## Section 1: Model ใหม่ — `sd_user_roles_access`

**ไฟล์:** `internal/models/sd_user_roles_access.go` (สร้างใหม่)

Composite PK (`role_id`, `role_type_id`) ตาม schema จริงของ gistdaapi (`entities/rolesaccess.entity.ts`) — ชื่อตาราง/คอลัมน์ต้องตรงกับ DB ปัจจุบัน:

```go
package models

import "time"

type SdUserRolesAccess struct {
	RoleID     int        `gorm:"column:role_id;primaryKey"`
	RoleTypeID int        `gorm:"column:role_type_id;primaryKey"`
	CreateDate *time.Time `gorm:"column:create"`
	UpdateDate *time.Time `gorm:"column:update"`
}

func (SdUserRolesAccess) TableName() string { return "sd_user_roles_access" }
```

## Section 2: Model ใหม่ — `sd_user_roles_permision`

**ไฟล์:** `internal/models/sd_user_roles_permission.go` (สร้างใหม่)

- ⚠️ ชื่อตารางสะกด `permision` (s เดียว) ตาม schema จริงของ gistdaapi — **ห้าม** สะกด `permission`
- ⚠️ คอลัมน์ `select` เป็น reserved word ฝั่ง Postgres (`insert/update/delete/log/config/truncate` เป็น unreserved) — GORM quote ให้ตอน write; raw SELECT ต้อง quote เอง (Section 6)

```go
package models

import "time"

type SdUserRolePermission struct {
	RoleTypeID int        `gorm:"column:role_type_id;primaryKey"`
	Name       string     `gorm:"column:name"`
	Detail     *string    `gorm:"column:detail"`
	CreatedAt  *time.Time `gorm:"column:created"`
	UpdatedAt  *time.Time `gorm:"column:updated"`
	Insert     int        `gorm:"column:insert"`
	Update     int        `gorm:"column:update"`
	Delete     int        `gorm:"column:delete"`
	Select     int        `gorm:"column:select"`
	Log        int        `gorm:"column:log"`
	Config     int        `gorm:"column:config"`
	Truncate   int        `gorm:"column:truncate"`
}

func (SdUserRolePermission) TableName() string { return "sd_user_roles_permision" }
```

## Section 3: Read-model types + interface — `internal/modules/users/pg_repository.go`

**From:** interface มีแค่ CRUD + auth lookups; joined-profile/filter/statistics มีอยู่ใน impl แต่ไม่อยู่ใน interface (handler เรียกไม่ได้); `GetByResetToken` query column `reset_token` ซึ่งไม่มีใน `sd_user`

**To:** types ระดับ package `users` ประกาศในไฟล์นี้ (กัน import cycle — impl ใน package `repository` import `users` อยู่แล้ว) แล้วขยาย interface

เพิ่ม types (append ท้ายไฟล์):

```go
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
```

Interface เปลี่ยน (ไฟล์เดียวกัน):

```diff
 type UserPgRepository interface {
 	internal.PgRepository[models.SdUser]
 	...
-	GetByResetToken(ctx context.Context, resetToken string) (*models.SdUser, error)
+
+	// RBAC profile (join chain: sd_user → sd_user_role → sd_user_roles_access → sd_user_roles_permision)
+	GetProfileWithPermissions(ctx context.Context, id uuid.UUID) (*UserProfileDetail, error)
+	// Admin listing: keyword/status/active_status/sort/page + total count
+	ListUsers(ctx context.Context, f UserFilter) ([]*models.SdUser, int64, error)
+	// Statistics (CASE aggregation)
+	GetStatistics(ctx context.Context) ([]UserStatusCount, error)
+	// Target lists สำหรับ fan-out ตาม notification preference
+	GetActiveByNotification(ctx context.Context, ch NotificationChannel) ([]*models.SdUser, error)
+	UpdateActiveStatus(ctx context.Context, id uuid.UUID, status int16) error
+	UpdateAvatar(ctx context.Context, id uuid.UUID, avatarPath, avatar string) error
```

import เพิ่ม: `"strings"`, `"github.com/google/uuid"`

**เหตุผลลบ `GetByResetToken`:** column `reset_token` ไม่มีใน `sd_user` (model มี `password_reset_token`) และไม่มี caller — reset flow ใช้ `GetByResetTokenResetAt` อยู่แล้ว

## Section 4: Infra — Redis `Expire`

**`internal/redis_repository.go`** — เพิ่มใน interface `RedisRepository[M]`:

```go
	Expire(ctx context.Context, key string, seconds int) error
```

**`internal/repository/redis.go`** — เพิ่ม impl:

```go
func (r *RedisRepo[M]) Expire(ctx context.Context, key string, seconds int) error {
	return r.RedisClient.Expire(ctx, key, time.Second*time.Duration(seconds)).Err()
}
```

> หลังเพิ่ม method: รัน `go build ./...` ทันที — ถ้ามี struct อื่น implement `internal.RedisRepository` ตรง ๆ (ไม่ผ่าน embed `repository.RedisRepo`) จะ compile fail ตรงนั้น ให้เติม method ให้ครบ

## Section 5: Config

**`config/config.go`:**

เพิ่ม field ใน struct `Config` (หลัง `VectorDB`):

```go
	Upload         UploadConfig        `mapstructure:"upload"`
```

เพิ่ม type:

```go
type UploadConfig struct {
	Dir            string `mapstructure:"dir"`              // root dir ของไฟล์อัปโหลด
	MaxAvatarBytes int64  `mapstructure:"max_avatar_bytes"` // default 2MB
}
```

**`config/config.default.yml`** — เพิ่ม block (ท้ายไฟล์):

```yaml
upload:
  dir: "./uploads"
  max_avatar_bytes: 2097152
```

(`BindEnvs` walk ทุก field อยู่แล้ว → env override `UPLOAD_DIR`, `UPLOAD_MAX_AVATAR_BYTES` ใช้ได้ทันที)

## Section 6: Repo impl — `internal/modules/users/repository/pg_repository.go`

### 6.1 ลบ dead code

**From:** methods ที่ไม่อยู่ใน interface + หลายตัวทำงานบน entity ผิด (`models.User` ตาราง `user`, preload relations ที่ไม่มีจริง):

- `InsertUserWithFields` (146-180), `InsertUserWithMap` (205-262), `BatchInsertUsers` (266-272)
- `FilterRequest` struct (285-297), `FilterUsers` (301-396)
- `UserStatusCount` type (400-405), `GetUserStatistics` (409-434) — logic port ไป method ใหม่ชื่อ `GetStatistics`
- `BulkUpdateStatus` (438-447), `UpdateRoleBasedOnSuperuser` (451-457), `GetUsersByDateRange` (461-469), `GetActiveUsersWithRoleNames` (473-481)
- `CreateUser(models.User)` (485-489), `GetUserByID(uint)` (493-508)
- `UserFilterScopes`+`ToScopes` (512-556), `GetByUsername` (557-573), `FindWithScopes` (577-581)
- `GetByResetToken` (119-125) ตาม Section 3

**To:** ลบทิ้งทั้งหมด; import ปรับ — ลบ `"icmongolang/pkg/cryptpass"` (ไม่ใช้แล้ว), เพิ่ม `"fmt"`

ก่อนลบรายชื่อ: รัน `grep -rn "<MethodName>" internal cmd` ยืนยันไม่มี caller — ถ้าเจอ caller ให้หยุดรายงาน ไม่ลบ

### 6.2 Methods ใหม่ (append ท้ายไฟล์)

```go
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
	// ⚠️ quote reserved-word columns ด้วย double quotes
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

// GetStatistics – CASE aggregation (port จาก GetUserStatistics เดิม)
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
		Updates(map[string]interface{}{"active_status": status, "updateddate": time.Now()})
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
		Updates(map[string]interface{}{"avatarpath": avatarPath, "avatar": avatar, "updateddate": time.Now()})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
```

import เพิ่มเพิ่มเติมสำหรับไฟล์นี้: `"github.com/google/uuid"`

## Section 7: Redis repo — profile cache methods

**`internal/modules/users/redis_repository.go`:**

```diff
 type UserRedisRepository interface {
 	internal.RedisRepository[models.SdUser]
+	SetDetail(ctx context.Context, key string, exp *UserProfileDetail, seconds int) error
+	GetDetail(ctx context.Context, key string) (*UserProfileDetail, error)
 }
```

(ประกาศใน package `users` จึงอ้าง `UserProfileDetail` ได้ตรง — ไม่แตะ generic interface)

**`internal/modules/users/repository/redis_repository.go`** — เพิ่ม impl (struct `UserRedisRepo` มี field `repository.RedisRepo[models.SdUser]` embed ซึ่งเข้าถึง `RedisClient`):

```go
func (r *UserRedisRepo) SetDetail(ctx context.Context, key string, d *users.UserProfileDetail, seconds int) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return r.RedisClient.Set(ctx, key, b, time.Second*time.Duration(seconds)).Err()
}

func (r *UserRedisRepo) GetDetail(ctx context.Context, key string) (*users.UserProfileDetail, error) {
	b, err := r.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var d users.UserProfileDetail
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	return &d, nil
}
```

import เพิ่ม: `"encoding/json"`, `"errors"`, `"time"`, `"github.com/redis/go-redis/v9"` (บางส่วนอาจมีอยู่แล้ว)

> หมายเหตุ: ถ้า `CreateUserRedisRepository(...)` return `users.UserRedisRepository` ผ่าน wrapper struct ที่ไม่ expose `RedisClient` ให้เพิ่ม field `RedisClient *redis.Client` ตอน construct ใน factory function

## Section 8: Usecase interface — `internal/modules/users/usecase.go`

เพิ่มต่อท้าย interface:

```go
	GenerateRedisProfileKey(id uuid.UUID) string
	Profile(ctx context.Context, id uuid.UUID) (*UserProfileDetail, error)
	ListUsers(ctx context.Context, f UserFilter) ([]*models.SdUser, int64, error)
	Statistics(ctx context.Context) ([]UserStatusCount, error)
	ActiveByNotification(ctx context.Context, ch NotificationChannel) ([]*models.SdUser, error)
	UpdateActiveStatus(ctx context.Context, id uuid.UUID, status int16) error
	UpdateAvatar(ctx context.Context, id uuid.UUID, fh *multipart.FileHeader) error
```

import เพิ่ม: `"mime/multipart"`

## Section 9: Usecase impl — `internal/modules/users/usecase/usecase.go`

### 9.1 Fixes ของเดิม

**(a) Account-enumeration leak — `SignInByUsername` (บรรทัด 62)**

```diff
-		return "", "", httpErrors.ErrNotFound(err)
+		return "", "", httpErrors.ErrWrongPassword(errors.New("invalid credentials"))
```

**(b) Refresh-token TTL** — เพิ่ม helper แล้วแทนที่ call site `redisRepo.Sadd(ctx, u.GenerateRedisRefreshTokenKey(...), ...)` ทั้ง 3 จุด (`SignInByUsername` ~76, `SignIn` ~336, `Refresh` ~482):

```go
func (u *userUseCase) storeRefreshToken(ctx context.Context, id uuid.UUID, refreshToken string) error {
	key := u.GenerateRedisRefreshTokenKey(id)
	if err := u.redisRepo.Sadd(ctx, key, refreshToken); err != nil {
		return err
	}
	ttlSeconds := int(u.Cfg.Jwt.RefreshTokenExpireDuration) * 60 // config หน่วยนาที
	if err := u.redisRepo.Expire(ctx, key, ttlSeconds); err != nil {
		u.logger.Warnf("Failed to set TTL on refresh token set %s: %v", key, err)
	}
	return nil
}
```

**(c) Plaintext password defensive** — ใน `Create` หลัง hash (บรรทัด ~192) เพิ่ม:

```go
	exp.PasswordTemp = nil // ห้ามเก็บ plaintext ค้างใน model
```

**(d) Hardcoded URL ใน email template (บรรทัด 249, 580)**

```diff
-		fmt.Sprintf("http://localhost:5000/auth/verifyemail?code=%s", verificationCode),
+		fmt.Sprintf("%s/auth/verifyemail?code=%s", strings.TrimRight(u.Cfg.Server.BaseUrl, "/"), verificationCode),
```
```diff
-		fmt.Sprintf("http://localhost:5000/auth/resetpassword?code=%s", resetToken),
+		fmt.Sprintf("%s/auth/resetpassword?code=%s", strings.TrimRight(u.Cfg.Server.BaseUrl, "/"), resetToken),
```

### 9.2 Cache-invalidation helper

```go
func (u *userUseCase) invalidateUserCaches(ctx context.Context, id uuid.UUID) {
	_ = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(id))
	_ = u.redisRepo.Delete(ctx, u.GenerateRedisProfileKey(id))
}
```

แทนที่ single `_ = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(...))` ใน `Update` (~156), `Delete` (~129), `Verify` (~539), `ForgotPassword` (~569), `ResetPassword` (~628), `UpdatePassword` (~422) ด้วย helper นี้

คงการลบ refresh-token set แบบเดิมเฉพาะจุด kill-session: `Delete`, `UpdatePassword`, `ResetPassword`

### 9.3 Methods ใหม่ (append ท้ายไฟล์)

```go
func (u *userUseCase) GenerateRedisProfileKey(id uuid.UUID) string {
	return fmt.Sprintf("sd_user_profile:%s", id.String())
}

// Profile – joined profile + role title + permission flags (cache-through TTL 1h)
func (u *userUseCase) Profile(ctx context.Context, id uuid.UUID) (*users.UserProfileDetail, error) {
	cacheKey := u.GenerateRedisProfileKey(id)
	if cached, _ := u.redisRepo.GetDetail(ctx, cacheKey); cached != nil {
		return cached, nil
	}
	detail, err := u.pgRepo.GetProfileWithPermissions(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := u.redisRepo.SetDetail(ctx, cacheKey, detail, 3600); err != nil {
		u.logger.Warnf("Failed to cache profile %s: %v", id, err)
	}
	return detail, nil
}

func (u *userUseCase) ListUsers(ctx context.Context, f users.UserFilter) ([]*models.SdUser, int64, error) {
	u.logger.Infof("ListUsers keyword=%q sort=%s limit=%d offset=%d", f.Keyword, f.SortField, f.Limit, f.Offset)
	return u.pgRepo.ListUsers(ctx, f)
}

func (u *userUseCase) Statistics(ctx context.Context) ([]users.UserStatusCount, error) {
	return u.pgRepo.GetStatistics(ctx)
}

func (u *userUseCase) ActiveByNotification(ctx context.Context, ch users.NotificationChannel) ([]*models.SdUser, error) {
	switch ch {
	case users.NotifyEmail, users.NotifySms, users.NotifyLine:
	default:
		return nil, httpErrors.ErrValidation(fmt.Errorf("invalid notification channel: %s", ch))
	}
	return u.pgRepo.GetActiveByNotification(ctx, ch)
}

func (u *userUseCase) UpdateActiveStatus(ctx context.Context, id uuid.UUID, status int16) error {
	if status != 0 && status != 1 {
		return httpErrors.ErrValidation(errors.New("active_status must be 0 or 1"))
	}
	if _, err := u.Get(ctx, id); err != nil {
		return err
	}
	if err := u.pgRepo.UpdateActiveStatus(ctx, id, status); err != nil {
		return err
	}
	u.invalidateUserCaches(ctx, id)
	return nil
}

// UpdateAvatar – validate size/type → save disk → update DB → invalidate cache
func (u *userUseCase) UpdateAvatar(ctx context.Context, id uuid.UUID, fh *multipart.FileHeader) error {
	maxBytes := u.Cfg.Upload.MaxAvatarBytes
	if maxBytes <= 0 {
		maxBytes = 2 << 20 // 2MB default
	}
	if fh.Size > maxBytes {
		return httpErrors.ErrValidation(fmt.Errorf("avatar exceeds max size %d bytes", maxBytes))
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png":
	default:
		return httpErrors.ErrValidation(errors.New("avatar must be jpg/jpeg/png"))
	}

	dir := u.Cfg.Upload.Dir
	if dir == "" {
		dir = "./uploads"
	}
	avatarDir := filepath.Join(dir, "avatar")
	if err := os.MkdirAll(avatarDir, 0o755); err != nil {
		return err
	}
	filename := id.String() + ext
	dstPath := filepath.Join(avatarDir, filename)

	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	publicPath := "/uploads/avatar/" + filename
	if err := u.pgRepo.UpdateAvatar(ctx, id, publicPath, filename); err != nil {
		_ = os.Remove(dstPath) // rollback file ถ้า DB fail
		return err
	}
	u.invalidateUserCaches(ctx, id)
	return nil
}
```

import เพิ่ม: `"encoding/json"` (ถ้าใช้), `"io"`, `"mime/multipart"`, `"os"`, `"path/filepath"`

## Section 10: Presenter — `internal/modules/users/presenter/presenters.go`

**(a) ตัด sensitive fields ออกจาก `UserResponse`** — ลบ 4 field (mapModelResponse เดิมไม่ populate อยู่แล้ว → ไม่ break client):

```diff
-	RefreshToken          *string    `json:"refresh_token,omitempty"`
-	VerifiedCode          *string    `json:"verification_code,omitempty"`
-	PasswordResetToken    *string    `json:"password_reset_token,omitempty"`
-	PasswordResetAt       *string    `json:"password_reset_at,omitempty"`
```

**(b) เพิ่ม DTO:**

```go
// PermissionFlag – response shape ของ permission flags
type PermissionFlag struct {
	Name     string `json:"name"`
	Detail   string `json:"detail,omitempty"`
	Insert   int    `json:"insert"`
	Update   int    `json:"update"`
	Delete   int    `json:"delete"`
	Select   int    `json:"select"`
	Log      int    `json:"log"`
	Config   int    `json:"config"`
	Truncate int    `json:"truncate"`
}

// UserProfileResponse – user + role_title + permissions
type UserProfileResponse struct {
	UserResponse
	RoleTitle   string           `json:"role_title,omitempty"`
	Permissions []PermissionFlag `json:"permissions,omitempty"`
}

// UserUpdateActiveStatus – PATCH /user/{id}/activestatus
type UserUpdateActiveStatus struct {
	ActiveStatus *int16 `json:"active_status" validate:"required,oneof=0 1"`
}
```

## Section 11: Handler — `internal/modules/users/delivery/http/handlers.go`

**(a) Fix mapModel (บรรทัด ~631):**

```diff
 		Password:     req.Password,
-		PasswordTemp: stringPtr(req.Password),
```

**(b) Mapper ใหม่ (append ท้ายไฟล์ ใน zone Mapping functions):**

```go
func mapProfileResponse(d *users.UserProfileDetail) *presenter.UserProfileResponse {
	base := mapModelResponse(d.User)
	if base == nil {
		return nil
	}
	resp := &presenter.UserProfileResponse{UserResponse: *base, RoleTitle: d.RoleTitle}
	resp.Permissions = make([]presenter.PermissionFlag, 0, len(d.Permissions))
	for _, p := range d.Permissions {
		item := presenter.PermissionFlag{
			Name: p.Name, Insert: p.Insert, Update: p.Update,
			Delete: p.Delete, Select: p.Select, Log: p.Log,
			Config: p.Config, Truncate: p.Truncate,
		}
		if p.Detail != nil {
			item.Detail = *p.Detail
		}
		resp.Permissions = append(resp.Permissions, item)
	}
	return resp
}
```

**(c) Handlers ใหม่ 7 ตัว (pattern เดิมทั้งหมด: parse → validate → usecase → render):**

```go
// Profile – GET /user/profile/{id} : usersUC.Profile(id) → mapProfileResponse
// ProfileMe – GET /user/me ปรับ: เรียก usersUC.Profile(user.ID) แทนคืน raw ctx user
// ListUsers – GET /user/list : query keyword,status=csv,active_status=csv,sort=field-ASC|DESC,page,page_size → usersUC.ListUsers(UserFilter{...})
//   - page<1→1, page_size<1→20 (>100 clamp 100), offset=(page-1)*page_size
//   - status/active_status csv → []int16 ผ่าน strconv.ParseInt(...,10,16)
//   - sort "field-DESC" → SortField=field, SortOrder=DESC (split จากท้าย)
// Statistics – GET /user/statistics : usersUC.Statistics() → payload stats
// NotifyList – GET /user/notify/{channel} : chi.URLParam("channel") → NotificationChannel(ch) → ActiveByNotification → mapModelsResponse
// UpdateActiveStatus – PATCH /user/{id}/activestatus : body presenter.UserUpdateActiveStatus → ValidateStruct → usersUC.UpdateActiveStatus
// UploadAvatar – POST /user/me/avatar :
//   r.ParseMultipartForm(h.cfg.Upload.MaxAvatarBytes+1024)
//   file, header, err := r.FormFile("file") ; defer file.Close()
//   h.usersUC.UpdateAvatar(ctx, user.ID, header)  (user จาก middleware.GetUserFromCtx)
```

Swagger annotation ทุก endpoint: `@Security BearerAuth`, `@Tags users`, `@Success/@Failure` ตาม pattern เดิม

## Section 12: Routes — `internal/modules/users/delivery/http/routes.go`

```diff
 func MapUserRoute(router *chi.Mux, h users.Handlers, mw *middleware.MiddlewareManager) {
 	// Public
 	router.Post("/register", h.Register())
-	router.Post("/users", h.Create())   // ❌ ลบ — swagger บอก admin only; ตัวจริงคือ /api/user/ ใน SuperUser group
 	router.Post("/signin", h.SignInEmail())
 	router.Post("/login", h.SignInUsername())

 	router.Route("/user", func(r chi.Router) {
 		...middleware เดิม...

 		r.Get("/me", h.Me())                    // Me handler ปรับให้เรียก Profile (Section 11c)
 		r.Put("/me", h.UpdateMe())
 		r.Patch("/me/updatepass", h.UpdatePasswordMe())
+		r.Post("/me/avatar", h.UploadAvatar())
+		r.Get("/profile/{id}", h.Profile())

 		r.Group(func(admin chi.Router) {
 			admin.Use(mw.SuperUser())
 			admin.Get("/", h.GetMulti())
 			admin.Post("/", h.Create())
 			admin.Patch("/{id}/role", h.UpdateRole())
+			admin.Get("/list", h.ListUsers())
+			admin.Get("/statistics", h.Statistics())
+			admin.Get("/notify/{channel}", h.NotifyList())
+			admin.Patch("/{id}/activestatus", h.UpdateActiveStatus())
 		})

 		r.Route("/{id}", func(r chi.Router) { ...เดิม... })
 	})
 }
```

> ⚠️ ลำดับ route: `/list`, `/statistics`, `/notify/{channel}` ต้องลงทะเบียน **ก่อน** `r.Route("/{id}", ...)` ไม่งั้น chi จะ match `{id}` ก่อน (chi static beats param จริง แต่วางก่อนเพื่อความชัด)

**Interface `users.Handlers` (`internal/modules/users/handler.go`)** — เพิ่ม method signature ให้ครบ 7 handler ใหม่

## Section 13: Static serve — `internal/server/handlers.go`

หลัง global middleware block (~หลังบรรทัด 176) เพิ่ม:

```go
	// --- Static avatar uploads ---
	uploadRoot := cfg.Upload.Dir
	if uploadRoot == "" {
		uploadRoot = "./uploads"
	}
	avatarFS := http.FileServer(http.Dir(filepath.Join(uploadRoot, "avatar")))
	r.Handle("/uploads/avatar/*", http.StripPrefix("/uploads/avatar/", avatarFS))
```

(`http.Dir` clean path + ปฏิเสธ `..` อยู่แล้ว; import `"path/filepath"` ถ้ายังไม่มี)

## Section 14: Unit Tests — เขียน

Test files ใหม่ (mock hand-written):

| ไฟล์ | Cases |
|------|-------|
| `internal/modules/users/pg_repository_test.go` | `ResolveUserSortColumn`: whitelist คืน column ตรง / unknown → `createddate` / case-insensitive + trim; `UserFilter.SortOrderEffective`: asc→ASC, ""→DESC, DESC→DESC |
| `internal/modules/users/usecase/usecase_test.go` | mock pgRepo+redisRepo: `Profile` cache hit ไม่ยิง DB / cache miss เรียก repo + set cache / `UpdateActiveStatus` status=2 → ErrValidation, success เรียก invalidateUserCaches / `ActiveByNotification` channel="web" → ErrValidation, "email" → pass-through / `UpdateAvatar` size เกิน → error, ext=.gif → error, success → pgRepo.UpdateAvatar path `/uploads/avatar/<uuid>.png` + invalidate / `SignInByUsername` user-not-found → error message = "invalid credentials" (เท่ากับ SignIn) |
| `internal/modules/users/delivery/http/handlers_test.go` | httptest: `/api/user/list` ไม่มี superuser → 403 / `/api/user/statistics` superuser → 200 payload array / UploadAvatar multipart เกิน max size → 400 / route `POST /api/users` ถูกลบ → 404 |

Mock pattern: struct implement `users.UserPgRepository`, `users.UserRedisRepository`, `users.UserUseCaseI` (ฝั่ง handler test) — embed interface เพื่อไม่ต้อง implement ครบ:

```go
type mockPgRepo struct {
	users.UserPgRepository // panic ถ้าเรียก method ที่ไม่ stub
	profileFn func(ctx context.Context, id uuid.UUID) (*users.UserProfileDetail, error)
}
```

## Verification (รันท้ายงาน)

```bash
gofmt ./...
go build ./...
go vet ./...
go test ./internal/modules/users/...
swag init -g cmd/serve.go   # regenerate swagger docs
```

Manual smoke: register → signin → GET /api/user/me (มี role_title+permissions) → POST /api/user/me/avatar (jpg) → GET /uploads/avatar/<file> → GET /api/user/list?keyword=&sort=username-ASC → GET /api/user/statistics → PATCH /api/user/{id}/activestatus → GET /api/user/notify/email
