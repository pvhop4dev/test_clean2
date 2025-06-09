package helper

import (
	"fmt"

	"clean-arch-go/internal/pkg/i18n"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// GetLocalizer returns the localizer instance with the correct language
func GetLocalizer(c *gin.Context) *i18n.Localizer {
	localizer := i18n.GetLocalizer()
	if lang, exists := c.Get("language"); exists {
		if langTag, ok := lang.(language.Tag); ok {
			localizer.SetLanguage(langTag)
		}
	}
	return localizer
}

// Translate translates a message with the given ID and template data
func Translate(c *gin.Context, messageID string, templateData map[string]interface{}) (string, error) {
	localizer := GetLocalizer(c)
	lang, _ := c.Get("language")
	if langTag, ok := lang.(language.Tag); ok {
		return localizer.Translate(langTag, messageID, templateData)
	}
	return "", fmt.Errorf("language not found")
}

// MustTranslate translates a message or panics if there's an error
func MustTranslate(c *gin.Context, messageID string, templateData map[string]interface{}) string {
	msg, err := Translate(c, messageID, templateData)
	if err != nil {
		panic(fmt.Sprintf("failed to translate message %s: %v", messageID, err))
	}
	return msg
}
