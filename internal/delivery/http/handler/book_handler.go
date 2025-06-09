package handler

import (
	"net/http"
	"strconv"

	"clean-arch-go/internal/domain/entities"
	"clean-arch-go/internal/domain/service"
	"clean-arch-go/internal/pkg/i18n"

	"golang.org/x/text/language"

	"github.com/gin-gonic/gin"
)

type CreateBookRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Author      string `json:"author" binding:"required"`
}

type UpdateBookRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

type BookHandler struct {
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) GetBook(c *gin.Context) {
	// Get book ID from URL
	bookID := c.Param("id")

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.unauthorized", nil)})
		return
	}

	// Check if the user owns the book
	if err := h.bookService.CheckBookOwnership(c.Request.Context(), bookID, userID.(string)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.forbidden", nil)})
		return
	}

	book, err := h.bookService.GetBookByID(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.not_found", nil)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": i18n.GetLocalizer().MustTranslate(language.English, "book.get_success", nil),
		"data":    book,
	})
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.unauthorized", nil)})
		return
	}

	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := &entities.Book{
		Title:       req.Title,
		Description: req.Description,
		Author:      req.Author,
		UserID:      userID.(string),
	}

	if err := h.bookService.CreateBook(c.Request.Context(), book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.internal_server_error", nil)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": i18n.GetLocalizer().MustTranslate(language.English, "book.created", nil),
		"data":    book,
	})
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	// Get book ID from URL
	bookID := c.Param("id")

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.unauthorized", nil)})
		return
	}

	// Check if the user owns the book
	if err := h.bookService.CheckBookOwnership(c.Request.Context(), bookID, userID.(string)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.forbidden", nil)})
		return
	}

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing book
	book, err := h.bookService.GetBookByID(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.not_found", nil)})
		return
	}

	// Update book fields
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.Author != "" {
		book.Author = req.Author
	}

	if err := h.bookService.UpdateBook(c.Request.Context(), bookID, book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.internal_server_error", nil)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": i18n.GetLocalizer().MustTranslate(language.English, "book.updated", nil),
		"data":    book,
	})
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	// Get book ID from URL
	bookID := c.Param("id")

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.unauthorized", nil)})
		return
	}

	// Check if the user owns the book
	if err := h.bookService.CheckBookOwnership(c.Request.Context(), bookID, userID.(string)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.forbidden", nil)})
		return
	}

	if err := h.bookService.DeleteBook(c.Request.Context(), bookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.internal_server_error", nil)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": i18n.GetLocalizer().MustTranslate(language.English, "book.deleted", nil),
	})
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.unauthorized", nil)})
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	books, err := h.bookService.ListBooksByUserID(c.Request.Context(), userID.(string), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.GetLocalizer().MustTranslate(language.English, "error.internal_server_error", nil)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": i18n.GetLocalizer().MustTranslate(language.English, "book.list_success", nil),
		"data":    books,
	})
}
