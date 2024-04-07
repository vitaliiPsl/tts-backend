package token

import "time"

type TokenDto struct {
	Id        string       `json:"id"`
	UserID    string       `json:"user_id"`
	Token     string       `json:"token"`
	Purpose   TokenPurpose `json:"purpose"`
	CreatedAt time.Time    `json:"created_at"`
	ExpiresAt time.Time    `json:"expires_at"`
}

func ToVerificationTokenModel(dto *TokenDto) *Token {
	return &Token{
		Id:        dto.Id,
		UserID:    dto.UserID,
		Token:     dto.Token,
		Purpose:   dto.Purpose,
		CreatedAt: dto.CreatedAt,
		ExpiresAt: dto.ExpiresAt,
	}
}

func ToVerificationTokenDto(model *Token) *TokenDto {
	return &TokenDto{
		Id:        model.Id,
		UserID:    model.UserID,
		Token:     model.Token,
		Purpose:   model.Purpose,
		CreatedAt: model.CreatedAt,
		ExpiresAt: model.ExpiresAt,
	}
}
