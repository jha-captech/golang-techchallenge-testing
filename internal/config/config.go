package config

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config holds the application configuration settings. The configuration is loaded from
// environment variables.
type Config struct {
	DBHost         string     `env:"DATABASE_HOST,required"`
	DBUserName     string     `env:"DATABASE_USER,required"`
	DBUserPassword string     `env:"DATABASE_PASSWORD,required"`
	DBName         string     `env:"DATABASE_NAME,required"`
	DBPort         string     `env:"DATABASE_PORT,required"`
	ClientOrigin   string     `env:"CLIENT_ORIGIN,required"`
	Host           string     `env:"HOST,required"`
	Port           string     `env:"PORT,required"`
	LogLevel       slog.Level `env:"LOG_LEVEL,required"`
	CacheHost      string     `env:"CACHE_HOST,required"`
	CachePort      int        `env:"CACHE_PORT,required"`
	CacheDB        int        `env:"CACHE_DB,required"`
	CachePassword  string     `env:"CACHE_PASSWORD,required"`
}

// New loads configuration from environment variables and a .env file, and returns a
// Config struct or error.
func New() (Config, error) {
	// Load values from a .env file and add them to system environment variables.
	// Discard errors coming from this function. This allows us to call this
	// function without a .env file which will by default load values directly
	// from system environment variables.
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found")
	} else {
		fmt.Println(".env file found")
	}

	// Once values have been loaded into system env vars, parse those into our
	// config struct and validate them returning any errors.
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("[in config.New] failed to parse config: %w", err)
	}

	return cfg, nil
}
