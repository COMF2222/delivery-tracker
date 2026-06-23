package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load .env file: %w", err)
	}

	dbHost, err := getRequiredEnv("POSTGRES_HOST")
	if err != nil {
		return Config{}, err
	}
	dbPort, err := getRequiredEnv("POSTGRES_PORT")
	if err != nil {
		return Config{}, err
	}
	dbName, err := getRequiredEnv("POSTGRES_DB")
	if err != nil {
		return Config{}, err
	}
	dbUser, err := getRequiredEnv("POSTGRES_USER")
	if err != nil {
		return Config{}, err
	}
	dbPassword, err := getRequiredEnv("POSTGRES_PASSWORD")
	if err != nil {
		return Config{}, err
	}

	serverPort, err := getRequiredEnv("SERVER_PORT")
	if err != nil {
		return Config{}, err
	}

	jwtSecret, err := getRequiredEnv("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}
	jwtTTL, err := getRequiredEnv("JWT_TTL")
	if err != nil {
		return Config{}, err
	}

	db := DatabaseConfig{
		Host:     dbHost,
		Port:     dbPort,
		Name:     dbName,
		User:     dbUser,
		Password: dbPassword,
	}

	server := ServerConfig{Port: serverPort}

	parseTTL, err := time.ParseDuration(jwtTTL)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse duration: %w", err)
	}

	jwt := JWTConfig{Secret: jwtSecret, TTL: parseTTL}

	config := Config{
		Database: db,
		Server:   server,
		JWT:      jwt,
	}

	return config, nil
}

func getRequiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return value, fmt.Errorf("env variable %s is required", key)
	}
	return value, nil
}
