package common

import (
	b64 "encoding/base64"

	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/partyscript/bledger/internal/model"
)

// GenIdempotencyKey generates a unique idempotency key
func GenIdempotencyKey(
	accountID string,
	money model.TransactionMoney,
	direction string,
) (string, error) {
	if accountID == "" {
		return "", errors.New("account id cannot be empty")
	}

	if money.Amount == 0 {
		return "", errors.New("money amount cannot be empty")
	}

	if money.Currency == "" {
		return "", errors.New("money currency cannot be empty")
	}

	if direction == "" {
		return "", errors.New("direction cannot be empty")
	}

	return b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v-%v-%v-%v", accountID, money.Amount, money.Currency, direction))), nil
}

// CheckID checks if gin context has id
func CheckID(c *gin.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", errors.New("id cannot be empty")
	}
	return id, nil
}

// IsCompletableState checks if transaction is completable
func IsCompletableState(tx model.Transaction) bool {
	return tx.State == model.TransactionStatePending
}
