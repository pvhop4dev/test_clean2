// Package http provides HTTP handlers for the application.
//
// This package contains the HTTP handlers for the application's API endpoints.
// It includes documentation for all the available routes and their parameters,
// request/response formats, and error responses.
//
// The API follows RESTful principles and uses JSON for request and response bodies.
// All endpoints require authentication unless otherwise specified.
//
// BasePath: /api
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package http

import "clean-arch-go/internal/domain/entities"

// Generic error response
// swagger:response errorResponse

type errorResponseWrapper struct {
	// Error message
	// in: body
	Body struct {
		Error string `json:"error"`
	}
}

// Success response with authentication token
// swagger:response authResponse

type authResponseWrapper struct {
	// Authentication token and user info
	// in: body
	Body struct {
		Message string        `json:"message"`
		Token   string        `json:"token"`
		User    entities.User `json:"user"`
	}
}

// Success response with book data
// swagger:response bookResponse

type bookResponseWrapper struct {
	// Book data
	// in: body
	Body entities.Book
}

// Success response with list of books
// swagger:response booksResponse

type booksResponseWrapper struct {
	// List of books
	// in: body
	Body []entities.Book
}

// swagger:parameters registerUser

type registerUserParamsWrapper struct {
	// User registration data
	// in: body
	// required: true
	Body struct {
		// User's full name
		// required: true
		// example: John Doe
		Name string `json:"name"`

		// User's email address
		// required: true
		// example: john@example.com
		Email string `json:"email"`

		// User's password (must be at least 6 characters)
		// required: true
		// minLength: 6
		// example: password123
		Password string `json:"password"`
	}
}

// swagger:parameters loginUser

type loginUserParamsWrapper struct {
	// User login credentials
	// in: body
	// required: true
	Body struct {
		// User's email address
		// required: true
		// example: john@example.com
		Email string `json:"email"`

		// User's password
		// required: true
		// example: password123
		Password string `json:"password"`
	}
}

// swagger:parameters createBook updateBook

type bookParamsWrapper struct {
	// Book data
	// in: body
	// required: true
	Body struct {
		// Book title
		// required: true
		// example: The Go Programming Language
		Title string `json:"title"`

		// Book description
		// required: false
		// example: An introduction to programming in Go
		Description string `json:"description,omitempty"`

		// Book author
		// required: true
		// example: Alan A. A. Donovan, Brian W. Kernighan
		Author string `json:"author"`
	}
}
