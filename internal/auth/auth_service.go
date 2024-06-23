package auth

import (
	"errors"
	"os"
	"time"
	"vitaliiPsl/synthesizer/internal/auth/jwt"
	"vitaliiPsl/synthesizer/internal/auth/sso"
	"vitaliiPsl/synthesizer/internal/email"
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/requests"
	"vitaliiPsl/synthesizer/internal/token"
	"vitaliiPsl/synthesizer/internal/users"

	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLength = 8

type AuthService struct {
	emailVerificationUrl string
	passwordResetUrl     string

	userService  users.UserService
	tokenService token.TokenService
	emailService email.EmailService
	jwtService   jwt.JwtService
	providers    map[string]sso.SSOProvider
}

func NewAuthService(
	userService users.UserService,
	tokenService token.TokenService,
	emailService email.EmailService,
	jwtService jwt.JwtService,
	providers map[string]sso.SSOProvider,
) *AuthService {
	emailVerificationUrl := os.Getenv("EMAIL_VERIFICATION_URL")
	passwordResetUrl := os.Getenv("PASSWORD_RESET_URL")

	return &AuthService{
		emailVerificationUrl: emailVerificationUrl,
		passwordResetUrl:     passwordResetUrl,
		userService:          userService,
		tokenService:         tokenService,
		emailService:         emailService,
		jwtService:           jwtService,
		providers:            providers,
	}
}

func (s *AuthService) HandleSignUp(req *requests.SignUpRequest) error {
	logger.Logger.Info("Handling sign up", "email", req.Email)

	existingUser, err := s.userService.FindByEmail(req.Email)
	if err != nil {
		var errNotFound *service_errors.ErrNotFound
		if !errors.As(err, &errNotFound) {
			logger.Logger.Error("User with given email already exists", "email", req.Email)
			return service_errors.NewErrBadRequest("User with given email already exists")
		}
	} else if existingUser != nil {
		logger.Logger.Error("User with given email already exists", "email", req.Email)
		return service_errors.NewErrBadRequest("User with given email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error("Failed to hash password", "email", req.Email)
		return err
	}

	userDto := &users.UserDto{
		Email:    req.Email,
		Password: string(hashedPassword),
		Username: req.Username,
		Role:     users.RoleUser,
		Status:   users.StatusPending,
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
			return "", service_errors.NewErrUnauthorized("Invalid username or password")
		}

		logger.Logger.Error("Failed to fetch user", "email", req.Email)
		return "", err
	}

	if user.Status != users.StatusActive {
		logger.Logger.Error("User is not active", "email", req.Email, "status", user.Status)
		return "", service_errors.NewErrUnauthorized("Email not verified")

	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Logger.Error("Incorrect password", "email", req.Email)
		return "", service_errors.NewErrUnauthorized("Invalid username or password")
	}

	jwtToken, err := s.jwtService.GenerateJWT(user)
	if err != nil {
		return "", err
	}

	logger.Logger.Info("Handled sign in.")
	return jwtToken, nil
}

func (s *AuthService) HandleSsoSignIn(providerName string) (string, error) {
	logger.Logger.Info("Handling SSO sign in", "provider", providerName)

	provider, exists := s.providers[providerName]
	if !exists {
		return "", service_errors.NewErrBadRequest("Unsupported SSO provider")
	}

	return provider.AuthCodeURL("state-string"), nil
}

func (s *AuthService) HandleSSOCallback(providerName, code string) (string, error) {
	logger.Logger.Info("Handling SSO callback", "provider", providerName)

	provider, exists := s.providers[providerName]
	if !exists {
		return "", service_errors.NewErrBadRequest("Unsupported SSO provider")
	}

	token, err := provider.Exchange(code)
	if err != nil {
		return "", err
	}

	var user *users.UserDto
	user, err = provider.FetchUserInfo(token)
	if err != nil {
		return "", err
	}

	user.Provider = providerName
	user.Status = users.StatusActive

	user, err = s.userService.UpsertUser(user)
	if err != nil {
		return "", err
	}

	var jwtToken string
	jwtToken, err = s.jwtService.GenerateJWT(user)
	if err != nil {
		return "", err
	}

	logger.Logger.Info("Handled SSO sign in.")
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
		return service_errors.NewErrBadRequest("Invalid token purpose")
	}

	if time.Now().After(verificationToken.ExpiresAt) {
		logger.Logger.Error("Email token expired", "token", req.Token, "expiredAt", verificationToken.ExpiresAt)
		return service_errors.NewErrBadRequest("Email token expired")
	}

	user, err := s.userService.FindById(verificationToken.UserID)
	if err != nil {
		return err
	}

	user.Status = users.StatusActive
	s.userService.UpdateUser(user.Id, user)

	err = s.tokenService.DeleteTokensForUser(verificationToken.UserID)
	if err != nil {
		return err
	}

	logger.Logger.Info("Verified email address", "userId", verificationToken.UserID)
	return nil
}

func (s *AuthService) HandlePasswordReset(req *requests.PasswordResetRequest) error {
	logger.Logger.Info("Handling password reset", "token", req.Token)

	verificationToken, err := s.tokenService.GetToken(req.Token)
	if err != nil {
		return err
	}

	if verificationToken.Purpose != token.PurposePasswordReset {
		logger.Logger.Error("Invalid verification token purpose", "token", req.Token, "purpose", verificationToken.Purpose)
		return service_errors.NewErrBadRequest("Invalid token purpose")
	}

	if time.Now().After(verificationToken.ExpiresAt) {
		logger.Logger.Error("Password reset token expired", "token", req.Token, "expiredAt", verificationToken.ExpiresAt)
		return service_errors.NewErrBadRequest("Password token expired")
	}

	user, err := s.userService.FindById(verificationToken.UserID)
	if err != nil {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error("Failed to hash password", "token", req.Token)
		return err
	}

	user.Password = string(hashedPassword)
	_, err = s.userService.UpdateUser(user.Id, user)
	if err != nil {
		return err
	}

	err = s.tokenService.DeleteTokensForUser(verificationToken.UserID)
	if err != nil {
		return err
	}

	logger.Logger.Info("Reset password", "userId", verificationToken.UserID)
	return nil
}

func (s *AuthService) HandleSendPasswordResetToken(req *requests.VerificationTokenRequest) error {
	logger.Logger.Info("Handling resend of password verification token", "email", req.Email)

	user, err := s.userService.FindByEmail(req.Email)
	if err != nil {
		return err
	}

	if user == nil {
		logger.Logger.Error("User not found", "email", req.Email)
		return service_errors.NewErrNotFound("User not found")
	}

	if err = s.sendResetPasswordEmail(user); err != nil {
		return err
	}

	logger.Logger.Info("Handled resend of password verification token", "email", req.Email)
	return nil
}

func (s *AuthService) sendVerificationEmail(user *users.UserDto) error {
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

func (s *AuthService) sendResetPasswordEmail(user *users.UserDto) error {
	token, err := s.tokenService.CreateVerificationToken(user.Id, token.PurposePasswordReset)
	if err != nil {
		return err
	}

	emailVariables := map[string]string{
		"password_reset_link": s.passwordResetUrl + token.Token,
	}

	return s.emailService.SendTemplatedEmail(user.Email, "Password reset", "reset_password.html", emailVariables)
}
