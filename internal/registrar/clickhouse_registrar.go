package registrar

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/clickhouse"
)

type ClickhouseRegistrar struct {
	config clickhouse.Config
}

func (cr *ClickhouseRegistrar) Boot() {
	config.Registry.Register("clickhouse", &cr.config)
}

func (cr *ClickhouseRegistrar) Register() {
	conn, err := clickhouse.New(cr.config)
	if err != nil {
		panic(err)
	}

	current := deps.CurrentConfig()
	current.ClickhouseConfig = &cr.config
	deps.SetConfig(current)

	serviceState := deps.CurrentService()
	serviceState.ClickhouseValue = conn
	deps.SetService(serviceState)
}
