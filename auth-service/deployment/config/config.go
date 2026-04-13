package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string

	JWTSecret          string
	JWTExpirationHours int
}

func Load() (*Config, error) {
	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./deployment/config")
	viper.AddConfigPath(".")

	// Читаем config.yaml
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("No config.yaml found")
	}

	cfg := &Config{
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpirationHours: viper.GetInt("jwt.expiration_hours"),
		ServerPort:         viper.GetString("server.port"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.JWTExpirationHours <= 0 {
		return nil, fmt.Errorf("JWT_EXPIRATION_HOURS must be positive")
	}

	return cfg, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}
