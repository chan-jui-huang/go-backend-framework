package registrar_test

import (
	"testing"

	_ "github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/registrar"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/config"
	"github.com/stretchr/testify/suite"
)

type RegisterExecutorTestSuite struct {
	suite.Suite
	booterConfig booter.Config
	viper        viper.Viper
}

func (suite *RegisterExecutorTestSuite) SetupSuite() {
	suite.booterConfig = deps.BooterConfig()
	suite.viper = config.Registry.GetViper()
}

func (suite *RegisterExecutorTestSuite) TestRegisterExecutor() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	deps.SetConfig(deps.ConfigState{BooterConfig: &suite.booterConfig})
	deps.SetService(deps.ServiceState{})

	registrar.RegisterExecutor.BeforeExecute()
	registrar.RegisterExecutor.Execute()
	registrar.RegisterExecutor.AfterExecute()

	suite.NotNil(config.Registry.Get("httpServer"))
	suite.NotNil(config.Registry.Get("middleware.csrf"))
	suite.NotNil(config.Registry.Get("middleware.rateLimit"))
	suite.NotNil(config.Registry.Get("authentication.authenticator"))
	suite.NotNil(config.Registry.Get("database"))
	suite.NotNil(config.Registry.Get("logger.console"))
	suite.NotNil(config.Registry.Get("logger.file"))
	suite.NotNil(config.Registry.Get("logger.access"))
	suite.NotNil(config.Registry.Get("redis"))
	suite.NotNil(config.Registry.Get("clickhouse"))

	currentConfig := deps.CurrentConfig()
	currentService := deps.CurrentService()
	suite.NotNil(currentConfig.AuthenticationConfig)
	suite.NotNil(currentService.AuthenticatorValue)
	suite.NotNil(currentConfig.DatabaseConfig)
	suite.NotNil(currentService.DatabaseValue)
	suite.NotNil(currentService.CasbinEnforcerValue)
	suite.NotNil(currentService.LoggerValue)
	suite.NotNil(currentService.ConsoleLogger)
	suite.NotNil(currentService.FileLogger)
	suite.NotNil(currentService.AccessLoggerValue)
	suite.NotNil(currentConfig.RedisConfig)
	suite.NotNil(currentService.RedisValue)
	suite.NotNil(currentService.FormDecoder)
	suite.NotNil(currentService.Modifier)
	suite.NotNil(currentService.MapstructureDecoder)
	suite.NotNil(currentConfig.ClickhouseConfig)
	suite.NotNil(currentService.ClickhouseValue)
}

func (suite *RegisterExecutorTestSuite) TestSimpleRegisterExecutor() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	deps.SetConfig(deps.ConfigState{BooterConfig: &suite.booterConfig})
	deps.SetService(deps.ServiceState{})

	registrar.SimpleRegisterExecutor.BeforeExecute()
	registrar.SimpleRegisterExecutor.Execute()
	registrar.SimpleRegisterExecutor.AfterExecute()

	suite.NotNil(config.Registry.Get("httpServer"))
	suite.NotNil(config.Registry.Get("middleware.csrf"))
	suite.NotNil(config.Registry.Get("middleware.rateLimit"))
	suite.NotNil(config.Registry.Get("authentication.authenticator"))
	suite.NotNil(config.Registry.Get("logger.console"))
	suite.NotNil(config.Registry.Get("logger.file"))
	suite.NotNil(config.Registry.Get("logger.access"))

	currentConfig := deps.CurrentConfig()
	currentService := deps.CurrentService()
	suite.NotNil(currentConfig.AuthenticationConfig)
	suite.NotNil(currentService.AuthenticatorValue)
	suite.NotNil(currentService.ConsoleLogger)
	suite.NotNil(currentService.FileLogger)
	suite.NotNil(currentService.AccessLoggerValue)
	suite.NotNil(currentService.FormDecoder)
	suite.NotNil(currentService.Modifier)
	suite.NotNil(currentService.MapstructureDecoder)
}

func (suite *RegisterExecutorTestSuite) TestValidatorRegistrarTagNameFunc() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	deps.SetConfig(deps.ConfigState{BooterConfig: &suite.booterConfig})
	deps.SetService(deps.ServiceState{})

	registrar.RegisterExecutor.BeforeExecute()
	registrar.RegisterExecutor.Execute()
	registrar.RegisterExecutor.AfterExecute()

	type dummy struct {
		Email string `json:"email" binding:"required,email"`
		Page  int    `form:"page" binding:"required,gt=0"`
	}

	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		suite.Fail("validator engine is not *validator.Validate")
		return
	}

	err := engine.Struct(&dummy{})
	suite.Error(err)

	ves, ok := err.(validator.ValidationErrors)
	if !ok {
		suite.Fail("error is not ValidationErrors")
		return
	}

	fields := map[string]string{}
	for _, fe := range ves {
		fields[fe.Field()] = fe.Tag()
	}

	suite.Equal(map[string]string{
		"email": "required",
		"page":  "required",
	}, fields)
}

func (suite *RegisterExecutorTestSuite) TearDownSuite() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	deps.SetConfig(deps.ConfigState{BooterConfig: &suite.booterConfig})
	deps.SetService(deps.ServiceState{})

	registrar.RegisterExecutor.BeforeExecute()
	registrar.RegisterExecutor.Execute()
	registrar.RegisterExecutor.AfterExecute()
}

func TestRegistrarTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterExecutorTestSuite))
}
