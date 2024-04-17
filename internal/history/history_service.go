package history

import (
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
)

type HistoryService struct {
	repository *HistoryRepository
}

func NewHistoryService(repository *HistoryRepository) *HistoryService {
	return &HistoryService{repository: repository}
}

func (s *HistoryService) SaveHistoryRecord(dto *HistoryRecordDto) (*HistoryRecordDto, error) {
	logger.Logger.Info("Saving history record...", "userId", dto.UserId)

	historyRecord := ToHistoryRecordModel(dto)

	err := s.repository.Save(historyRecord)
	if err != nil {
		logger.Logger.Error("Failed to save history record", "userId", dto.UserId)
		return nil, &service_errors.ErrInternalServer{Message: "Failed to save history record"}
	}

	logger.Logger.Info("Saved history record.", "id", historyRecord.Id, "userId", dto.UserId)
	return ToHistoryRecordDto(historyRecord), nil
}

func (s *HistoryService) GetHistoryRecordsByUserId(userId string, page, limit int) (*PaginatedHistoryResponse, error) {
	logger.Logger.Info("Fetching history records...", "userId", userId, "page", page, "limit", limit)

	offset := (page - 1) * limit

	records, err := s.repository.FindByUserId(userId, offset, limit)
	if err != nil {
		logger.Logger.Error("Failed to fetch history records", "userId", userId)
		return nil, &service_errors.ErrInternalServer{Message: "Failed to fetch history records"}
	}

	totalRecords, err := s.repository.CountByUserId(userId)
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

	logger.Logger.Info("Fetched history records", "userId", userId, "size", len(dtos))
	return response, nil
}

func (s *HistoryService) DeleteHistoryRecordsByUserId(userId string) error {
	logger.Logger.Info("Deleting history records...", "userId", userId)

	err := s.repository.DeleteByUserId(userId)
	if err != nil {
		logger.Logger.Error("Failed to delete history records", "userId", userId)
		return &service_errors.ErrInternalServer{Message: "Failed to delete history records"}
	}

	logger.Logger.Info("Deleted history records.", "userId", userId)
	return nil
}
