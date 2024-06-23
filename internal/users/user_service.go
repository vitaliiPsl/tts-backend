package users

import (
	"errors"

	"gorm.io/gorm"
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
)

type UserService interface {
	SaveUser(userDto *UserDto) (*UserDto, error)
	UpdateUser(id string, userDto *UserDto) (*UserDto, error)
	UpsertUser(userDto *UserDto) (*UserDto, error)
	FindById(id string) (*UserDto, error)
	FindByEmail(email string) (*UserDto, error)
}

type UserServiceImpl struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repository: repository}
}

func (s *UserServiceImpl) SaveUser(userDto *UserDto) (*UserDto, error) {
	logger.Logger.Info("Saving user...", "id", userDto.Id, "email", userDto.Email)

	existing, err := s.repository.FindByEmail(userDto.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("Failed to fetch user", "email", userDto.Email)
		return nil, service_errors.NewErrInternalServer("Failed to fetch user by email")
	}

	if existing != nil {
		logger.Logger.Error("User with given email already exists", "email", userDto.Email)
		return nil, service_errors.NewErrBadRequest("User with this email already exists")
	}

	if userDto.Role == "" {
		userDto.Role = RoleUser
	}

	user := ToUserModel(userDto)
	err = s.repository.Save(user)
	if err != nil {
		logger.Logger.Error("Failed to save user", "email", user.Email)
		return nil, service_errors.NewErrInternalServer("Failed to save user")
	}

	logger.Logger.Info("Saved user.", "id", user.Id)
	return ToUserDto(user), nil
}

func (s *UserServiceImpl) UpdateUser(id string, userDto *UserDto) (*UserDto, error) {
	logger.Logger.Info("Updating user...", "id", id, "email", userDto.Email)

	existingUser, err := s.repository.FindById(id)
	if err != nil {
		logger.Logger.Error("Failed to find user", "id", id)
		return nil, service_errors.NewErrInternalServer("Failed to find user")
	}

	if userDto.Email != "" {
		existingUser.Email = userDto.Email
	}
	if userDto.Username != "" {
		existingUser.Username = userDto.Username
	}
	if userDto.PictureUrl != "" {
		existingUser.PictureURL = userDto.PictureUrl
	}
	if userDto.Provider != "" {
		existingUser.Provider = userDto.Provider
	}
	if userDto.Status != "" {
		existingUser.Status = userDto.Status
	}
	if userDto.Role != "" {
		existingUser.Role = userDto.Role
	}
	if userDto.Password != "" {
		existingUser.Password = userDto.Password
	}

	if err := s.repository.Save(existingUser); err != nil {
		logger.Logger.Error("Failed to update user", "id", existingUser.Id)
		return nil, service_errors.NewErrInternalServer("Failed to update user")
	}

	logger.Logger.Info("Updated user successfully.", "id", id)
	updatedDto := ToUserDto(existingUser)
	return updatedDto, nil
}

func (s *UserServiceImpl) UpsertUser(userDto *UserDto) (*UserDto, error) {
	logger.Logger.Info("Updating user...", "email", userDto.Email)

	existing, err := s.repository.FindByEmail(userDto.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("Failed to fetch user", "email", userDto.Email)
		return nil, service_errors.NewErrInternalServer("Failed to fetch user by email")
	}

	if existing != nil {
		return s.UpdateUser(existing.Id, userDto)
	} else {
		return s.SaveUser(userDto)
	}
}

func (s *UserServiceImpl) FindById(id string) (*UserDto, error) {
	logger.Logger.Info("Fetching user by id...", "id", id)

	user, err := s.repository.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("User not found", "id", id)
			return nil, service_errors.NewErrNotFound("User not found")
		}

		logger.Logger.Error("Failed to fetch user", "id", id)
		return nil, service_errors.NewErrInternalServer("Failed to fetch user by id")
	}

	logger.Logger.Info("Fetched user by id", "id", user.Id)
	return ToUserDto(user), nil
}

func (s *UserServiceImpl) FindByEmail(email string) (*UserDto, error) {
	logger.Logger.Info("Fetching user by email...", "email", email)

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("User not found", "email", email)
			return nil, service_errors.NewErrNotFound("User not found")
		}

		logger.Logger.Error("Failed to fetch user", "email", email)
		return nil, service_errors.NewErrInternalServer("Failed to fetch user by email")
	}

	logger.Logger.Info("Fetched user by email.", "id", user.Id)
	return ToUserDto(user), nil
}
