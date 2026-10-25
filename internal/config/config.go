package config

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	Port        string
	SessionKey  string
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("error loading .env file, not found, reading system variables")
	}

	cfg := Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", ""),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", ""),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", ""),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		Port:        getEnv("PORT", "8080"),
		SessionKey:  getEnv("SESSION_KEY", ""),
	}

	if cfg.DatabaseURL == "" && (cfg.DBHost == "" || cfg.DBName == "" || cfg.DBUser == "" || cfg.DBPassword == "") {
		return cfg, fmt.Errorf("missing environment variables for database")
	}

	if cfg.SessionKey == "" {
		return cfg, fmt.Errorf("environment variable SESSION_KEY is required")
	}

	return cfg, nil
}

func (c *Config) DBUrl() string {

	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   net.JoinHostPort(c.DBHost, c.DBPort),
		Path:   c.DBName,
	}

	q := u.Query()
	q.Set("sslmode", c.DBSSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
