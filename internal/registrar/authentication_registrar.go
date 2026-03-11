package registrar

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/authentication"
)

type AuthenticationRegistrar struct {
	config authentication.Config
}

func (ar *AuthenticationRegistrar) Boot() {
	config.Registry.Register("authentication.authenticator", &ar.config)
}

func (ar *AuthenticationRegistrar) Register() {
	authenticator, err := authentication.NewAuthenticator(ar.config)
	if err != nil {
		panic(err)
	}

	current := deps.CurrentConfig()
	current.AuthenticationConfig = &ar.config
	deps.SetConfig(current)

	serviceState := deps.CurrentService()
	serviceState.AuthenticatorValue = authenticator
	deps.SetService(serviceState)
}
