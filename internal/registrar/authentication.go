package registrar

import (
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/authentication"
	booterconfig "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
)

func NewAuthenticationConfig(loader *booterconfig.Loader) *authentication.Config {
	config := &authentication.Config{}
	loader.Unmarshal("authentication.authenticator", config)

	return config
}

func NewAuthenticator(config *authentication.Config) (*authentication.Authenticator, error) {
	return authentication.NewAuthenticator(*config)
}
