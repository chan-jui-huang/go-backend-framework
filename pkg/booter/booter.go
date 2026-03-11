package booter

import (
	"os"
	"path"
	"strings"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/config"
	booterpkg "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/spf13/viper"
)

type Config = booterpkg.Config

type NewConfigFunc func() *Config
type LoadEnvFunc func()

type Registrar interface {
	Boot()
	Register()
}

type RegistrarCenter struct {
	registrars []Registrar
}

func NewConfig(rootDir string, configFileName string, debug bool) *Config {
	return booterpkg.NewConfig(rootDir, configFileName, debug)
}

func NewConfigWithCommand() *Config {
	return booterpkg.NewConfigWithCommand()
}

func NewDefaultConfig() *Config {
	return booterpkg.NewDefaultConfig()
}

func NewRegistrarCenter(registrars []Registrar) *RegistrarCenter {
	return &RegistrarCenter{registrars: registrars}
}

func (r *RegistrarCenter) GetRegistrars() []Registrar {
	return r.registrars
}

func (r *RegistrarCenter) Execute() {
	for _, registrar := range r.registrars {
		registrar.Boot()
		registrar.Register()
	}
}

type RegisterExecutor interface {
	BeforeExecute()
	Execute()
	AfterExecute()
}

func BootConfigLoader(booterConfig *Config) *config.Loader {
	byteYaml, err := os.ReadFile(path.Join(booterConfig.RootDir, booterConfig.ConfigFileName))
	if err != nil {
		panic(err)
	}
	stringYaml := os.ExpandEnv(string(byteYaml))

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(strings.NewReader(stringYaml)); err != nil {
		panic(err)
	}

	return config.NewLoader(v, booterConfig)
}

func Boot(loadEnvFunc LoadEnvFunc, newConfigFunc NewConfigFunc, registrarCenter RegisterExecutor) {
	loadEnvFunc()
	booterConfig := newConfigFunc()
	loader := BootConfigLoader(booterConfig)
	current := deps.CurrentConfig()
	current.BooterConfig = booterConfig
	deps.SetConfig(current)
	config.Registry.Set("booter", loader.BooterConfig())
	config.Registry.SetViper(loader.Viper())
	registrarCenter.BeforeExecute()
	registrarCenter.Execute()
	registrarCenter.AfterExecute()
}
