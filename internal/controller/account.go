package controller

import (
	"encoding/json"

	"github.com/partyscript/bledger/internal/cache"
	"github.com/partyscript/bledger/internal/config"
	"github.com/partyscript/bledger/internal/db"
	"github.com/partyscript/bledger/internal/model"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

// AccountController is the struct that the constructor implements
type AccountController struct {
	logger *zap.SugaredLogger
	cfg    *config.GlobalConfig
	cache  *cache.Manager
	db     *db.Manager
}

// NewAccountController initializes a AccountController instance
func NewAccountController(
	logger *zap.SugaredLogger,
	cfg *config.GlobalConfig,
	cache *cache.Manager,
	db *db.Manager,
) AccountController {
	return AccountController{
		logger: logger,
		cfg:    cfg,
		cache:  cache,
		db:     db,
	}
}

// GetAccount returns a transaction by id
func (ac *AccountController) GetAccount(id string) (*model.Account, error) {
	var acct model.Account

	find := ac.db.Gorm.First(&acct, &model.Account{ID: id})

	if find.Error != nil {
		return nil, find.Error
	}

	return &acct, nil
}

// CreateAccount creates a new account
func (ac *AccountController) CreateAccount(acctReq model.CreateAccountRequest) (*model.Account, error) {
	bal := model.AccountMoney{Currency: "USD"}

	amtToStore, err := json.Marshal(bal)
	if err != nil {
		return nil, err
	}

	acct := &model.Account{
		Name:        acctReq.Name,
		Description: acctReq.Description,
		Balance:     datatypes.JSON(amtToStore),
	}

	create := ac.db.Gorm.Create(acct)
	if create.Error != nil {
		return nil, create.Error
	}

	return acct, nil
}
