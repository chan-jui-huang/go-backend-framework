package registrar_test

import (
	"testing"

	_ "github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/registrar"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/stretchr/testify/suite"
)

type RegisterExecutorTestSuite struct {
	suite.Suite
	booterConfig booter.Config
	viper        viper.Viper
}

func (suite *RegisterExecutorTestSuite) SetupSuite() {
	suite.booterConfig = config.Registry.Get("booter").(booter.Config)
	suite.viper = config.Registry.GetViper()
}

func (suite *RegisterExecutorTestSuite) TestRegisterExecutor() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	service.Registry = service.NewRegistry()

	registrar.RegisterExecutor.AfterExecute()
	registrar.RegisterExecutor.Execute()
	registrar.RegisterExecutor.BeforeExecute()

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

	suite.NotNil(service.Registry.Get("authentication.authenticator"))
	suite.NotNil(service.Registry.Get("database"))
	suite.NotNil(service.Registry.Get("casbinEnforcer"))
	suite.NotNil(service.Registry.Get("logger"))
	suite.NotNil(service.Registry.Get("logger.console"))
	suite.NotNil(service.Registry.Get("logger.file"))
	suite.NotNil(service.Registry.Get("logger.access"))
	suite.NotNil(service.Registry.Get("redis"))
	suite.NotNil(service.Registry.Get("formDecoder"))
	suite.NotNil(service.Registry.Get("modifier"))
	suite.NotNil(service.Registry.Get("mapstructureDecoder"))
	suite.NotNil(service.Registry.Get("clickhouse"))
}

func (suite *RegisterExecutorTestSuite) TestSimpleRegisterExecutor() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	service.Registry = service.NewRegistry()

	registrar.SimpleRegisterExecutor.AfterExecute()
	registrar.SimpleRegisterExecutor.Execute()
	registrar.SimpleRegisterExecutor.BeforeExecute()

	suite.NotNil(config.Registry.Get("httpServer"))
	suite.NotNil(config.Registry.Get("middleware.csrf"))
	suite.NotNil(config.Registry.Get("middleware.rateLimit"))
	suite.NotNil(config.Registry.Get("authentication.authenticator"))
	suite.NotNil(config.Registry.Get("logger.console"))
	suite.NotNil(config.Registry.Get("logger.file"))
	suite.NotNil(config.Registry.Get("logger.access"))

	suite.NotNil(service.Registry.Get("authentication.authenticator"))
	suite.NotNil(service.Registry.Get("logger.console"))
	suite.NotNil(service.Registry.Get("logger.file"))
	suite.NotNil(service.Registry.Get("logger.access"))
	suite.NotNil(service.Registry.Get("formDecoder"))
	suite.NotNil(service.Registry.Get("modifier"))
	suite.NotNil(service.Registry.Get("mapstructureDecoder"))
}

func (suite *RegisterExecutorTestSuite) TestValidatorRegistrarTagNameFunc() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	service.Registry = service.NewRegistry()

	registrar.RegisterExecutor.AfterExecute()
	registrar.RegisterExecutor.Execute()
	registrar.RegisterExecutor.BeforeExecute()

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
	service.Registry = service.NewRegistry()

	registrar.RegisterExecutor.AfterExecute()
	registrar.RegisterExecutor.Execute()
	registrar.RegisterExecutor.BeforeExecute()
}

func TestRegistrarTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterExecutorTestSuite))
}
