package model

import (
	"gorm.io/gorm"
)

type ModelRepository interface {
	Save(model *Model) error
	FindById(id string) (*Model, error)
	FindByNameAndLanguage(name, language string) (*Model, error)
	FindAll() ([]Model, error)
	DeleteById(id string) error
}

type ModelRepositoryImpl struct {
	db *gorm.DB
}

func NewModelRepository(db *gorm.DB) *ModelRepositoryImpl {
	return &ModelRepositoryImpl{db: db}
}

func (r *ModelRepositoryImpl) Save(model *Model) error {
	return r.db.Save(model).Error
}

func (r *ModelRepositoryImpl) FindById(id string) (*Model, error) {
	var model Model

	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *ModelRepositoryImpl) FindByNameAndLanguage(name, language string) (*Model, error) {
	var model *Model

	if err := r.db.Where("name = ? AND language = ?", name, language).First(&model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *ModelRepositoryImpl) FindAll() ([]Model, error) {
	var models []Model

	result := r.db.Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return models, nil
}

func (r *ModelRepositoryImpl) DeleteById(id string) error {
	return r.db.Delete(&Model{}, "id = ?", id).Error
}
