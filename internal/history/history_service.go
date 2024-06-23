package history

import (
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/users"
)

type HistoryService interface {
	SaveHistoryRecord(dto *HistoryRecordDto) (*HistoryRecordDto, error)
	GetHistoryRecordsByUserId(userDto *users.UserDto, page, limit int) (*PaginatedHistoryResponse, error)
	DeleteHistory(userId string) error
	DeleteHistoryRecordById(id, userId string) error
}

type HistoryServiceImpl struct {
	repository HistoryRepository
}

func NewHistoryService(repository HistoryRepository) *HistoryServiceImpl {
	return &HistoryServiceImpl{repository: repository}
}

func (s *HistoryServiceImpl) SaveHistoryRecord(dto *HistoryRecordDto) (*HistoryRecordDto, error) {
	logger.Logger.Info("Saving history record...", "userId", dto.UserId)

	historyRecord := ToHistoryRecordModel(dto)

	err := s.repository.Save(historyRecord)
	if err != nil {
		logger.Logger.Error("Failed to save history record", "userId", dto.UserId)
		return nil, service_errors.NewErrInternalServer("Failed to save history record")
	}

	logger.Logger.Info("Saved history record.", "id", historyRecord.Id, "userId", dto.UserId)
	return ToHistoryRecordDto(historyRecord), nil
}

func (s *HistoryServiceImpl) GetHistoryRecordsByUserId(userDto *users.UserDto, page, limit int) (*PaginatedHistoryResponse, error) {
	logger.Logger.Info("Fetching history records...", "userId", userDto.Id, "page", page, "limit", limit)

	offset := (page - 1) * limit

	records, err := s.repository.FindByUserId(userDto.Id, offset, limit)
	if err != nil {
		logger.Logger.Error("Failed to fetch history records", "userId", userDto.Id)
		return nil, service_errors.NewErrInternalServer("Failed to fetch history records")
	}

	totalRecords, err := s.repository.CountByUserId(userDto.Id)
	if err != nil {
		return nil, err
	}

	totalPages := totalRecords / limit
	if totalRecords%limit != 0 {
		totalPages++
	}

	dtos := make([]HistoryRecordDto, len(records))
	for i, record := range records {
		dtos[i] = *ToHistoryRecordDto(&record)
	}

	response := &PaginatedHistoryResponse{
		Records:      dtos,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
		HasMore:      page < totalPages,
	}

	logger.Logger.Info("Fetched history records", "userId", userDto.Id, "size", len(dtos))
	return response, nil
}

func (s *HistoryServiceImpl) DeleteHistory(userId string) error {
	logger.Logger.Info("Deleting history...", "userId", userId)

	err := s.repository.DeleteByUserId(userId)
	if err != nil {
		logger.Logger.Error("Failed to delete history", "userId", userId)
		return service_errors.NewErrInternalServer("Failed to delete history")
	}

	logger.Logger.Info("Deleted history.", "userId", userId)
	return nil
}

func (s *HistoryServiceImpl) DeleteHistoryRecordById(id, userId string) error {
	logger.Logger.Info("Deleting history record...", "id", id, "userId", userId)

	err := s.repository.DeleteById(id, userId)
	if err != nil {
		logger.Logger.Error("Failed to delete history record", "id", id, "userId", userId)
		return service_errors.NewErrInternalServer("Failed to delete history record")
	}

	logger.Logger.Info("Deleted history record.", "id", id, "userId", userId)
	return nil
}
