package handler

import (
	"clean-arch-go/internal/domain/entities"
	"clean-arch-go/internal/domain/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

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

func (h *BookHandler) CreateBook(c *gin.Context) {
	// Lấy user ID từ context (đã được set bởi auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
		UserID:      userID.(uint),
	}

	if err := h.bookService.CreateBook(c.Request.Context(), book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	// Lấy book ID từ URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	// Lấy user ID từ context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Lấy thông tin sách hiện tại
	book, err := h.bookService.GetBookByID(c.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// Cập nhật thông tin
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.Author != "" {
		book.Author = req.Author
	}

	if err := h.bookService.UpdateBook(c.Request.Context(), book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	// Lấy book ID từ URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	// Lấy user ID từ context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	book, err := h.bookService.GetBookByID(c.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	// Lấy book ID từ URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	// Lấy user ID từ context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.bookService.DeleteBook(c.Request.Context(), uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	// Lấy user ID từ context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	books, err := h.bookService.ListBooksByUserID(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}
