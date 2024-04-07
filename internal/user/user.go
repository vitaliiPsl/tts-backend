package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id         string     `gorm:"type:varchar(256);primaryKey;"`
	Email      string     `gorm:"type:varchar(256);uniqueIndex;not null"`
	Password   string     `gorm:"type:varchar(256);"`
	Username   string     `gorm:"type:varchar(256);"`
	Role       UserRole   `gorm:"type:varchar(256);default:'User'"`
	Status     UserStatus `gorm:"type:varchar(256);default:'Pending'"`
	Provider   string     `gorm:"type:varchar(256);"`
	PictureURL string     `gorm:"type:varchar(512);"`
	CreatedAt  time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Id = uuid.NewString()
	return
}
