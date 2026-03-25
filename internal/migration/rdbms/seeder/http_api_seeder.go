package seeder

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/route"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/permission"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func runHttpApiSeeder(tx *gorm.DB) error {
	gin.SetMode(gin.ReleaseMode)
	httpApis, err := permission.GetHttpApis(tx, "")
	if err != nil {
		return err
	}

	engine := gin.New()
	routers := []route.Router{
		route.NewApiRouter(engine),
		route.NewSwaggerRouter(engine),
	}
	for _, router := range routers {
		router.AttachRoutes()
	}

	newHttpApis := []model.HttpApi{}
	doesNotInsert := true
	for _, routeInfo := range engine.Routes() {
		for _, httpApi := range httpApis {
			if routeInfo.Method == httpApi.Method && routeInfo.Path == httpApi.Path {
				doesNotInsert = false
				break
			}
		}
		if doesNotInsert {
			newHttpApis = append(newHttpApis, model.HttpApi{Method: routeInfo.Method, Path: routeInfo.Path})
		}
		doesNotInsert = true
	}

	if len(newHttpApis) == 0 {
		return nil
	}

	if err := permission.CreateHttpApi(tx, &newHttpApis); err != nil {
		return err
	}

	return nil
}
