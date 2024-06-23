package history

import (
	"gorm.io/gorm"
)

type HistoryRepository interface {
	Save(record *HistoryRecord) error
	FindByUserId(userId string, offset, limit int) ([]HistoryRecord, error)
	CountByUserId(userId string) (int, error)
	DeleteByUserId(userId string) error
	DeleteById(id, userId string) error
}

type HistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) *HistoryRepositoryImpl {
	return &HistoryRepositoryImpl{db: db}
}

func (r *HistoryRepositoryImpl) Save(record *HistoryRecord) error {
	result := r.db.Save(record)
	return result.Error
}

func (r *HistoryRepositoryImpl) FindByUserId(userId string, offset, limit int) ([]HistoryRecord, error) {
	var records []HistoryRecord

	result := r.db.Where("user_id = ?", userId).Offset(offset).Limit(limit).Order("created_at desc").Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}

func (rep *HistoryRepositoryImpl) CountByUserId(userId string) (int, error) {
	var count int64
	result := rep.db.Model(&HistoryRecord{}).Where("user_id = ?", userId).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

func (r *HistoryRepositoryImpl) DeleteByUserId(userId string) error {
	result := r.db.Delete(&HistoryRecord{}, "user_id = ?", userId)
	return result.Error
}

func (r *HistoryRepositoryImpl) DeleteById(id, userId string) error {
	result := r.db.Delete(&HistoryRecord{}, "id = ? AND user_id = ?", id, userId)
	return result.Error
}
