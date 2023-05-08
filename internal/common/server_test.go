package common

import (
	"testing"

	"github.com/partyscript/bledger/internal/config"
	"github.com/partyscript/bledger/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestFormatPort(t *testing.T) {
	cfg := config.EnvironmentConfig{
		Env:  model.ApplicationEnvironmentDev,
		Port: "8080",
	}
	port := FormatPort(cfg)
	assert.Equal(t, port, ":8080")
}
