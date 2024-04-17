package history

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HistoryRecord struct {
	Id        string    `gorm:"type:varchar(256);not null;primaryKey;"`
	UserId    string    `gorm:"type:varchar(256);not null;index"`
	Text      string    `gorm:"type:varchar(2056);not null"`
	Language  string    `gorm:"type:varchar(256);"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (record *HistoryRecord) BeforeCreate(tx *gorm.DB) (err error) {
	record.Id = uuid.NewString()
	return
}
