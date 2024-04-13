package jwt

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/user"
)

type JwtService struct {
	application     string
	secretKey       string
	expirationHours int
}

func NewJwtService() *JwtService {
	appName := os.Getenv("APP_NAME")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	expirationHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil {
		logger.Logger.Error("Invalid JWT_EXPIRATION_HOURS value.", "value", expirationHours)
		panic(fmt.Sprintf("Invalid JWT_EXPIRATION_HOURS value: %v", expirationHours))
	}

	return &JwtService{
		application:     appName,
		secretKey:       jwtSecretKey,
		expirationHours: expirationHours,
	}
}

func (s *JwtService) GenerateJWT(user *user.UserDto) (string, error) {
	logger.Logger.Info("Generating JWT token...", "userId", user.Id)

	claims := s.createUserClaims(user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		logger.Logger.Error("Failed to sign JWT")
		return "", &service_errors.ErrInternalServer{}
	}

	logger.Logger.Info("Generated JWT token.")
	return signedToken, nil
}

func (s *JwtService) ValidateToken(tokenString string) (*UserClaims, error) {
	logger.Logger.Info("Validating JWT token...")

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Logger.Error("Unexpected signing method", "method", token.Header["alg"])
			return nil, &service_errors.ErrInternalServer{}
		}
		return []byte(s.secretKey), nil
	})

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		logger.Logger.Info("Token passed validation.")
		return claims, nil
	} else {
		logger.Logger.Info("Token failed validation.")
		return nil, err
	}
}

func (s *JwtService) createUserClaims(user *user.UserDto) *UserClaims {
	return &UserClaims{
		Id: user.Id,
		StandardClaims: jwt.StandardClaims{
			Issuer:    s.application,
			Subject:   user.Id,
			ExpiresAt: time.Now().Add(time.Duration(s.expirationHours) * time.Hour).Unix(),
		},
	}
}
