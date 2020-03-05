package models

import "time"

type PurposeType int

const (
	PurposeRegistration   PurposeType = 0
	PurposeForgotPassword PurposeType = 1
	PurposeMFA            PurposeType = 2
)

type Verification struct {
	BaseModel
	UserID    uint `gorm:"unique_index:'user_verification_purpose'"`
	User      User
	Code      string
	ExpiresOn time.Time
	Purpose   PurposeType `gorm:"unique_index:'user_verification_purpose'"`
}
