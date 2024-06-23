package users

import (
	"time"
)

type UserDto struct {
	Id         string     `json:"id"`
	Email      string     `json:"email"`
	Password   string     `json:"-"`
	Username   string     `json:"username"`
	Role       UserRole   `json:"role"`
	Status     UserStatus `json:"status"`
	Provider   string     `json:"-"`
	PictureUrl string     `json:"picture"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"-"`
}

func ToUserModel(dto *UserDto) *User {
	return &User{
		Id:         dto.Id,
		Email:      dto.Email,
		Password:   dto.Password,
		Username:   dto.Username,
		Role:       dto.Role,
		Status:     dto.Status,
		Provider:   dto.Provider,
		PictureURL: dto.PictureUrl,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
	}
}

func ToUserDto(model *User) *UserDto {
	return &UserDto{
		Id:         model.Id,
		Email:      model.Email,
		Password:   model.Password,
		Username:   model.Username,
		Role:       model.Role,
		Status:     model.Status,
		Provider:   model.Provider,
		PictureUrl: model.PictureURL,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}
