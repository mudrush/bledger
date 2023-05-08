package middleware

import (
	"github.com/partyscript/bledger/internal/common"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// JSONLogger logs a gin HTTP request in JSON format, with some additional custom key/values
func JSONLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process Request
		c.Next()

		entry := log.WithFields(log.Fields{
			"method":             c.Request.Method,
			"path":               c.Request.RequestURI,
			"status":             c.Writer.Status(),
			"referrer":           c.Request.Referer(),
			"agent":              c.Request.UserAgent(),
			"idempotency_header": c.GetHeader(common.IdempotencyHeader),
			"ip":                 c.ClientIP(),
		})

		if c.Writer.Status() >= 500 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("")
		}
	}
}
