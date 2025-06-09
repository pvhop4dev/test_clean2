package handler

import (
	"net/http"

	"clean-arch-go/internal/domain/service"
	"clean-arch-go/internal/pkg/redis"
	"clean-arch-go/internal/pkg/server/http/httpconfig"

	"github.com/gin-gonic/gin"
)

// BookInput represents the book creation/update request body
// swagger:model BookInput
type BookInput struct {
	// Title of the book
	// required: true
	// example: The Go Programming Language
	Title string `json:"title" binding:"required"`

	// Author of the book
	// required: true
	// example: Alan A. A. Donovan, Brian W. Kernighan
	Author string `json:"author" binding:"required"`

	// Description of the book
	// example: The Go Programming Language is the authoritative resource for any programmer who wants to learn Go.
	Description string `json:"description"`

	// Published year of the book
	// example: 2015
	PublishedYear int `json:"published_year"`
}

// BookResponse represents a book response
// swagger:response bookResponse
type BookResponse struct {
	// ID of the book
	// example: 507f1f77bcf86cd799439011
	ID string `json:"id"`

	// Title of the book
	// example: The Go Programming Language
	Title string `json:"title"`

	// Author of the book
	// example: Alan A. A. Donovan, Brian W. Kernighan
	Author string `json:"author"`

	// Description of the book
	// example: The Go Programming Language is the authoritative resource for any programmer who wants to learn Go.
	Description string `json:"description,omitempty"`

	// Published year of the book
	// example: 2015
	PublishedYear int `json:"published_year,omitempty"`

	// CreatedAt timestamp
	// example: 2023-01-01T00:00:00Z
	CreatedAt string `json:"created_at"`

	// UpdatedAt timestamp
	// example: 2023-01-01T00:00:00Z
	UpdatedAt string `json:"updated_at"`
}

// BooksListResponse represents a list of books response
// swagger:response booksListResponse
type BooksListResponse struct {
	// List of books
	Data []BookResponse `json:"data"`

	// Total number of books
	// example: 42
	Total int `json:"total"`
}

// TranslateInput represents the translation request body
// swagger:model TranslateInput
type TranslateInput struct {
	// Text to translate
	// required: true
	// example: Hello, world!
	Text string `json:"text" binding:"required"`
	// Source language code (ISO 639-1)
	// example: en
	SourceLang string `json:"source_lang"`

	// Target language code (ISO 639-1)
	// required: true
	// example: vi
	TargetLang string `json:"target_lang" binding:"required"`
}

// TranslateResponse represents the translation response
// swagger:response translateResponse
type TranslateResponse struct {
	// Translated text
	// example: Xin chào thế giới!
	Text string `json:"text"`

	// Source language code (ISO 639-1)
	// example: en
	SourceLang string `json:"source_lang"`

	// Target language code (ISO 639-1)
	// example: vi
	TargetLang string `json:"target_lang"`
}

// LanguagesResponse represents supported languages response
// swagger:response languagesResponse
type LanguagesResponse struct {
	// List of supported language codes (ISO 639-1)
	// example: ["en", "vi", "fr", "es"]
	Languages []string `json:"languages"`
}

type Handler struct {
	authSvc        service.AuthService
	bookSvc        service.BookService
	translationSvc service.TranslationService
	redisClient    *redis.RedisClient
	HTTPConfig     *httpconfig.HTTPConfig
	AuthHandler    *AuthHandler
}

func NewHandler(
	authSvc service.AuthService,
	bookSvc service.BookService,
	translationSvc service.TranslationService,
	redisClient *redis.RedisClient,
	HTTPConfig *httpconfig.HTTPConfig,
) *Handler {
	h := &Handler{
		authSvc:        authSvc,
		bookSvc:        bookSvc,
		translationSvc: translationSvc,
		redisClient:    redisClient,
		HTTPConfig:     HTTPConfig,
	}
	h.AuthHandler = NewAuthHandler(authSvc)
	return h
}

// RegisterBookRoutes registers the book routes
// @Summary Register book routes
// @Description Register all book related routes
// @Tags books
// @Security BearerAuth
// @Router /api/books [get]
// @Router /api/books [post]
// @Router /api/books/{id} [get]
func (h *Handler) RegisterBookRoutes(router *gin.RouterGroup) {
	books := router.Group("/books")
	{
		books.GET("", h.ListBooks)
		books.POST("", h.CreateBook)
		books.GET("/:id", h.GetBook)
	}
}

// RegisterTranslationRoutes registers the translation routes
// @Summary Register translation routes
// @Description Register all translation related routes
// @Tags translations
// @Router /api/translations/translate [post]
// @Router /api/translations/languages [get]
func (h *Handler) RegisterTranslationRoutes(router *gin.RouterGroup) {
	translations := router.Group("/translations")
	{
		translations.POST("/translate", h.Translate)
		translations.GET("/languages", h.GetSupportedLanguages)
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

// ListBooks returns a list of books with pagination
// @Summary List all books
// @Description Get a paginated list of books
// @Tags books
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} BooksListResponse "Successfully retrieved books"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/books [get]
func (h *Handler) ListBooks(c *gin.Context) {
	// TODO: Implement list books with pagination
	c.JSON(http.StatusOK, BooksListResponse{
		Data:  []BookResponse{},
		Total: 0,
	})
}

// CreateBook creates a new book
// @Summary Create a new book
// @Description Create a new book with the provided information
// @Tags books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param book body BookInput true "Book data"
// @Success 201 {object} BookResponse "Successfully created book"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/books [post]
func (h *Handler) CreateBook(c *gin.Context) {
	var input BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// TODO: Implement create book
	c.JSON(http.StatusCreated, BookResponse{
		ID:            "generated-id",
		Title:         input.Title,
		Author:        input.Author,
		Description:   input.Description,
		PublishedYear: input.PublishedYear,
		CreatedAt:     "2023-01-01T00:00:00Z",
		UpdatedAt:     "2023-01-01T00:00:00Z",
	})
}

// GetBook gets a book by ID
// @Summary Get a book by ID
// @Description Get detailed information about a specific book
// @Tags books
// @Security BearerAuth
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} BookResponse "Successfully retrieved book"
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Book not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/books/{id} [get]
func (h *Handler) GetBook(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Book ID is required"})
		return
	}

	// TODO: Implement get book by ID
	c.JSON(http.StatusOK, BookResponse{
		ID:            id,
		Title:         "Sample Book",
		Author:        "Sample Author",
		Description:   "This is a sample book description",
		PublishedYear: 2023,
		CreatedAt:     "2023-01-01T00:00:00Z",
		UpdatedAt:     "2023-01-01T00:00:00Z",
	})
}

// Translate translates text from one language to another
// @Summary Translate text
// @Description Translate text from source language to target language
// @Tags translations
// @Accept json
// @Produce json
// @Param input body TranslateInput true "Translation input"
// @Success 200 {object} TranslateResponse "Successfully translated text"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Translation service error"
// @Router /api/translations/translate [post]
func (h *Handler) Translate(c *gin.Context) {
	var input TranslateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// TODO: Implement translation
	c.JSON(http.StatusOK, TranslateResponse{
		Text:       "Translated: " + input.Text,
		SourceLang: input.SourceLang,
		TargetLang: input.TargetLang,
	})
}

// GetSupportedLanguages returns a list of supported languages
// @Summary Get supported languages
// @Description Get a list of all supported languages for translation
// @Tags translations
// @Produce json
// @Success 200 {object} LanguagesResponse "Successfully retrieved supported languages"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/translations/languages [get]
func (h *Handler) GetSupportedLanguages(c *gin.Context) {
	// TODO: Get supported languages from service
	c.JSON(http.StatusOK, LanguagesResponse{
		Languages: []string{"en", "vi", "fr", "es"},
	})
}
