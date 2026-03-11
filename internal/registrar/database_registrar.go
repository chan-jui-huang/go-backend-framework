package registrar

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/database"
)

type DatabaseRegistrar struct {
	config database.Config
}

func (dr *DatabaseRegistrar) Boot() {
	config.Registry.Register("database", &dr.config)
}

func (dr *DatabaseRegistrar) Register() {
	current := deps.CurrentConfig()
	current.DatabaseConfig = &dr.config
	deps.SetConfig(current)

	serviceState := deps.CurrentService()
	serviceState.DatabaseValue = database.New(dr.config)
	deps.SetService(serviceState)
}
