package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig       `mapstructure:",squash"`
	Database  DatabaseConfig  `mapstructure:",squash"`
	Redis     RedisConfig     `mapstructure:",squash"`
	JWT       JWTConfig       `mapstructure:",squash"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type RateLimitConfig struct {
	Limit int `mapstructure:"limit"`
	Burst int `mapstructure:"burst"`
}

type AppConfig struct {
	Name        string
	Env         string
	Port        string
	GRPCPort    string
	Secret      string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret           string
	ExpirationMinute int
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Set default values
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("JWT_EXPIRATION_MINUTE", 1440)
	viper.SetDefault("RATE_LIMIT", 100)
	viper.SetDefault("RATE_BURST", 30)

	// Set default values for app config
	viper.SetDefault("APP_NAME", "Clean Arch Go")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_SECRET", "your-secret-key")

	config := &Config{
		App: AppConfig{
			Name:   viper.GetString("APP_NAME"),
			Env:    viper.GetString("APP_ENV"),
			Port:   viper.GetString("APP_PORT"),
			Secret: viper.GetString("APP_SECRET"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		Redis: RedisConfig{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:           viper.GetString("JWT_SECRET"),
			ExpirationMinute: viper.GetInt("JWT_EXPIRATION_MINUTE"),
		},
		RateLimit: RateLimitConfig{
			Limit: viper.GetInt("RATE_LIMIT"),
			Burst: viper.GetInt("RATE_BURST"),
		},
	}

	return config
}
