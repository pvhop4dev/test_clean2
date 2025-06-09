package service

import (
	"context"
	"sync"

	"clean-arch-go/internal/application/translation"
	"clean-arch-go/internal/errors"
)

// TranslationService defines the interface for translation operations
type TranslationService interface {
	Translate(ctx context.Context, text string, sourceLang string, targetLang string) (string, error)
	GetSupportedLanguages(ctx context.Context) ([]string, error)
}

type translationService struct {
	usecase translation.TranslationUsecase
	cache    map[string]string
	mu       sync.RWMutex
}

// NewTranslationService creates a new translation service
func NewTranslationService(usecase translation.TranslationUsecase) TranslationService {
	return &translationService{
		usecase: usecase,
		cache:   make(map[string]string),
	}
}

func (s *translationService) Translate(ctx context.Context, text string, sourceLang string, targetLang string) (string, error) {
	// Check cache first
	cacheKey := s.generateCacheKey(text, sourceLang, targetLang)
	s.mu.RLock()
	if translation, ok := s.cache[cacheKey]; ok {
		s.mu.RUnlock()
		return translation, nil
	}
	s.mu.RUnlock()

	// Get translation from usecase
	translation, err := s.usecase.Translate(ctx, text, targetLang)
	if err != nil {
		return "", errors.NewInternalServerError(err.Error())
	}

	// Cache the result
	s.mu.Lock()
	s.cache[cacheKey] = translation
	s.mu.Unlock()

	return translation, nil
}

func (s *translationService) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	// In a real implementation, this would fetch supported languages from a configuration or API
	return []string{"en", "vi", "fr", "es"}, nil
}

func (s *translationService) generateCacheKey(text, sourceLang, targetLang string) string {
	return text + "_" + sourceLang + "_" + targetLang
}
