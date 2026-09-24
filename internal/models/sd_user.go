package models

import (
	"time"

	"github.com/google/uuid"
)

type SdUser struct {
	ID                    uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedDate           time.Time  `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate           time.Time  `gorm:"column:updateddate;not null;default:now()"`
	DeleteDate            *time.Time `gorm:"column:deletedate;type:date"`
	RoleID                int        `gorm:"column:role_id;not null"`
	Email                 string     `gorm:"column:email;not null;uniqueIndex"`
	Username              string     `gorm:"column:username;not null;uniqueIndex"`
	Password              string     `gorm:"column:password;not null"`
	PasswordTemp          *string    `gorm:"column:password_temp"`
	Firstname             *string    `gorm:"column:firstname"`
	Lastname              *string    `gorm:"column:lastname"`
	Fullname              *string    `gorm:"column:fullname"`
	Nickname              *string    `gorm:"column:nickname"`
	IDCard                *string    `gorm:"column:idcard"`
	Lastsignindate        time.Time  `gorm:"column:lastsignindate;not null;default:now()"`
	Status                int16      `gorm:"column:status;not null"` // 1 = active, 0 = inactive
	ActiveStatus          *int16     `gorm:"column:active_status"`
	NetworkID             *int       `gorm:"column:network_id;default:1"`
	Remark                *string    `gorm:"column:remark"`
	InfomationAgreeStatus *int16     `gorm:"column:infomation_agree_status;default:0"`
	Gender                *string    `gorm:"column:gender"`
	Birthday              *time.Time `gorm:"column:birthday;type:date"`
	OnlineStatus          *string    `gorm:"column:online_status;default:'0'"`
	Message               *string    `gorm:"column:message"`
	NetworkTypeID         *int       `gorm:"column:network_type_id;default:0"`
	PublicStatus          *int16     `gorm:"column:public_status;default:0"`
	TypeID                *int       `gorm:"column:type_id;default:0"`
	Avatarpath            *string    `gorm:"column:avatarpath"`
	Avatar                *string    `gorm:"column:avatar"`
	RefreshToken          *string    `gorm:"column:refresh_token;type:text"`
	Loginfailed           *int16     `gorm:"column:loginfailed"`
	PublicNotification    *int16     `gorm:"column:public_notification;default:0"`
	SmsNotification       *int16     `gorm:"column:sms_notification;default:0"`
	EmailNotification     *int16     `gorm:"column:email_notification;default:0"`
	LineNotification      *int16     `gorm:"column:line_notification;default:0"`
	MobileNumber          *string    `gorm:"column:mobile_number;default:'0'"`
	PhoneNumber           *string    `gorm:"column:phone_number;default:'0'"`
	LineID                *string    `gorm:"column:lineid;default:'0'"`
	SystemID              *string    `gorm:"column:system_id;default:'1'"`
	LocationID            *string    `gorm:"column:location_id;default:'1'"`

	// Fields for authentication and verification
	Verified           bool       `gorm:"column:verified;default:false"`
	VerificationCode   string     `gorm:"column:verification_code;type:varchar(64);index"`
	PasswordResetToken string     `gorm:"column:password_reset_token;type:varchar(64);index"`
	PasswordResetAt    *time.Time `gorm:"column:password_reset_at"`

	// Added for superuser role
	IsSuperUser bool `gorm:"column:is_superuser;default:false"`
}

func (SdUser) TableName() string {
	return "sd_user"
}
