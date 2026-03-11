package registrar_test

import (
	"testing"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
)

type RegistrarTestSuite struct {
	suite.Suite
}

func (suite *RegistrarTestSuite) SetupSuite() {
	test.Setup(suite.T())
}

func (suite *RegistrarTestSuite) TestDependenciesRegistered() {
	currentConfig := deps.CurrentConfig()
	currentService := deps.CurrentService()
	suite.NotNil(currentConfig.BooterConfig)
	suite.NotNil(currentConfig.CsrfConfigValue)
	suite.NotNil(currentConfig.RateLimitConfigValue)
	suite.NotNil(currentConfig.AuthenticationConfig)
	suite.NotNil(currentConfig.DatabaseConfig)
	suite.NotNil(currentConfig.RedisConfig)
	suite.NotNil(currentConfig.ClickhouseConfig)
	suite.NotNil(currentConfig.ConsoleLoggerConfig)
	suite.NotNil(currentConfig.FileLoggerConfig)
	suite.NotNil(currentConfig.AccessLoggerConfig)
	suite.NotNil(currentService.AuthenticatorValue)
	suite.NotNil(currentService.DatabaseValue)
	suite.NotNil(currentService.RedisValue)
	suite.NotNil(currentService.ClickhouseValue)
	suite.NotNil(currentService.CasbinEnforcerValue)
	suite.NotNil(currentService.LoggerValue)
	suite.NotNil(currentService.ConsoleLogger)
	suite.NotNil(currentService.FileLogger)
	suite.NotNil(currentService.AccessLoggerValue)
	suite.NotNil(currentService.FormDecoder)
	suite.NotNil(currentService.Modifier)
	suite.NotNil(currentService.MapstructureDecoder)
}

func (suite *RegistrarTestSuite) TestValidatorTagNameFunc() {
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

func (suite *RegistrarTestSuite) TearDownSuite() {
	test.Shutdown()
}

func TestRegistrarTestSuite(t *testing.T) {
	suite.Run(t, new(RegistrarTestSuite))
}
