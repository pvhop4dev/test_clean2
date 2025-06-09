// internal/pkg/i18n/i18n.go
package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*
var localesFS embed.FS

var (
	instance *Localizer
	once     sync.Once
)

// Localizer is a wrapper around i18n.Localizer with additional functionality
type Localizer struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
	mu        sync.RWMutex
}

// GetLocalizer returns the singleton instance of Localizer
func GetLocalizer() *Localizer {
	once.Do(func() {
		defaultLang := language.English
		bundle := i18n.NewBundle(defaultLang)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		// Load all translation files
		entries, err := localesFS.ReadDir("locales")
		if err != nil {
			panic(fmt.Sprintf("failed to read locales directory: %v", err))
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				_, err := bundle.LoadMessageFileFS(localesFS, "locales/"+entry.Name())
				if err != nil {
					panic(fmt.Sprintf("failed to load message file %s: %v", entry.Name(), err))
				}
			}
		}

		instance = &Localizer{
			bundle:    bundle,
			localizer: i18n.NewLocalizer(bundle, defaultLang.String()),
		}
	})
	return instance
}

// SetLanguage changes the language of the localizer
func (l *Localizer) SetLanguage(lang language.Tag) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.localizer = i18n.NewLocalizer(l.bundle, lang.String())
}

// Translate translates a message with the given ID and template data
func (l *Localizer) Translate(lang language.Tag, messageID string, templateData map[string]interface{}) (string, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
}

// MustTranslate translates a message or panics if there's an error
func (l *Localizer) MustTranslate(lang language.Tag, messageID string, templateData map[string]interface{}) string {
	msg, err := l.Translate(lang, messageID, templateData)
	if err != nil {
		panic(fmt.Sprintf("failed to translate message %s: %v", messageID, err))
	}
	return msg
}