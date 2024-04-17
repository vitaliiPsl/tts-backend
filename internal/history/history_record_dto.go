package history

import "time"

type HistoryRecordDto struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Text      string    `json:"text"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"created_at"`
}

func ToHistoryRecordModel(dto *HistoryRecordDto) *HistoryRecord {
	return &HistoryRecord{
		Id:        dto.Id,
		UserId:    dto.UserId,
		Text:      dto.Text,
		Language:  dto.Language,
		CreatedAt: dto.CreatedAt,
	}
}

func ToHistoryRecordDto(model *HistoryRecord) *HistoryRecordDto {
	return &HistoryRecordDto{
		Id:        model.Id,
		UserId:    model.UserId,
		Text:      model.Text,
		Language:  model.Language,
		CreatedAt: model.CreatedAt,
	}
}
