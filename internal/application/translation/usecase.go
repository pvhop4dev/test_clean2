package translation

import (
	"context"

	"clean-arch-go/internal/domain/translation"
)

// TranslationUsecase defines the interface for translation use cases
//go:generate mockery --name=TranslationUsecase --output=../mocks/translation

type TranslationUsecase interface {
	Translate(ctx context.Context, text string, targetLang string) (string, error)
}

type translationUsecase struct {
	repo translation.Repository
}

func NewTranslationUsecase(repo translation.Repository) TranslationUsecase {
	return &translationUsecase{
		repo: repo,
	}
}

func (u *translationUsecase) Translate(ctx context.Context, text string, targetLang string) (string, error) {
	// Business logic here
	return u.repo.Translate(ctx, text, targetLang)
}
