package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	FrontendOrigin  string
	UploadDir       string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	loadEnvFiles(".env", filepath.Join("..", ".env"))

	cfg := Config{
		Port:            appPort(),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		JWTSecret:       env("JWT_SECRET", "change-me-in-env"),
		FrontendOrigin:  env("FRONTEND_ORIGIN", "http://localhost:5173"),
		UploadDir:       env("UPLOAD_DIR", "uploads"),
		AccessTokenTTL:  envDuration("ACCESS_TOKEN_TTL_MINUTES", 30),
		RefreshTokenTTL: envDuration("REFRESH_TOKEN_TTL_HOURS", 24*7),
		ReadTimeout:     envDuration("HTTP_READ_TIMEOUT_SECONDS", 10),
		WriteTimeout:    envDuration("HTTP_WRITE_TIMEOUT_SECONDS", 15),
		ShutdownTimeout: envDuration("HTTP_SHUTDOWN_TIMEOUT_SECONDS", 10),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func loadEnvFiles(paths ...string) {
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			if key != "" && os.Getenv(key) == "" {
				_ = os.Setenv(key, value)
			}
		}
		_ = file.Close()
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func appPort() string {
	if value := os.Getenv("APP_PORT"); value != "" {
		return value
	}
	if value := os.Getenv("PORT"); value != "" {
		return value
	}
	return "8080"
}

func envDuration(key string, fallback int) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(fallback) * time.Minute
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return time.Duration(fallback) * time.Minute
	}

	switch key {
	case "REFRESH_TOKEN_TTL_HOURS":
		return time.Duration(number) * time.Hour
	case "HTTP_READ_TIMEOUT_SECONDS", "HTTP_WRITE_TIMEOUT_SECONDS", "HTTP_SHUTDOWN_TIMEOUT_SECONDS":
		return time.Duration(number) * time.Second
	default:
		return time.Duration(number) * time.Minute
	}
}
