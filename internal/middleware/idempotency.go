package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/partyscript/bledger/internal/cache"
	"github.com/partyscript/bledger/internal/common"
	"github.com/partyscript/bledger/internal/model"
	"github.com/partyscript/bledger/pkg"
)

// Idempotency is a middleware that checks for idempotency keys
func Idempotency(ic common.IdempotencyConfig, cm *cache.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request method is whitelisted
		if contains(ic.WhitelistedMethods, c.Request.Method) {
			c.Next()
			return
		}

		// Check if the request path is whitelisted
		if contains(ic.WhitelistedRoutes, c.Request.URL.Path) {
			c.Next()
			return
		}

		// Check if the request has an idempotency key
		reqID := c.GetHeader(ic.Header)
		if reqID == "" {
			c.AbortWithStatusJSON(
				common.WrapAPIError("idempotency key not found on request",
					common.BLedgerIdempotencyError,
					pkg.APIVersion,
				))
			return
		}

		var transaction model.CreateTransactionRequest

		ByteBody, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(ByteBody))

		err := json.Unmarshal(ByteBody, &transaction)
		if err != nil {
			c.AbortWithStatusJSON(
				common.WrapAPIError(err.Error(),
					common.BLedgerBadRequestError,
					pkg.APIVersion,
				),
			)
			return
		}

		// Check idempotency validity
		returnedID, err := common.GenIdempotencyKey(transaction.AccountID, transaction.Money, string(transaction.Direction))
		if err != nil {
			c.AbortWithStatusJSON(
				common.WrapAPIError(err.Error(),
					common.BLedgerBadRequestError,
					pkg.APIVersion,
				),
			)
			return
		}

		// Simple check for hmac-like validation
		if returnedID != reqID {
			c.AbortWithStatusJSON(
				common.WrapAPIError("idempotency key does not match request body",
					common.BLedgerIdempotencyError,
					pkg.APIVersion,
				))
			return
		}

		// Check if the idempotency key exists in the cache
		key, err := cm.Client.Get(c, reqID)
		if err != nil {
			c.AbortWithStatusJSON(
				common.WrapAPIError("failed to retrieve idempotency key from cache",
					common.BLedgerIdempotencyError,
					pkg.APIVersion,
				))
			return
		}

		// If the key exists, return conflict
		if key != nil {
			c.AbortWithStatusJSON(
				common.WrapAPIError("idempotency key already exists",
					common.BLedgerIdempotencyError,
					pkg.APIVersion,
				))
			return
		}

		// If the key does not exist, set it in the cache
		err = cm.Client.Set(c, reqID, common.IdempotencyLock{
			Lock: true,
			Body: transaction,
		}, time.Minute*60)

		// If the key could not be set, return an error
		if err != nil {
			c.AbortWithStatusJSON(
				common.WrapAPIError("idempotency key set failure",
					common.BLedgerIdempotencyError,
					pkg.APIVersion,
				))
			return
		}

		// Set the idempotency key in the context for use in the handler
		ic.SetIdempotencyKey(c, reqID)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
