package token

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
)

type TokenService struct {
	tokenDurationMins int
	repository        *TokenRepository
}

func NewTokenService(repo *TokenRepository) *TokenService {
	tokenDuration, err := strconv.Atoi(os.Getenv("VERIFICATION_TOKEN_DURATION_MIN"))
	if err != nil {
		logger.Logger.Error("Invalid verification token duration", "error", err)
		panic(1)
	}

	return &TokenService{tokenDurationMins: tokenDuration, repository: repo}
}

func (s *TokenService) CreateVerificationToken(userId string, purpose TokenPurpose) (*TokenDto, error) {
	logger.Logger.Info("Creating new verification token", "userId", userId)

	expiration := time.Now().Add(time.Minute * time.Duration(s.tokenDurationMins))
	token := uuid.NewString()

	verificationToken := Token{
		UserID:    userId,
		Token:     token,
		Purpose:   purpose,
		ExpiresAt: expiration,
	}

	if err := s.repository.Save(&verificationToken); err != nil {
		logger.Logger.Error("Failed to save verification token", "userId", userId, "error", err)
		return nil, &service_errors.ErrInternalServer{Message: "Failed to save verification token"}
	}

	logger.Logger.Info("Created new verification token", "userId", userId, "tokenId", verificationToken.Id)
	return ToVerificationTokenDto(&verificationToken), nil
}

func (s *TokenService) GetToken(token string) (*TokenDto, error) {
	logger.Logger.Info("Fetching token model", "token", token)

	verificationToken, err := s.repository.FindByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("Verification token not found", "token", token, "error", err)
			return nil, &service_errors.ErrNotFound{Message: "Token not found"}
		}

		logger.Logger.Error("Failed to fetch verification token", "token", token, "error", err)
		return nil, &service_errors.ErrInternalServer{Message: "Failed to fetch verification token"}
	}

	logger.Logger.Info("Fetched token model")
	return ToVerificationTokenDto(verificationToken), nil
}

func (s *TokenService) DeleteTokensForUser(userId string) error {
	logger.Logger.Info("Deleting user's tokens", "userId", userId)

	if err := s.repository.DeleteByUserID(userId); err != nil {
		logger.Logger.Error("Failed to delete verification token", "userId", userId, "error", err)
		return &service_errors.ErrInternalServer{Message: "Failed to delete user tokens"}
	}

	logger.Logger.Info("Deleted user's tokens", "userId", userId)
	return nil
}
