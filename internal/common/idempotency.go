package common

import (
	"github.com/gin-gonic/gin"
)

// IdempotencyConfig is a struct that contains the configuration for the idempotency middleware
type IdempotencyConfig struct {
	WhitelistedMethods []string
	WhitelistedRoutes  []string
	ContextKey         string
	StatusCode         int
	Header             string
}

// IdempotencyLock is a struct that contains the lock and the body of the request
type IdempotencyLock struct {
	Lock bool
	Body interface{}
}

// NewIdempotencyConfig returns a new instance of IdempotencyConfig
func NewIdempotencyConfig(whitelistedMethods []string, whitelistedRoutes []string, contextKey string, statusCode int) IdempotencyConfig {
	return IdempotencyConfig{
		WhitelistedMethods: whitelistedMethods,
		WhitelistedRoutes:  whitelistedRoutes,
		ContextKey:         contextKey,
		StatusCode:         statusCode,
		Header:             IdempotencyHeader,
	}
}

// SetIdempotencyKey sets the idempotency key in the context
func (ic *IdempotencyConfig) SetIdempotencyKey(c *gin.Context, key string) {
	c.Set(ic.ContextKey, key)
}

// IdempotencyKey returns the idempotency key from the context
func (ic *IdempotencyConfig) IdempotencyKey(c *gin.Context) string {
	return c.GetString(ic.ContextKey)
}
