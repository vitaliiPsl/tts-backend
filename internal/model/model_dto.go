package model

import "time"

type ModelDto struct {
	Id        string    `json:"id"`
	Url       string    `json:"url"`
	Name      string    `json:"name"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"created_at"`
}

func ToModelModel(dto *ModelDto) *Model {
	return &Model{
		Id:        dto.Id,
		Url:       dto.Url,
		Name:      dto.Name,
		Language:  dto.Language,
		CreatedAt: dto.CreatedAt,
	}
}

func ToModelDto(model *Model) *ModelDto {
	return &ModelDto{
		Id:        model.Id,
		Url:       model.Url,
		Name:      model.Name,
		Language:  model.Language,
		CreatedAt: model.CreatedAt,
	}
}
