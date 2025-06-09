package translation

import (
	"clean-arch-go/internal/errors"
	"context"

	"gorm.io/gorm"
)

// Repository defines the interface for translation data access
//go:generate mockery --name=Repository --output=../mocks/translation

type Repository interface {
	Translate(ctx context.Context, text string, targetLang string) (string, error)
	GetTranslation(ctx context.Context, text string, targetLang string) (string, error)
	SaveTranslation(ctx context.Context, text string, targetLang string, translation string) error
}

// ErrTranslationNotFound is returned when translation is not found
var ErrTranslationNotFound = errors.NewAppError("TRANSLATION_NOT_FOUND", "Translation not found", nil)

// NewRepository creates a new translation repository
func NewRepository(db *gorm.DB) Repository {
	return &translationRepository{
		db: db,
	}
}

type translationRepository struct {
	db *gorm.DB
}

func (r *translationRepository) Translate(ctx context.Context, text string, targetLang string) (string, error) {
	// Implementation here
	return "", ErrTranslationNotFound
}

func (r *translationRepository) GetTranslation(ctx context.Context, text string, targetLang string) (string, error) {
	// Implementation here
	return "", ErrTranslationNotFound
}

func (r *translationRepository) SaveTranslation(ctx context.Context, text string, targetLang string, translation string) error {
	// Implementation here
	return nil
}
