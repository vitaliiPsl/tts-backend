package auth

import (
	"errors"
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/requests"
	user_package "vitaliiPsl/synthesizer/internal/user"

	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLength = 8

type AuthService struct {
	userService *user_package.UserService
}

func NewAuthService(userService *user_package.UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (s *AuthService) HandleSignUp(req *requests.SignUpRequest) error {
	logger.Logger.Info("Handling sign up", "email", req.Email)

	existingUser, err := s.userService.FindByEmail(req.Email)
	if err != nil {
		var errNotFound *service_errors.ErrNotFound
		if !errors.As(err, &errNotFound) {
			logger.Logger.Error("User with given email already exists", "email", req.Email)
			return &service_errors.ErrBadRequest{Message: "User with given email already exists"}
		}
	} else if existingUser != nil {
		logger.Logger.Error("User with given email already exists", "email", req.Email)
		return &service_errors.ErrBadRequest{Message: "User with given email already exists"}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error("Failed to hash password", "email", req.Email)
		return err
	}

	userDto := &user_package.UserDto{
		Email:    req.Email,
		Password: string(hashedPassword),
		Username: req.Username,
		Role:     user_package.RoleUser,
		Status:   user_package.StatusPending,
	}

	_, err = s.userService.SaveUser(userDto)
	if err != nil {
		return err
	}

	// TODO: sent email verification link
	return nil
}
