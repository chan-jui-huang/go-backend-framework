package registrar

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	booterconfig "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/clickhouse"
)

func NewClickhouseConfig(loader *booterconfig.Loader) *clickhouse.Config {
	config := &clickhouse.Config{}
	loader.Unmarshal("clickhouse", config)

	return config
}

func NewClickhouse(config *clickhouse.Config) (driver.Conn, error) {
	return clickhouse.New(*config)
}
