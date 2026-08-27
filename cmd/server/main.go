package main

import (
	"HailowSellerService/internal/domain"
	"HailowSellerService/internal/transport/grpc/server"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			slog.Warn(fmt.Sprintf("Failed to load .env file: %v", err))
		}
	}

	debug := false
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	host := os.Getenv("SELLER_SERVICE_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("SELLER_SERVICE_PORT")
	if portStr == "" {
		portStr = "50002"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		slog.Error(fmt.Sprintf("Invalid SELLER_SERVICE_PORT value '%s': %v", portStr, err))
	}

	srv, err := server.New(debug, host, port)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: %v", domain.ErrInitServer.Message, err))
	}

	if err := srv.Run(); err != nil {
		slog.Error(fmt.Sprintf("%s: %v", domain.ErrRunServer.Message, err))
	}
}
