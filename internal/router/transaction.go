package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/partyscript/bledger/internal/model"

	"github.com/partyscript/bledger/internal/common"

	"github.com/partyscript/bledger/pkg"
)

// RegisterTransactionsRouter is a router register method that applies the routes to a router
func (m *Manager) RegisterTransactionsRouter(router *gin.RouterGroup) {
	router.GET("/:id", m.GetTransaction)
	router.POST("/", m.CreatePendingTransaction)
	router.PUT("/:id", m.ExecutePendingTransaction)
	router.POST("/immediate", m.CreateTransaction)
	router.DELETE("/:id", m.ReverseTransaction)
}

// ReverseTransaction is a router method that reverses a transaction
func (m *Manager) ReverseTransaction(c *gin.Context) {
	id, err := common.CheckID(c)
	if err != nil {
		c.JSON(
			common.WrapAPIError("id not found on request",
				common.BLedgerBadRequestError,
				pkg.APIVersion,
			),
		)
		return
	}

	err = m.Controller.Transactions.ReverseTransaction(id)
	if err != nil {
		c.JSON(
			common.WrapAPIError(err.Error(),
				common.BLedgerBadRequestError,
				pkg.APIVersion,
			),
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction reversed"})
}

// CreatePendingTransaction is a router method that creates a pending transaction
func (m *Manager) CreatePendingTransaction(c *gin.Context) {
	var txReq model.CreateTransactionRequest
	err := c.ShouldBindJSON(&txReq)
	if err != nil {
		c.JSON(
			common.WrapAPIError(err.Error(),
				common.BLedgerBadRequestError,
				pkg.APIVersion,
			),
		)
		return

	}

	transaction, err := m.Controller.Transactions.CreatePendingTransaction(txReq)
	if err != nil {
		c.JSON(
			common.WrapAPIError(err.Error(),
				common.BLedgerBadRequestError,
				pkg.APIVersion,
			),
		)
		return
	}

	c.JSON(http.StatusOK, &transaction)
}

// ExecutePendingTransaction is a router method that executes a pending transaction
func (m *Manager) ExecutePendingTransaction(c *gin.Context) {
	id, err := common.CheckID(c)
	if err != nil {
		c.JSON(
			common.WrapAPIError("id not found on request",
				common.BLedgerBadRequestError,
				pkg.APIVersion,
			),
		)
		return
	}

	transaction, err := m.Controller.Transactions.ExecutePendingTransaction(id)
	if err != nil {
		c.JSON(
			common.WrapAPIError(err.Error(),
				common.BLedgerBadRequestError,
				pkg.APIVersion,
			),
		)
		return
	}

	c.JSON(http.StatusOK, &transaction)
}

// GetTransaction is a router method that returns a single transaction
func (m *Manager) GetTransaction(c *gin.Context) {
	id, err := common.CheckID(c)
	if err != nil {
		c.JSON(
			common.WrapAPIError("id not found on request",
				common.BLedgerBadRequestError,
				pkg.APIVersion,
			),
		)
		return
	}

	transaction, err := m.Controller.Transactions.GetTransaction(id)
	if err != nil {
		c.JSON(
			common.WrapAPIError(err.Error(),
				common.BLedgerNotFoundError,
				pkg.APIVersion,
			),
		)
		return

	}

	c.JSON(http.StatusOK, &transaction)
}

// CreateTransaction is a router method that creates a new transaction
func (m *Manager) CreateTransaction(c *gin.Context) {
	var transaction model.CreateTransactionRequest
	err := c.ShouldBindJSON(&transaction)
	if err != nil {
		c.JSON(
			common.WrapAPIError(err.Error(),
				common.BLedgerBadRequestError,
				pkg.APIVersion,
			),
		)
		return

	}
	newTx, err := m.Controller.Transactions.CreateTransaction(transaction)
	if err != nil {
		c.JSON(
			common.WrapAPIError(err.Error(),
				common.BLedgerBadRequestError,
				pkg.APIVersion,
			),
		)
		return

	}
	c.JSON(http.StatusCreated, newTx)
}
