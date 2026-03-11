package registrar

import (
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewCasbinEnforcer(database *gorm.DB) (*casbin.SyncedCachedEnforcer, error) {
	adapter, err := gormadapter.NewAdapterByDBUseTableName(database, "", "casbin_rules")
	if err != nil {
		return nil, err
	}

	modelConfig, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act
`)
	if err != nil {
		return nil, err
	}

	return casbin.NewSyncedCachedEnforcer(modelConfig, adapter)
}
