package repository

import (
	"context"

	"clean-arch-go/internal/domain/entities"
	"clean-arch-go/internal/errors"

	"clean-arch-go/internal/pkg/database"

	"gorm.io/gorm"
)

type TranslationRepository interface {
	BaseRepository[entities.Translation]
	Translate(ctx context.Context, text string, targetLang string) (string, error)
	GetTranslation(ctx context.Context, text string, targetLang string) (string, error)
	SaveTranslation(ctx context.Context, text string, targetLang string, translation string) error
}

type translationRepository struct {
	*baseRepository[entities.Translation]
}

func NewTranslationRepository(db *database.Database) TranslationRepository {
	return &translationRepository{
		baseRepository: NewBaseRepository[entities.Translation](db.DB).(*baseRepository[entities.Translation]),
	}
}

func (r *translationRepository) Translate(ctx context.Context, text string, targetLang string) (string, error) {
	// In a real implementation, this would call an external translation API
	// For now, we'll just return the text as-is
	return text, nil
}

func (r *translationRepository) GetTranslation(ctx context.Context, text string, targetLang string) (string, error) {
	var translation entities.Translation
	if err := r.baseRepository.db.WithContext(ctx).
		Where("source_text = ? AND target_lang = ?", text, targetLang).
		First(&translation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.NewNotFoundError("translation not found")
		}
		return "", errors.NewInternalServerError(err.Error())
	}
	return translation.TranslatedText, nil
}

func (r *translationRepository) SaveTranslation(ctx context.Context, text string, targetLang string, translation string) error {
	newTranslation := entities.Translation{
		SourceText:   text,
		TargetLang:   targetLang,
		TranslatedText: translation,
	}
	if err := r.baseRepository.db.WithContext(ctx).Create(&newTranslation).Error; err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}
