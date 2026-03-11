package registrar

import (
	"fmt"
	"path"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/logger"
)

type LoggerRegistrar struct {
	consoleConfig logger.Config
	fileConfig    logger.Config
	accessConfig  logger.Config
}

func (lr *LoggerRegistrar) Boot() {
	config.Registry.RegisterMany(map[string]any{
		"logger.console": &logger.Config{},
		"logger.file":    &logger.Config{},
		"logger.access":  &logger.Config{},
	})

	lr.consoleConfig = config.Registry.Get("logger.console").(logger.Config)
	lr.fileConfig = config.Registry.Get("logger.file").(logger.Config)
	lr.accessConfig = config.Registry.Get("logger.access").(logger.Config)
}

func (lr *LoggerRegistrar) Register() {
	consoleLogger, err := logger.NewLogger(
		lr.consoleConfig,
		logger.ConsoleEncoder,
		logger.DefaultZapOptions...,
	)
	if err != nil {
		panic(err)
	}

	booterConfig := deps.BooterConfig()
	lr.fileConfig.LogPath = path.Join(booterConfig.RootDir, lr.fileConfig.LogPath)
	fileLogger, err := logger.NewLogger(
		lr.fileConfig,
		logger.JsonEncoder,
		logger.DefaultZapOptions...,
	)
	if err != nil {
		panic(err)
	}

	lr.accessConfig.LogPath = path.Join(booterConfig.RootDir, lr.accessConfig.LogPath)
	accessLogger, err := logger.NewLogger(
		lr.accessConfig,
		logger.JsonEncoder,
	)
	if err != nil {
		panic(err)
	}
	v := config.Registry.GetViper()
	settings := v.Sub("logger").AllSettings()
	defaultSetting := v.GetString("logger.default")
	defaultLogger := consoleLogger

	for setting := range settings {
		if defaultSetting == setting {
			switch fmt.Sprintf("logger.%s", defaultSetting) {
			case "logger.file":
				defaultLogger = fileLogger
			case "logger.access":
				defaultLogger = accessLogger
			default:
				defaultLogger = consoleLogger
			}
		}
	}

	current := deps.CurrentConfig()
	current.ConsoleLoggerConfig = &lr.consoleConfig
	current.FileLoggerConfig = &lr.fileConfig
	current.AccessLoggerConfig = &lr.accessConfig
	deps.SetConfig(current)

	serviceState := deps.CurrentService()
	serviceState.LoggerValue = defaultLogger
	serviceState.ConsoleLogger = consoleLogger
	serviceState.FileLogger = fileLogger
	serviceState.AccessLoggerValue = accessLogger
	deps.SetService(serviceState)
}
