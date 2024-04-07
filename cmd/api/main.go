package main

import (
	"fmt"
	"os"
	"strconv"

	"vitaliiPsl/synthesizer/internal/auth"
	"vitaliiPsl/synthesizer/internal/database"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/router"
	"vitaliiPsl/synthesizer/internal/server"
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

	validationService := validation.NewValidationService()

	authenticationService := auth.NewAuthService(userService)
	authenticationControler := auth.NewAuthController(authenticationService, validationService)

	router.SetupRoutes(server.App, authenticationControler)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
