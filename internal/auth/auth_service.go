package auth

import (
	"errors"
	"os"
	"vitaliiPsl/synthesizer/internal/email"
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/requests"
	"vitaliiPsl/synthesizer/internal/token"
	user_package "vitaliiPsl/synthesizer/internal/user"

	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLength = 8

type AuthService struct {
	emailVerificationUrl string

	userService  *user_package.UserService
	tokenService *token.TokenService
	emailService *email.EmailService
}

func NewAuthService(userService *user_package.UserService, tokenService *token.TokenService, emailService *email.EmailService) *AuthService {
	emailVerificationUrl := os.Getenv("EMAIL_VERIFICATION_URL")

	return &AuthService{
		emailVerificationUrl: emailVerificationUrl,
		userService:          userService,
		tokenService:         tokenService,
		emailService:         emailService,
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

	savedUser, err := s.userService.SaveUser(userDto)
	if err != nil {
		return err
	}

	return s.sendVerificationEmail(savedUser)
}

func (s *AuthService) sendVerificationEmail(user *user_package.UserDto) error {
	token, err := s.tokenService.CreateVerificationToken(user.Id, token.PurposeEmailVerification)
	if err != nil {
		return err
	}

	emailVariables := map[string]string{
		"user_name":         user.Username,
		"verification_link": s.emailVerificationUrl + token.Token,
	}

	return s.emailService.SendTemplatedEmail(user.Email, "Email verification", "email_verification.html", emailVariables)
}
