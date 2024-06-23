package main

import (
	"fmt"
	"os"
	"strconv"

	"vitaliiPsl/synthesizer/internal/auth"
	"vitaliiPsl/synthesizer/internal/auth/jwt"
	"vitaliiPsl/synthesizer/internal/auth/sso"
	"vitaliiPsl/synthesizer/internal/database"
	"vitaliiPsl/synthesizer/internal/email"
	"vitaliiPsl/synthesizer/internal/history"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/model"
	"vitaliiPsl/synthesizer/internal/router"
	"vitaliiPsl/synthesizer/internal/server"
	"vitaliiPsl/synthesizer/internal/synthesis"
	"vitaliiPsl/synthesizer/internal/token"
	"vitaliiPsl/synthesizer/internal/users"
	"vitaliiPsl/synthesizer/internal/validation"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger.Logger.Info("Spinning up Synthesizer...")

	database.SetupDatabase()

	server := server.New()

	userRepository := users.NewUserRepository(database.DB)
	userService := users.NewUserService(userRepository)

	tokenRepository := token.NewTokenRepository(database.DB)
	tokenService := token.NewTokenService(tokenRepository)

	emailService := email.NewEmailService()

	validationService := validation.NewValidationService()

	jwtService := jwt.NewJwtService()

	githubConfig := sso.GithubSSOConfig()
	githubProvider := sso.NewGithubProvider(githubConfig)
	ssoProviders := map[string]sso.SSOProvider{"github": githubProvider}

	authenticationService := auth.NewAuthService(userService, tokenService, emailService, jwtService, ssoProviders)
	authenticationControler := auth.NewAuthController(authenticationService, validationService)
	authenticationMiddleware := auth.NewAuthMiddleware(jwtService, userService)

	modelRepository := model.NewModelRepository(database.DB)
	modelService := model.NewModelService(modelRepository)
	modelController := model.NewModelController(modelService, validationService)

	historyRepository := history.NewHistoryRepository(database.DB)
	historyService := history.NewHistoryService(historyRepository)
	historyController := history.NewHistoryController(historyService)

	synthesisService := synthesis.NewSynthesisService(modelService, historyService)
	synthesisController := synthesis.NewSynthesisController(synthesisService, validationService)

	router.SetupRoutes(server.App, authenticationMiddleware, authenticationControler, modelController, synthesisController, historyController)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
