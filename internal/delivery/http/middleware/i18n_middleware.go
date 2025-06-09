package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get language from header, cookie, or URL parameter
		acceptLang := c.GetHeader("Accept-Language")
		if acceptLang == "" {
			acceptLang = "en" // default language
		}

		// Parse the language tag
		tag, _ := language.MatchStrings(language.NewMatcher([]language.Tag{
			language.English,
			language.Vietnamese,
			// add more supported languages
		}), acceptLang)

		c.Set("language", tag)
		c.Next()
	}
}
