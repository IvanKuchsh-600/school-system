// package config

// import (
// 	"os"
// 	"strconv"

// 	"github.com/spf13/viper"
// )

// type Config struct {
// 	Server ServerConfig
// 	JWT    JWTConfig
// }

// type ServerConfig struct {
// 	Port string
// }

// type JWTConfig struct {
// 	Secret          string
// 	ExpirationHours int
// }

// func Load() (*Config, error) {
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("yaml")
// 	viper.AddConfigPath("./deployment/config")

// 	// Читаем YAML если есть (если нет - не страшно)
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		// Handle error appropriately
// 	}

// 	// Берем значения (с приоритетом: env > yaml > default)
// 	cfg := &Config{
// 		Server: ServerConfig{
// 			Port: getEnv("SERVER_PORT", viper.GetString("server.port")),
// 		},
// 		JWT: JWTConfig{
// 			Secret:          getEnv("JWT_SECRET", viper.GetString("jwt.secret")),
// 			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", viper.GetInt("jwt.expiration_hours")),
// 		},
// 	}

// 	return cfg, nil
// }

// func getEnv(key, defaultValue string) string {
// 	if value := os.Getenv(key); value != "" {
// 		return value
// 	}
// 	return defaultValue
// }

// func getEnvAsInt(key string, defaultValue int) int {
// 	if value := os.Getenv(key); value != "" {
// 		if intVal, err := strconv.Atoi(value); err == nil {
// 			return intVal
// 		}
// 	}
// 	return defaultValue
// }

package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort         string
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
