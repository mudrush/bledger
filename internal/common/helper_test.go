package common

import (
	"testing"

	"github.com/partyscript/bledger/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGenIdempotencyKey(t *testing.T) {
	// Test cases
	testCases := []struct {
		accountID string
		money     model.TransactionMoney
		direction string
		expected  string
	}{
		{
			accountID: "123",
			money:     model.TransactionMoney{Amount: 100, Currency: "USD"},
			direction: "debit",
			expected:  "MTIzLTEwMC1VU0QtZGViaXQ=",
		},
		{
			accountID: "123",
			money:     model.TransactionMoney{Amount: 100, Currency: "USD"},
			direction: "credit",
			expected:  "MTIzLTEwMC1VU0QtY3JlZGl0",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		actual, err := GenIdempotencyKey(tc.accountID, tc.money, tc.direction)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}
