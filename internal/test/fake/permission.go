package fake

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
)

type PermissionPreset struct {
	Permissions []model.Permission
	CasbinRules []gormadapter.CasbinRule
	RoleName    string
}

func AdminPermissionPreset() PermissionPreset {
	return PermissionPreset{
		Permissions: []model.Permission{
			{Name: "http-api-read"},
			{Name: "permission-create"},
			{Name: "permission-read"},
			{Name: "permission-update"},
			{Name: "permission-delete"},
			{Name: "permission-reload"},
			{Name: "role-create"},
			{Name: "role-read"},
			{Name: "role-update"},
			{Name: "role-delete"},
			{Name: "user-role-update"},
		},
		CasbinRules: []gormadapter.CasbinRule{
			{Ptype: "p", V0: "http-api-read", V1: "/api/admin/http-api", V2: "GET"},
			{Ptype: "p", V0: "permission-create", V1: "/api/admin/permission", V2: "POST"},
			{Ptype: "p", V0: "permission-read", V1: "/api/admin/permission", V2: "GET"},
			{Ptype: "p", V0: "permission-read", V1: "/api/admin/permission/:id", V2: "GET"},
			{Ptype: "p", V0: "permission-update", V1: "/api/admin/permission/:id", V2: "PUT"},
			{Ptype: "p", V0: "permission-delete", V1: "/api/admin/permission", V2: "DELETE"},
			{Ptype: "p", V0: "permission-reload", V1: "/api/admin/permission/reload", V2: "POST"},
			{Ptype: "p", V0: "role-create", V1: "/api/admin/role", V2: "POST"},
			{Ptype: "p", V0: "role-read", V1: "/api/admin/role", V2: "GET"},
			{Ptype: "p", V0: "role-update", V1: "/api/admin/role/:id", V2: "PUT"},
			{Ptype: "p", V0: "role-delete", V1: "/api/admin/role", V2: "DELETE"},
			{Ptype: "p", V0: "user-role-update", V1: "/api/admin/user-role", V2: "PUT"},
		},
		RoleName: "admin",
	}
}
