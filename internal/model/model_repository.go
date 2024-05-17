package model

import (
	"gorm.io/gorm"
)

type ModelRepository struct {
	db *gorm.DB
}

func NewModelRepository(db *gorm.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

func (r *ModelRepository) Save(model *Model) error {
	return r.db.Save(model).Error
}

func (r *ModelRepository) FindById(id string) (*Model, error) {
	var model *Model

	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *ModelRepository) FindByNameAndLanguage(name, language string) (*Model, error) {
	var model *Model

	if err := r.db.Where("name = ? AND language = ?", name, language).First(&model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *ModelRepository) FindAll() ([]Model, error) {
	var models []Model

	result := r.db.Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return models, nil
}

func (r *ModelRepository) DeleteById(id string) error {
	return r.db.Delete(&Model{}, "id = ?", id).Error
}
