package common

import (
	"fmt"

	"github.com/partyscript/bledger/internal/config"
)

// FormatPort returns a formatted port for httpServer consumption
func FormatPort(cfg config.EnvironmentConfig) string {
	return fmt.Sprintf(":%v", cfg.Port)
}
