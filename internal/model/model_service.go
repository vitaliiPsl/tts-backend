package model

import (
	"errors"
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/requests"

	"gorm.io/gorm"
)

type ModelService interface {
	SaveModel(req *requests.ModelRequest) (*ModelDto, error)
	UpdateModel(id string, req *requests.ModelRequest) (*ModelDto, error)
	DeleteModel(id string) error
	GetModelById(modelId string) (*ModelDto, error)
	GetModels() ([]ModelDto, error)
}

type ModelServiceImpl struct {
	repository ModelRepository
}

func NewModelService(repo ModelRepository) *ModelServiceImpl {

	return &ModelServiceImpl{repository: repo}
}

func (s *ModelServiceImpl) SaveModel(req *requests.ModelRequest) (*ModelDto, error) {
	logger.Logger.Info("Saving model...", "name", req.Name, "language", req.Language)

	existing, err := s.repository.FindByNameAndLanguage(req.Name, req.Language)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("Failed to model", "name", req.Name, "language", req.Language)
		return nil, service_errors.NewErrInternalServer("Failed to model by name and language")
	}

	if existing != nil {
		logger.Logger.Error("Model with given name and language already exists", "name", req.Name, "language", req.Language)
		return nil, service_errors.NewErrBadRequest("Model with this name and language already exists")
	}

	model := &Model{
		Url:      req.Url,
		Name:     req.Name,
		Language: req.Language,
	}

	err = s.repository.Save(model)
	if err != nil {
		logger.Logger.Error("Failed to save model", "name", model.Name, "language", "model.Language")
		return nil, service_errors.NewErrInternalServer("Failed to save model")
	}

	logger.Logger.Info("Saved model.", "id", model.Id, "name", model.Name, "language", model.Language)
	return ToModelDto(model), nil
}

func (s *ModelServiceImpl) UpdateModel(id string, req *requests.ModelRequest) (*ModelDto, error) {
	logger.Logger.Info("Updating model...", "id", id, "url", req.Url, "name", req.Name, "language", req.Language)

	model, err := s.repository.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("Model not found", "id", id)
			return nil, service_errors.NewErrNotFound("Model not found")
		}

		logger.Logger.Error("Failed to fetch model", "id", id)
		return nil, service_errors.NewErrInternalServer("Failed to fetch model")
	}

	existing, err := s.repository.FindByNameAndLanguage(req.Name, req.Language)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("Failed to fetch existing model", "name", req.Name, "language", req.Language)
		return nil, service_errors.NewErrInternalServer("Failed to fetch model by name and language")
	}

	if existing != nil && existing.Id != model.Id {
		logger.Logger.Error("Model with given name and language already exists", "name", req.Name, "language", req.Language)
		return nil, service_errors.NewErrBadRequest("Model with this name and language already exists")
	}

	if req.Url != "" {
		model.Url = req.Url
	}

	if req.Name != "" {
		model.Name = req.Name
	}

	if req.Language != "" {
		model.Language = req.Language
	}

	err = s.repository.Save(model)
	if err != nil {
		logger.Logger.Error("Failed to update model", "id", model.Id)
		return nil, service_errors.NewErrInternalServer("Failed to update model")
	}

	logger.Logger.Info("Updated model.", "id", model.Id, "url", model.Url, "name", model.Name, "language", model.Language)
	return ToModelDto(model), nil
}

func (s *ModelServiceImpl) DeleteModel(id string) error {
	logger.Logger.Info("Deleting model...", "id", id)

	model, err := s.repository.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("Model not found", "id", id)
			return service_errors.NewErrNotFound("Model not found")
		}

		logger.Logger.Error("Failed to fetch model", "id", id)
		return service_errors.NewErrInternalServer("Failed to fetch model by Id")
	}

	err = s.repository.DeleteById(model.Id)
	if err != nil {
		logger.Logger.Error("Failed to delete model", "id", id)
		return service_errors.NewErrInternalServer("Failed to delete model")
	}

	logger.Logger.Info("Deleted model.", "id", id)
	return nil
}

func (s *ModelServiceImpl) GetModelById(modelId string) (*ModelDto, error) {
	logger.Logger.Info("Fetching model...", "modelId", modelId)

	model, err := s.repository.FindById(modelId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("Model not found", "id", modelId)
			return nil, service_errors.NewErrNotFound("Model not found")
		}

		logger.Logger.Error("Failed to fetch models")
		return nil, service_errors.NewErrInternalServer("Failed to fetch models")
	}

	logger.Logger.Info("Fetched model", "modelId", modelId)
	return ToModelDto(model), nil
}

func (s *ModelServiceImpl) GetModels() ([]ModelDto, error) {
	logger.Logger.Info("Fetching models...")

	records, err := s.repository.FindAll()
	if err != nil {
		logger.Logger.Error("Failed to fetch models")
		return nil, service_errors.NewErrInternalServer("Failed to fetch models")
	}

	dtos := make([]ModelDto, len(records))
	for i, record := range records {
		dtos[i] = *ToModelDto(&record)
	}

	logger.Logger.Info("Fetched models", "size", len(dtos))
	return dtos, nil
}
