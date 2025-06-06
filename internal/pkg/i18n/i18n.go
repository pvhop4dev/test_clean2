package i18n

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	// "github.com/nicksnyder/go-i18n/v2/i18n/json"
	"golang.org/x/text/language"
)

//go:embed locales/*
var localesFS embed.FS

type Localizer struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
}

func NewLocalizer(defaultLang language.Tag) (*Localizer, error) {
	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("json", func(data []byte, v interface{}) error {
		return json.Unmarshal(data, v)
	})

	// Load all translation files
	entries, err := localesFS.ReadDir("locales")
	if err != nil {
		return nil, fmt.Errorf("failed to read locales directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, err := localesFS.ReadFile(fmt.Sprintf("locales/%s", entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read locale file %s: %w", entry.Name(), err)
		}

		if _, err := bundle.ParseMessageFileBytes(data, entry.Name()); err != nil {
			return nil, fmt.Errorf("failed to parse locale file %s: %w", entry.Name(), err)
		}
	}

	return &Localizer{
		bundle:    bundle,
		localizer: i18n.NewLocalizer(bundle, defaultLang.String()),
	}, nil
}

func (l *Localizer) SetLanguage(lang string) {
	l.localizer = i18n.NewLocalizer(l.bundle, lang)
}

func (l *Localizer) Translate(ctx context.Context, messageID string, templateData map[string]interface{}) (string, error) {
	acceptLang := ctx.Value("Accept-Language")
	if acceptLang != nil {
		if lang, ok := acceptLang.(string); ok && lang != "" {
			l.SetLanguage(lang)
		}
	}

	return l.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
}

// Helper function to create a new localizer with default settings
func NewDefaultLocalizer() (*Localizer, error) {
	return NewLocalizer(language.English)
}
