package main

import (
	"fmt"
	"os"
	"strconv"

	"vitaliiPsl/synthesizer/internal/auth"
	"vitaliiPsl/synthesizer/internal/database"
	"vitaliiPsl/synthesizer/internal/email"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/router"
	"vitaliiPsl/synthesizer/internal/server"
	"vitaliiPsl/synthesizer/internal/token"
	"vitaliiPsl/synthesizer/internal/user"
	"vitaliiPsl/synthesizer/internal/validation"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger.Logger.Info("Spinning up Synthesizer...")

	database.SetupDatabase()

	server := server.New()

	userRepository := user.NewUserRepository(database.DB)
	userService := user.NewUserService(userRepository)


	tokenRepository := token.NewTokenRepository(database.DB)
	tokenService := token.NewTokenService(tokenRepository)

	emailService := email.NewEmailService()

	validationService := validation.NewValidationService()

	authenticationService := auth.NewAuthService(userService, tokenService, emailService)
	authenticationControler := auth.NewAuthController(authenticationService, validationService)

	router.SetupRoutes(server.App, authenticationControler)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
