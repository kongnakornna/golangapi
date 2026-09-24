package presenter

import (
	"time"

	"icmongolang/pkg/helpers"

	"github.com/google/uuid"
)

// UserCreate – used when creating a new user
type UserCreate struct {
	Username        string `json:"username" validate:"omitempty" example:"kongnakornna"`
	Email           string `json:"email" validate:"required,email" example:"kongnakornna@gmail.com"`
	Password        string `json:"password" validate:"required,min=8" example:"password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8" example:"password"`
	RoleID          int    `json:"role_id" validate:"required,min=1" example:"2"`
	Firstname       string `json:"firstname,omitempty" example:"Kongnakorn"`
	Lastname        string `json:"lastname,omitempty" example:"Jantakun"`
	Fullname        string `json:"fullname,omitempty" example:"Kongnakorn Jantakun"`
	MobileNumber    string `json:"mobile_number,omitempty" example:"0812345678"`
	PhoneNumber     string `json:"phone_number,omitempty" example:"021234567"`
	LineID          string `json:"line_id,omitempty" example:"kongnakorn_line"`
	LocationID      string `json:"location_id,omitempty" example:"loc_001"`
}

// UserUpdate – used when updating a user (all fields optional)
type UserUpdate struct {
	Firstname    *string `json:"firstname,omitempty"`
	Lastname     *string `json:"lastname,omitempty"`
	Fullname     *string `json:"fullname,omitempty"`
	MobileNumber *string `json:"mobile_number,omitempty"`
	PhoneNumber  *string `json:"phone_number,omitempty"`
	LineID       *string `json:"line_id,omitempty"`
	LocationID   *string `json:"location_id,omitempty"`
}

// UserUpdateRole – used by admin to change user role
type UserUpdateRole struct {
	RoleID int `json:"role_id" validate:"required"`
}

// UserUpdatePassword – used when changing password
type UserUpdatePassword struct {
	OldPassword     string `json:"old_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

// UserResponse – full user data returned by API
type UserResponse struct {
	ID                    uuid.UUID          `json:"id"`
	Email                 string             `json:"email"`
	Username              string             `json:"username"`
	RoleID                int                `json:"role_id"`
	Firstname             *string            `json:"firstname,omitempty"`
	Lastname              *string            `json:"lastname,omitempty"`
	Fullname              *string            `json:"fullname,omitempty"`
	Nickname              *string            `json:"nickname,omitempty"`
	MobileNumber          *string            `json:"mobile_number,omitempty"`
	PhoneNumber           *string            `json:"phone_number,omitempty"`
	LineID                *string            `json:"line_id,omitempty"`
	LocationID            *string            `json:"location_id,omitempty"`
	Avatar                *string            `json:"avatar,omitempty"`
	Gender                *string            `json:"gender,omitempty"`
	Birthday              *time.Time         `json:"birthday,omitempty"`
	Status                int16              `json:"status"`
	IsSuperUser           bool               `json:"is_superuser"`
	Verified              bool               `json:"verified"`
	LastSignIn            *helpers.LocalTime `json:"last_sign_in,omitempty"`
	PublicStatus          *string            `json:"public_status,omitempty"`
	IdCard                *string            `json:"idcard,omitempty"`
	LastSignInDate        *string            `json:"lastsignindate,omitempty"`
	ActiveStatus          *string            `json:"active_status,omitempty"`
	NetworkId             *string            `json:"network_id,omitempty"`
	Remark                *string            `json:"remark,omitempty"`
	InfomationAgreeStatus *string            `json:"infomation_agree_status,omitempty"`
	OnlineStatus          *string            `json:"online_status,omitempty"`
	Message               *string            `json:"message,omitempty"`
	NetworkTypeId         *string            `json:"network_type_id,omitempty"`
	TypeId                *string            `json:"type_id,omitempty"`
	AvatarPath            *string            `json:"avatarpath,omitempty"`
	LoginFailed           *string            `json:"loginfailed,omitempty"`
	PublicNotification    *string            `json:"public_notification,omitempty"`
	SmsNotification       *string            `json:"sms_notification,omitempty"`
	LineNotification      *string            `json:"line_notification,omitempty"`
	LineId                *string            `json:"lineid,omitempty"`
	SystemId              *string            `json:"system_id,omitempty"`
	CreatedAt             helpers.LocalTime  `json:"createddate"`
	UpdatedAt             helpers.LocalTime  `json:"updateddate"`
}

// Auth related DTOs (unchanged)
type UserSignIn struct {
	Email    string `json:"email" validate:"required" example:"kongnakornna@gmail.com"`
	Password string `json:"password" validate:"required,min=8" example:"password"`
}

// Auth related DTOs (unchanged)
type SignIn struct {
	Username string `json:"username" validate:"required" example:"kongnakornna"`
	Password string `json:"password" validate:"required,min=8" example:"password"`
}

type Token struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
}

type PublicKey struct {
	PublicKeyAccessToken  string `json:"public_key_access_token,omitempty"`
	PublicKeyRefreshToken string `json:"public_key_refresh_token,omitempty"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required" example:"kongnakornna@gmail.com"`
}

type ResetPassword struct {
	NewPassword     string `json:"new_password" validate:"required,min=8" example:"password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8" example:"password"`
}

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
