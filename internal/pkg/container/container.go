package container

import (
	"clean-arch-go/internal/domain/repository"
	"clean-arch-go/internal/domain/service"
	"clean-arch-go/internal/pkg/config"
	"clean-arch-go/internal/pkg/database"
	"clean-arch-go/internal/pkg/redis"
)

type Container struct {
	DB             *database.Database
	RedisClient    *redis.RedisClient
	Config         *config.Config
	AuthSvc        service.AuthService
	BookSvc        service.BookService
	TranslationSvc service.TranslationService
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// Initialize database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		return nil, err
	}

	// Initialize Redis
	redisClient, err := redis.NewRedisClient(&cfg.Redis)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	translationRepo := repository.NewTranslationRepository(db)

	// Initialize services
	authSvc := service.NewAuthService(
		userRepo,
		cfg.App.Secret,
		redisClient,
	)
	bookSvc := service.NewBookService(bookRepo)
	translationSvc := service.NewTranslationService(translationRepo)

	return &Container{
		DB:             db,
		RedisClient:    redisClient,
		Config:         cfg,
		AuthSvc:        authSvc,
		BookSvc:        bookSvc,
		TranslationSvc: translationSvc,
	}, nil
}

func (c *Container) Close() error {
	// Đóng kết nối Redis
	if err := c.RedisClient.Close(); err != nil {
		return err
	}

	// Đóng kết nối database
	if err := c.DB.Close(); err != nil {
		return err
	}
	return nil
}
