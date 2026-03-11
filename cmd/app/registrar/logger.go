package registrar

import (
	"path"

	booter "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	booterconfig "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type LoggerConfigs struct {
	fx.Out

	Default string         `mapstructure:"default"`
	Console *logger.Config `name:"logger.console" mapstructure:"console"`
	File    *logger.Config `name:"logger.file" mapstructure:"file"`
	Access  *logger.Config `name:"logger.access" mapstructure:"access"`
}

type LoggerServices struct {
	fx.Out

	Logger  *zap.Logger
	Console *zap.Logger `name:"logger.console"`
	File    *zap.Logger `name:"logger.file"`
	Access  *zap.Logger `name:"logger.access"`
}

type NewLoggersParams struct {
	fx.In

	BooterConfig  *booter.Config
	LoggerConfig  LoggerConfigs
	ConsoleConfig *logger.Config `name:"logger.console"`
	FileConfig    *logger.Config `name:"logger.file"`
	AccessConfig  *logger.Config `name:"logger.access"`
}

func NewLoggerConfigs(loader *booterconfig.Loader) LoggerConfigs {
	config := LoggerConfigs{
		Console: &logger.Config{},
		File:    &logger.Config{},
		Access:  &logger.Config{},
	}
	loader.Unmarshal("logger", &config)

	return config
}

func NewLoggers(params NewLoggersParams) (LoggerServices, error) {
	consoleLogger, err := logger.NewLogger(
		*params.ConsoleConfig,
		logger.ConsoleEncoder,
		logger.DefaultZapOptions...,
	)
	if err != nil {
		return LoggerServices{}, err
	}

	fileConfig := *params.FileConfig
	fileConfig.LogPath = path.Join(params.BooterConfig.RootDir, fileConfig.LogPath)
	fileLogger, err := logger.NewLogger(
		fileConfig,
		logger.JsonEncoder,
		logger.DefaultZapOptions...,
	)
	if err != nil {
		return LoggerServices{}, err
	}

	accessConfig := *params.AccessConfig
	accessConfig.LogPath = path.Join(params.BooterConfig.RootDir, accessConfig.LogPath)
	accessLogger, err := logger.NewLogger(
		accessConfig,
		logger.JsonEncoder,
	)
	if err != nil {
		return LoggerServices{}, err
	}

	defaultLogger := consoleLogger
	switch params.LoggerConfig.Default {
	case "file":
		defaultLogger = fileLogger
	case "access":
		defaultLogger = accessLogger
	}

	return LoggerServices{
		Logger:  defaultLogger,
		Console: consoleLogger,
		File:    fileLogger,
		Access:  accessLogger,
	}, nil
}
