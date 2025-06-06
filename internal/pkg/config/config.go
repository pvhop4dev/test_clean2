package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Rate     RateLimitConfig
}

type AppConfig struct {
	Name   string
	Env    string
	Port   string
	Secret string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type RateLimitConfig struct {
	Limit int
	Burst int
}

func LoadConfig() *Config {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(b))))

	viper.SetConfigFile(basepath + "/.env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file", err)
	}

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
			DBName:   viper.GetString("DB_NAME"),
		},
		Redis: RedisConfig{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		Rate: RateLimitConfig{
			Limit: viper.GetInt("RATE_LIMIT"),
			Burst: viper.GetInt("RATE_BURST"),
		},
	}

	return config
}
