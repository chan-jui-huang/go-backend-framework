package registrar

import (
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
)

type CasbinRegistrar struct{}

func (*CasbinRegistrar) Boot() {
}

func (*CasbinRegistrar) Register() {
	adapter, err := gormadapter.NewAdapterByDBUseTableName(
		deps.Database(),
		"",
		"casbin_rules",
	)
	if err != nil {
		panic(err)
	}

	m, err := model.NewModelFromString(`
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
		panic(err)
	}

	enforcer, err := casbin.NewSyncedCachedEnforcer(m, adapter)
	if err != nil {
		panic(err)
	}

	current := deps.CurrentService()
	current.CasbinEnforcerValue = enforcer
	deps.SetService(current)
}
