package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"clean-arch-go/internal/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	// In ra DSN để debug (không nên làm điều này trong môi trường production)
	log.Printf("Connecting to database: %s@%s:%s/%s", 
		cfg.User, cfg.Host, cfg.Port, cfg.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting DB instance: %v", err)
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to database")
	return &Database{db}, nil
}

func (db *Database) Migrate(models ...interface{}) error {
	return db.AutoMigrate(models...)
}

// WithContext returns a new DB instance with the given context
func (db *Database) WithContext(ctx context.Context) *gorm.DB {
	return db.DB.WithContext(ctx)
}
