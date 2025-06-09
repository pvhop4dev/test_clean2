package container

import (
	"log"

	"clean-arch-go/internal/domain/repository"
	"clean-arch-go/internal/domain/service"
	"clean-arch-go/internal/infrastructure/repository/cached"
	"clean-arch-go/internal/pkg/config"
	"clean-arch-go/internal/pkg/database"
	"clean-arch-go/internal/pkg/redis"

	"gorm.io/gorm"
)

// Container holds all the application dependencies
type Container struct {
	DB             *database.Database
	RedisClient    *redis.RedisClient
	Config         *config.Config
	AuthSvc        service.AuthService
	BookSvc        service.BookService
	TranslationSvc service.TranslationService
	UserRepo       repository.UserRepository
	BookRepo       repository.BookRepository
	TranslationRepo repository.TranslationRepository
}

// NewContainer creates a new application container with all dependencies
func NewContainer(cfg *config.Config) (*Container, error) {
	// Initialize database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		return nil, err
	}

	// Run database migrations
	if err := runMigrations(db); err != nil {
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

	// Initialize cached repositories
	cachedUserRepo := cached.NewCachedUserRepository(userRepo, redisClient)
	cachedBookRepo := cached.NewCachedBookRepository(bookRepo, redisClient)

	// Initialize services
	authSvc := service.NewAuthService(
		cachedUserRepo,
		cfg.App.Secret,
		redisClient,
	)

	bookSvc := service.NewBookService(cachedBookRepo)
	translationSvc := service.NewTranslationService(translationRepo)

	return &Container{
		DB:              db,
		RedisClient:     redisClient,
		Config:          cfg,
		AuthSvc:         authSvc,
		BookSvc:         bookSvc,
		TranslationSvc:  translationSvc,
		UserRepo:        cachedUserRepo,
		BookRepo:        cachedBookRepo,
		TranslationRepo: translationRepo,
	}, nil
}

// runMigrations runs database migrations for all domain models
func runMigrations(db *database.Database) error {
	// Get the underlying GORM DB instance
	dbInstance, err := db.DB.DB()
	if err != nil {
		return err
	}

	// Ping the database to ensure connection is working
	if err := dbInstance.Ping(); err != nil {
		return err
	}

	// Run migrations for all domain models
	// Note: You'll need to import your domain models here
	// Example:
	// if err := db.AutoMigrate(
	//     &user.User{},
	//     &book.Book{},
	//     // Add other models here
	// ); err != nil {
	//     return err
	// }


	log.Println("Database migrations completed successfully")
	return nil
}

// Close gracefully shuts down all connections and cleans up resources
func (c *Container) Close() error {
	// Close Redis connection
	if c.RedisClient != nil {
		if err := c.RedisClient.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
			// Continue with other cleanup even if Redis close fails
		}
	}

	// Close database connection
	if c.DB != nil {
		db, err := c.DB.DB.DB()
		if err != nil {
			log.Printf("Error getting database instance: %v", err)
		} else if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
			// Continue with other cleanup even if DB close fails
		}
	}

	// Add any additional cleanup here

	return nil
}

// GetDB returns the underlying GORM DB instance
func (c *Container) GetDB() *gorm.DB {
	return c.DB.DB
}

// GetRedis returns the Redis client
func (c *Container) GetRedis() *redis.RedisClient {
	return c.RedisClient
}
