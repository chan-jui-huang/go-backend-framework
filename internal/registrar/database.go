package registrar

import (
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/database"
	"gorm.io/gorm"
)

func NewDatabaseConfig(loader *config.Loader) *database.Config {
	config := &database.Config{}
	loader.Unmarshal("database", config)

	return config
}

func NewDatabase(config *database.Config) *gorm.DB {
	return database.New(*config)
}
