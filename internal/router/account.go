package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/partyscript/bledger/internal/common"
	"github.com/partyscript/bledger/internal/model"
	"github.com/partyscript/bledger/pkg"
)

// RegisterAccountsRouter is a router register method that applies the routes to a router
func (m *Manager) RegisterAccountsRouter(router *gin.RouterGroup) {
	router.GET("/:id", m.GetAccount)
	router.POST("/", m.CreateAccount)
}

// GetAccount is a router method that returns a single account
func (m *Manager) GetAccount(c *gin.Context) {
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

	acct, err := m.Controller.Accounts.GetAccount(id)
	if err != nil {
		c.JSON(
			common.WrapAPIError(err.Error(),
				common.BLedgerNotFoundError,
				pkg.APIVersion,
			),
		)
		return

	}

	c.JSON(http.StatusOK, &acct)
}

// CreateAccount is a router method that creates a new account
func (m *Manager) CreateAccount(c *gin.Context) {
	var acctReq model.CreateAccountRequest
	err := c.ShouldBindJSON(&acctReq)
	if err != nil {
		c.JSON(
			common.WrapAPIError(err.Error(),
				common.BLedgerBadRequestError,
				pkg.APIVersion,
			),
		)
		return
	}

	acct, err := m.Controller.Accounts.CreateAccount(acctReq)
	if err != nil {
		c.JSON(
			common.WrapAPIError(err.Error(),
				common.BLedgerBadRequestError,
				pkg.APIVersion,
			),
		)
		return

	}

	c.JSON(http.StatusOK, &acct)
}
