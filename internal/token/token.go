package token

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenPurpose string

const (
	PurposeEmailVerification TokenPurpose = "email_verification"
	PurposePasswordReset     TokenPurpose = "password_reset"
)

type Token struct {
	Id        string       `gorm:"type:varchar(256);primaryKey;"`
	UserID    string       `gorm:"type:varchar(255);index:idx_user_id;"`
	Token     string       `gorm:"type:varchar(255);index:idx_token,unique;"`
	Purpose   TokenPurpose `gorm:"type:varchar(255);"`
	CreatedAt time.Time    `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	ExpiresAt time.Time    `gorm:"type:timestamp;"`
}

func (token *Token) BeforeCreate(tx *gorm.DB) (err error) {
	token.Id = uuid.NewString()
	return
}
