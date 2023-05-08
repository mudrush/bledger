package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/partyscript/bledger/internal/cache"
	"github.com/partyscript/bledger/internal/common"
	"github.com/partyscript/bledger/internal/controller"
	"github.com/partyscript/bledger/internal/middleware"
	"github.com/partyscript/bledger/internal/model"
)

// Manager is a struct that holds the controller manager and a router
type Manager struct {
	Controller controller.Manager
	IdemConfig common.IdempotencyConfig
	Router     *gin.Engine
	Cache      *cache.Manager
}

const (
	versionRouterGroup     = "/v1"
	transactionRouterGroup = "/transactions"
	accountRouterGroup     = "/accounts"
)

// RegisterRouters is a router method to add all nested routers to the engine
func (m *Manager) RegisterRouters(router *gin.Engine) {
	// V1 API Router Group
	v1 := router.Group(versionRouterGroup)

	// Transaction Router Group
	m.RegisterTransactionsRouter(v1.Group(transactionRouterGroup))

	// Account Router Group
	m.RegisterAccountsRouter(v1.Group(accountRouterGroup))
}

// NewRouterManager is a constructor that returns a new instance of RouterManager
func NewRouterManager(
	cm controller.Manager,
	ic common.IdempotencyConfig,
	cache *cache.Manager,
) *Manager {

	gin.SetMode(gin.ReleaseMode)

	if (cm.Cfg.Environment.Env == model.ApplicationEnvironmentDev) ||
		(cm.Cfg.Environment.Env == model.ApplicationEnvironmentLocal) {
		gin.SetMode(gin.DebugMode)
	}

	return &Manager{
		Controller: cm,
		Router:     gin.New(),
		IdemConfig: ic,
		Cache:      cache,
	}

}

// InitRouter is a method that initializes the router
func (m *Manager) InitRouter() error {
	m.Router.Use(middleware.JSONLogger())

	// Recover from panics
	m.Router.Use(gin.Recovery())

	// Cors
	m.Router.Use(cors.Default())

	// Idempotency
	// m.Router.Use(middleware.Idempotency(m.IdemConfig, m.Cache))

	// No proxies
	err := m.Router.SetTrustedProxies(nil)
	if err != nil {
		return err
	}

	// Register all routers
	m.RegisterRouters(m.Router)

	// Health Check
	m.Router.GET("/health_check", func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	})

	// Default route
	m.Router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"running": true,
		})
	})

	return nil
}
