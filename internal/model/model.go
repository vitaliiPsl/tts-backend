package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	Id        string    `gorm:"type:varchar(256);primaryKey;"`
	Url       string    `gorm:"type:varchar(256);"`
	Name      string    `gorm:"type:varchar(255);index:idx_unique_name_language,unique;"`
	Language  string    `gorm:"type:varchar(255);index:idx_unique_name_language,unique;"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (model *Model) BeforeCreate(tx *gorm.DB) (err error) {
	model.Id = uuid.NewString()
	return
}
