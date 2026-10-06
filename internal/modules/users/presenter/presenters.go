package presenter

import (
	"time"

	"icmongolang/pkg/helpers"
)

// ============================================================
// AUTH
// ============================================================

type UserSignIn struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	NewPassword     string `json:"new_password"     validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// ============================================================
// USER
// ============================================================

type UserCreate struct {
	Email           string `json:"email"            validate:"required,email"`
	Username        string `json:"username"`
	Password        string `json:"password"         validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	RoleID          int    `json:"role_id"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Fullname        string `json:"fullname"`
	MobileNumber    string `json:"mobile_number"`
	PhoneNumber     string `json:"phone_number"`
	LineID          string `json:"lineid"`
	LocationID      string `json:"location_id"`
}

type UserUpdate struct {
	Firstname    *string `json:"firstname"`
	Lastname     *string `json:"lastname"`
	Fullname     *string `json:"fullname"`
	MobileNumber *string `json:"mobile_number"`
	PhoneNumber  *string `json:"phone_number"`
	LineID       *string `json:"lineid"`
	LocationID   *string `json:"location_id"`
}

type UserUpdatePassword struct {
	OldPassword     string `json:"old_password"     validate:"required"`
	NewPassword     string `json:"new_password"     validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type UserUpdateRole struct {
	RoleID int `json:"role_id" validate:"required,gt=0"`
}

type UserUpdateActiveStatus struct {
	ActiveStatus *int16 `json:"active_status" validate:"required,oneof=0 1"`
}

// UserResponse — ใช้ *helpers.LocalTime สำหรับ timestamp fields
type UserResponse struct {
	ID           interface{} `json:"id"`
	Email        string      `json:"email"`
	Username     string      `json:"username"`
	RoleID       int         `json:"role_id"`
	Firstname    *string     `json:"firstname,omitempty"`
	Lastname     *string     `json:"lastname,omitempty"`
	Fullname     *string     `json:"fullname,omitempty"`
	Nickname     *string     `json:"nickname,omitempty"`
	MobileNumber *string     `json:"mobile_number,omitempty"`
	PhoneNumber  *string     `json:"phone_number,omitempty"`
	LineID       *string     `json:"lineid,omitempty"`
	LocationID   *string     `json:"location_id,omitempty"`
	Avatar       *string     `json:"avatar,omitempty"`
	Gender       *string     `json:"gender,omitempty"`
	Birthday     interface{} `json:"birthday,omitempty"`
	Status       int16       `json:"status"`
	IsSuperUser  bool        `json:"is_superuser"`
	Verified     bool        `json:"verified"`

	// 🆕 ใช้ *helpers.LocalTime (ไม่ใช่ *string)
	LastSignIn *helpers.LocalTime `json:"last_sign_in,omitempty"`
	Remark     *string            `json:"remark,omitempty"`
	CreatedAt  *helpers.LocalTime `json:"createddate,omitempty"`
	UpdatedAt  *helpers.LocalTime `json:"updateddate,omitempty"`
}

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

type UserProfileResponse struct {
	UserResponse
	RoleTitle   string           `json:"role_title"`
	Permissions []PermissionFlag `json:"permissions"`
}

type UserListResponse struct {
	Data       []*UserResponse `json:"data"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	Total      int64           `json:"total"`
	TotalPages int             `json:"total_pages"`
}

func NewUserListResponse(list []*UserResponse, page, perPage int, total int64) *UserListResponse {
	totalPages := 0
	if perPage > 0 {
		totalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}
	return &UserListResponse{
		Data:       list,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

// ============================================================
// ACCESS MENU
// ============================================================

type AccessMenuCreate struct {
	UserTypeID int `json:"user_type_id" validate:"required"`
	RoleID     int `json:"role_id"      validate:"required"`
	MenuID     int `json:"menu_id"      validate:"required"`
	ParentID   int `json:"parent_id"`
}

type AccessMenuUpdate struct {
	UserTypeID *int `json:"user_type_id"`
	RoleID     *int `json:"role_id"`
	MenuID     *int `json:"menu_id"`
	ParentID   *int `json:"parent_id"`
}

type AccessMenuResponse struct {
	UserAccessID int `json:"user_access_id"`
	UserTypeID   int `json:"user_type_id"`
	RoleID       int `json:"role_id"`
	MenuID       int `json:"menu_id"`
	ParentID     int `json:"parent_id"`
}

type AccessMenuListResponse struct {
	Data       []*AccessMenuResponse `json:"data"`
	Page       int                   `json:"page"`
	PerPage    int                   `json:"per_page"`
	Total      int64                 `json:"total"`
	TotalPages int                   `json:"total_pages"`
}

func NewAccessMenuListResponse(list []*AccessMenuResponse, page, perPage int, total int64) *AccessMenuListResponse {
	totalPages := 0
	if perPage > 0 {
		totalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}
	return &AccessMenuListResponse{
		Data: list, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages,
	}
}

// ============================================================
// FILE
// ============================================================

type FileCreate struct {
	FileName   string `json:"file_name" validate:"required"`
	FileType   string `json:"file_type"`
	FilePath   string `json:"file_path"`
	FileTypeID int    `json:"file_type_id"`
	UID        string `json:"uid"`
	Status     int16  `json:"status"`
}

type FileUpdate struct {
	FileName *string `json:"file_name"`
	FilePath *string `json:"file_path"`
	Status   *int16  `json:"status"`
}

type FileResponse struct {
	ID         int       `json:"id"`
	FileName   string    `json:"file_name"`
	FileType   string    `json:"file_type,omitempty"`
	FilePath   string    `json:"file_path,omitempty"`
	FileTypeID int       `json:"file_type_id"`
	UID        string    `json:"uid"`
	FileDate   time.Time `json:"file_date"`
	Status     int16     `json:"status"`
}

type FileListResponse struct {
	Data       []*FileResponse `json:"data"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	Total      int64           `json:"total"`
	TotalPages int             `json:"total_pages"`
}

func NewFileListResponse(list []*FileResponse, page, perPage int, total int64) *FileListResponse {
	totalPages := 0
	if perPage > 0 {
		totalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}
	return &FileListResponse{Data: list, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}
}

// ============================================================
// LOG
// ============================================================

type LogCreate struct {
	LogTypeID    int     `json:"log_type_id" validate:"required"`
	UID          string  `json:"uid"         validate:"required,uuid"`
	Name         string  `json:"name"`
	Detail       string  `json:"detail"`
	IPAddress    *string `json:"ipaddress"`
	DeviceAccess *string `json:"device_access"`
	BrowserAgent *string `json:"browser_agent"`
	Status       int16   `json:"status"`
	Lang         string  `json:"lang"`
}

type LogResponse struct {
	ID           int       `json:"id"`
	LogTypeID    int       `json:"log_type_id"`
	UID          string    `json:"uid"`
	Name         string    `json:"name,omitempty"`
	Detail       string    `json:"detail,omitempty"`
	IPAddress    *string   `json:"ipaddress,omitempty"`
	DeviceAccess *string   `json:"device_access,omitempty"`
	BrowserAgent *string   `json:"browser_agent,omitempty"`
	Status       int16     `json:"status"`
	Create       time.Time `json:"create"`
	Update       time.Time `json:"update"`
	Lang         string    `json:"lang,omitempty"`
}

type LogListResponse struct {
	Data       []*LogResponse `json:"data"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"total_pages"`
}

func NewLogListResponse(list []*LogResponse, page, perPage int, total int64) *LogListResponse {
	totalPages := 0
	if perPage > 0 {
		totalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}
	return &LogListResponse{Data: list, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}
}

// ============================================================
// LOG TYPE
// ============================================================

type LogTypeCreate struct {
	TypeName   string `json:"type_name" validate:"required"`
	TypeDetail string `json:"type_detail"`
	Status     int16  `json:"status"`
}

type LogTypeUpdate struct {
	TypeName   *string `json:"type_name"`
	TypeDetail *string `json:"type_detail"`
	Status     *int16  `json:"status"`
}

// ============================================================
// MENU
// ============================================================

type MenuCreate struct {
	MenuParentID int    `json:"menu_parent_id"`
	ModulesID    int    `json:"modules_id" validate:"required"`
	Name         string `json:"name"        validate:"required"`
	Path         string `json:"path"`
	Icon         string `json:"icon"`
	Sort         int    `json:"sort"`
	Status       int16  `json:"status"`
}

type MenuUpdate struct {
	MenuParentID *int    `json:"menu_parent_id"`
	ModulesID    *int    `json:"modules_id"`
	Name         *string `json:"name"`
	Path         *string `json:"path"`
	Icon         *string `json:"icon"`
	Status       *int16  `json:"status"`
}

type MenuResponse struct {
	MenuID       int       `json:"menu_id"`
	MenuParentID int       `json:"menu_parent_id"`
	ModulesID    int       `json:"modules_id"`
	Name         string    `json:"name"`
	Path         string    `json:"path,omitempty"`
	Icon         string    `json:"icon,omitempty"`
	Sort         int       `json:"sort"`
	CreateDate   time.Time `json:"create_date"`
	UpdateDate   time.Time `json:"update_date"`
	Status       int16     `json:"status"`
}

type MenuListResponse struct {
	Data       []*MenuResponse `json:"data"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	Total      int64           `json:"total"`
	TotalPages int             `json:"total_pages"`
}

func NewMenuListResponse(list []*MenuResponse, page, perPage int, total int64) *MenuListResponse {
	totalPages := 0
	if perPage > 0 {
		totalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}
	return &MenuListResponse{Data: list, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}
}

// ============================================================
// MODULE
// ============================================================

type ModuleCreate struct {
	ModuleParentID int    `json:"module_parent_id"`
	Name           string `json:"name" validate:"required"`
	Path           string `json:"path"`
	Icon           string `json:"icon"`
	Sort           int    `json:"sort"`
	Status         int16  `json:"status"`
}

type ModuleUpdate struct {
	ModuleParentID *int    `json:"module_parent_id"`
	Name           *string `json:"name"`
	Path           *string `json:"path"`
	Icon           *string `json:"icon"`
	Status         *int16  `json:"status"`
}

type ModuleResponse struct {
	ModuleID       int       `json:"module_id"`
	ModuleParentID int       `json:"module_parent_id"`
	Name           string    `json:"name"`
	Path           string    `json:"path,omitempty"`
	Icon           string    `json:"icon,omitempty"`
	Sort           int       `json:"sort"`
	CreateDate     time.Time `json:"create_date"`
	UpdateDate     time.Time `json:"update_date"`
	Status         int16     `json:"status"`
}

type ModuleListResponse struct {
	Data       []*ModuleResponse `json:"data"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"total_pages"`
}

func NewModuleListResponse(list []*ModuleResponse, page, perPage int, total int64) *ModuleListResponse {
	totalPages := 0
	if perPage > 0 {
		totalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}
	return &ModuleListResponse{Data: list, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}
}

// ============================================================
// SORT
// ============================================================

type SortItem struct {
	ID       int  `json:"id"`
	Sort     int  `json:"sort"`
	ParentID *int `json:"parent_id,omitempty"`
}

type SortUpdateRequest struct {
	Items []SortItem `json:"items" validate:"required,dive"`
}