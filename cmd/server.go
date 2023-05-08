package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// This will automatically load/inject environment variables from a .env file
	_ "github.com/joho/godotenv/autoload"
	"github.com/partyscript/bledger/internal/cache"
	"github.com/partyscript/bledger/internal/cache/redis"
	"github.com/partyscript/bledger/internal/common"
	"github.com/partyscript/bledger/internal/config"
	"github.com/partyscript/bledger/internal/controller"
	"github.com/partyscript/bledger/internal/db"
	"github.com/partyscript/bledger/internal/router"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Set global config
	cfg, err := config.NewGlobalConfig()
	if err != nil {
		log.Fatal(err)
	}

	// New Prod Logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	// Flush stdout
	defer zapLogger.Sync()

	// SugaredLogger for pretty stdout
	sugaredLogger := zapLogger.Sugar()

	// Create new instance of redis cache
	redis, err := redis.NewRedisCache(cfg.Cache.URI, cfg.Cache.Password, time.Minute*15)
	if err != nil {
		log.Fatal(err)
	}

	// Create new instance of cache manager
	cacher := cache.NewCacheManager(redis)

	// Set ORM w/ upstream sql database
	gorm, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a new instance of a db manager
	dbManager, err := db.NewDBManager(gorm)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new instance of a controller manager
	cm := controller.NewControllerManager(
		sugaredLogger,
		cfg,
		cacher,
		dbManager,
	)

	ic := common.NewIdempotencyConfig(
		[]string{"GET", "HEAD", "OPTIONS", "TRACE"},
		[]string{},
		"idemKey",
		http.StatusForbidden,
	)

	// Create a new instance of a router manager
	rm := router.NewRouterManager(
		cm,
		ic,
		cacher,
	)

	// Initialize the router
	err = rm.InitRouter()
	if err != nil {
		log.Fatal(err)
	}

	// Start the server
	err = rm.Router.Run(fmt.Sprintf("0.0.0.0%v", common.FormatPort(*cfg.Environment)))
	if err != nil {
		log.Fatal(err)
	}
}
