package handler

import (
	"net/http"

	"clean-arch-go/internal/domain/service"
	"clean-arch-go/internal/pkg/redis"
	"clean-arch-go/internal/pkg/server/http/httpconfig"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authSvc      service.AuthService
	bookSvc      service.BookService
	translationSvc service.TranslationService
	redisClient  *redis.RedisClient
	HTTPConfig   *httpconfig.HTTPConfig
}

func NewHandler(
	authSvc service.AuthService,
	bookSvc service.BookService,
	translationSvc service.TranslationService,
	redisClient *redis.RedisClient,
	HTTPConfig *httpconfig.HTTPConfig,
) *Handler {
	return &Handler{
		authSvc:      authSvc,
		bookSvc:      bookSvc,
		translationSvc: translationSvc,
		redisClient:  redisClient,
		HTTPConfig:   HTTPConfig,
	}
}

func (h *Handler) RegisterAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}
}

func (h *Handler) RegisterBookRoutes(router *gin.Engine) {
	books := router.Group("/books")
	{
		books.GET("/", h.ListBooks)
		books.POST("/", h.CreateBook)
		books.GET("/:id", h.GetBook)
	}
}

func (h *Handler) RegisterTranslationRoutes(router *gin.Engine) {
	translations := router.Group("/translations")
	{
		translations.POST("/translate", h.Translate)
		translations.GET("/supported", h.GetSupportedLanguages)
	}
}

func (h *Handler) Login(c *gin.Context) {
	// TODO: Implement login
	c.JSON(http.StatusOK, gin.H{
		"message": "login",
	})
}

func (h *Handler) Register(c *gin.Context) {
	// TODO: Implement register
	c.JSON(http.StatusOK, gin.H{
		"message": "register",
	})
}

func (h *Handler) ListBooks(c *gin.Context) {
	// TODO: Implement list books
	c.JSON(http.StatusOK, gin.H{
		"message": "list books",
	})
}

func (h *Handler) CreateBook(c *gin.Context) {
	// TODO: Implement create book
	c.JSON(http.StatusOK, gin.H{
		"message": "create book",
	})
}

func (h *Handler) GetBook(c *gin.Context) {
	// TODO: Implement get book
	c.JSON(http.StatusOK, gin.H{
		"message": "get book",
	})
}

func (h *Handler) Translate(c *gin.Context) {
	// TODO: Implement translation
	c.JSON(http.StatusOK, gin.H{
		"message": "translate",
	})
}

func (h *Handler) GetSupportedLanguages(c *gin.Context) {
	// TODO: Implement get supported languages
	c.JSON(http.StatusOK, gin.H{
		"message": "supported languages",
	})
}
