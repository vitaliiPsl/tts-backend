package auth

import (
	"errors"
	"os"
	"time"
	"vitaliiPsl/synthesizer/internal/auth/jwt"
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
	jwtService   *jwt.JwtService
}

func NewAuthService(userService *user_package.UserService, tokenService *token.TokenService, emailService *email.EmailService, jwtService *jwt.JwtService) *AuthService {
	emailVerificationUrl := os.Getenv("EMAIL_VERIFICATION_URL")

	return &AuthService{
		emailVerificationUrl: emailVerificationUrl,
		userService:          userService,
		tokenService:         tokenService,
		emailService:         emailService,
		jwtService:           jwtService,
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

func (s *AuthService) HandleSignIn(req *requests.SignInRequest) (string, error) {
	logger.Logger.Info("Handling sing in req", "email", req.Email)

	user, err := s.userService.FindByEmail(req.Email)
	if err != nil {
		var errNotFound *service_errors.ErrNotFound
		if !errors.As(err, &errNotFound) {
			logger.Logger.Error("User with given email doesn't exist", "email", req.Email)
			return "", &service_errors.ErrUnauthorized{Message: "Invalid username or password"}
		}

		logger.Logger.Error("Failed to fetch user", "email", req.Email)
		return "", err
	}

	if user.Status != user_package.StatusActive {
		logger.Logger.Error("User is not active", "email", req.Email, "status", user.Status)
		return "", &service_errors.ErrUnauthorized{Message: "Email not verified"}

	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Logger.Error("Incorrect password", "email", req.Email)
		return "", &service_errors.ErrUnauthorized{Message: "Invalid username or password"}
	}

	jwtToken, err := s.jwtService.GenerateJWT(user)
	if err != nil {
		return "", err
	}

	logger.Logger.Info("Handled sign in.")
	return jwtToken, nil
}

func (s *AuthService) HandleEmailVerification(req *requests.EmailVerificationRequest) error {
	logger.Logger.Info("Handling email verification", "token", req.Token)

	verificationToken, err := s.tokenService.GetToken(req.Token)
	if err != nil {
		return err
	}

	if verificationToken.Purpose != token.PurposeEmailVerification {
		logger.Logger.Error("Invalid verification token purpose", "token", req.Token, "purpose", verificationToken.Purpose)
		return &service_errors.ErrBadRequest{Message: "Invalid token purpose"}
	}

	if time.Now().After(verificationToken.ExpiresAt) {
		logger.Logger.Error("Email token expired", "token", req.Token, "expiredAt", verificationToken.ExpiresAt)
		return &service_errors.ErrBadRequest{Message: "Email token expired"}
	}

	user, err := s.userService.FindById(verificationToken.UserID)
	if err != nil {
		return err
	}

	user.Status = user_package.StatusActive
	s.userService.UpdateUser(user.Id, user)

	err = s.tokenService.DeleteTokensForUser(verificationToken.UserID)
	if err != nil {
		return err
	}

	logger.Logger.Info("Verified email address", "userId", verificationToken.UserID)
	return nil
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
