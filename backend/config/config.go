package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/meh-hackathon/meh/apperror"
	"github.com/meh-hackathon/meh/logger"
)

var (
	ErrConfigValueNotSet  = apperror.Define("config:value_not_set", "Configuration value not set")
	ErrInvalidConfigValue = apperror.Define("config:invalid_value", "Invalid configuration value")
	ErrReadingConfigFile  = apperror.Define("config:reading_file", "Error reading configuration file")
)

type environment string

const (
	Prod environment = "prod"
	Dev  environment = "dev"
)

var (
	Env         environment
	Secret      []byte
	Port        int
	LogLevel    slog.Level
	DatabaseURL string
)

func Load(envFiles ...string) error {
	if len(envFiles) > 0 {
		err := godotenv.Load(envFiles...)
		if err != nil {
			return ErrReadingConfigFile.WithMessage("Failed to load environment variables").WithOrigin().WithCause(err)
		}
	} else {
		err := godotenv.Load(".env.local", ".env")
		if err != nil && !os.IsNotExist(err) {
			return ErrReadingConfigFile.WithMessage("Failed to load environment variables").WithOrigin().WithCause(err)
		}
	}

	env := fallback(os.Getenv("ENV"), "dev")
	secret := os.Getenv("SECRET")
	databaseURL := os.Getenv("DATABASE_URL")
	portStr := fallback(os.Getenv("PORT"), "8080")
	logLevelStr := fallback(os.Getenv("LOG_LEVEL"), "INFO")

	err := errors.Join(
		mustHaveMinLen(secret, 32, "SECRET"),
		mustBeSet(databaseURL, "DATABASE_URL"),
		mustBeInt(portStr, "PORT"),
		mustBeValidLogLevel(logLevelStr, "LOG_LEVEL"),
		mustBeOneOf(env, []string{"prod", "dev"}, "ENV"),
	)
	if err != nil {
		return err
	}

	Env = environment(env)
	Secret = []byte(secret)
	DatabaseURL = databaseURL
	Port, _ = strconv.Atoi(portStr)
	LogLevel, _ = logger.LevelFromString(logLevelStr)

	return nil
}

func fallback(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func mustBeSet(value, name string) error {
	if value == "" {
		return ErrConfigValueNotSet.WithMessage(fmt.Sprintf("%s must be set", name)).WithOrigin()
	}
	return nil
}

func mustBeInt(value, name string) error {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return ErrInvalidConfigValue.WithMessage(fmt.Sprintf("%s must be a valid integer", name)).WithOrigin().WithCause(err)
	}
	if parsed <= 0 {
		return ErrInvalidConfigValue.WithMessage(fmt.Sprintf("%s must be a valid positive integer", name)).WithOrigin()
	}
	return nil
}

func mustHaveMinLen(value string, minLen int, name string) error {
	if value == "" {
		return ErrConfigValueNotSet.WithMessage(fmt.Sprintf("%s must be set", name)).WithOrigin()
	}
	if len(value) < minLen {
		return ErrInvalidConfigValue.WithMessage(fmt.Sprintf("%s must be at least %d characters long", name, minLen)).WithOrigin()
	}
	return nil
}

func mustBeOneOf(value string, options []string, name string) error {
	if slices.Contains(options, value) {
		return nil
	}
	return ErrInvalidConfigValue.WithMessage(fmt.Sprintf("%s (%s) must be one of %v", name, value, options)).WithOrigin()
}

func mustBeValidLogLevel(value, name string) error {
	_, err := logger.LevelFromString(value)
	if err != nil {
		return ErrInvalidConfigValue.WithMessage(fmt.Sprintf("%s must be one of DEBUG, INFO, WARNING, ERROR", name)).WithOrigin().WithCause(err)
	}
	return nil
}

func anyIsSet(values ...string) bool {
	for _, value := range values {
		if value != "" {
			return true
		}
	}
	return false
}

func allAreSet(values ...string) bool {
	return !slices.Contains(values, "")
}
