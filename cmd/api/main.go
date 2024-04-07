package main

import (
	"fmt"
	"os"
	"strconv"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/server"
)

func main() {
	logger.Logger.Info("Spinning up Synthesizer...")

	server := server.New()

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
