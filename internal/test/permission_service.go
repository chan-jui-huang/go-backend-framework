package test

import (
	"fmt"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/user"
	"gorm.io/gorm"
)

type permissionService struct {
	permissions []model.Permission
	casbinRules []gormadapter.CasbinRule
	role        *model.Role
}

var PermissionService *permissionService

func NewPermissionService() *permissionService {
	permissions := []model.Permission{
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
	}
	casbinRules := []gormadapter.CasbinRule{
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
	}
	role := &model.Role{Name: "admin"}

	return &permissionService{
		permissions: permissions,
		casbinRules: casbinRules,
		role:        role,
	}
}

func (ps *permissionService) AddPermissions() {
	err := database.NewTx().Transaction(func(tx *gorm.DB) error {
		if err := permission.Create(tx, ps.permissions); err != nil {
			return err
		}

		if err := permission.CreateRole(tx, ps.role); err != nil {
			return err
		}

		rolePermissions := make([]model.RolePermission, len(ps.permissions))
		for i := 0; i < len(rolePermissions); i++ {
			rolePermissions[i].RoleId = ps.role.Id
			rolePermissions[i].PermissionId = ps.permissions[i].Id
		}
		if err := permission.CreateRolePermission(tx, rolePermissions); err != nil {
			return err
		}

		if err := permission.CreateCasbinRule(tx, ps.casbinRules); err != nil {
			panic(err)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	enforcer := deps.CasbinEnforcer()
	if err := enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}

func (ps *permissionService) GrantRoleToUser(userId uint, roleName string) {
	role, err := permission.GetRole(database.NewTx("Permissions"), "name = ?", roleName)
	if err != nil {
		panic(err)
	}
	userRole := &model.UserRole{
		UserId: userId,
		RoleId: role.Id,
	}

	casbinRules := make([]gormadapter.CasbinRule, len(role.Permissions))
	for i := 0; i < len(casbinRules); i++ {
		casbinRules[i].Ptype = "g"
		casbinRules[i].V0 = fmt.Sprintf("%d", userId)
		casbinRules[i].V1 = role.Permissions[i].Name
	}

	err = database.NewTx().Transaction(func(tx *gorm.DB) error {
		if err := permission.CreateUserRole(tx, userRole); err != nil {
			return err
		}

		if err := permission.CreateCasbinRule(tx, casbinRules); err != nil {
			panic(err)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	enforcer := deps.CasbinEnforcer()
	if err := enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}

func (ps *permissionService) GrantAdminToAdminUser() {
	u, err := user.Get(database.NewTx(), "email = ?", "admin@test.com")
	if err != nil {
		panic(err)
	}
	userRoles := []model.UserRole{
		{UserId: u.Id, RoleId: ps.role.Id},
	}
	casbinRules := make([]gormadapter.CasbinRule, len(ps.permissions))
	for i := 0; i < len(casbinRules); i++ {
		casbinRules[i].Ptype = "g"
		casbinRules[i].V0 = fmt.Sprintf("%d", u.Id)
		casbinRules[i].V1 = ps.permissions[i].Name
	}

	err = database.NewTx().Transaction(func(tx *gorm.DB) error {
		if err := permission.CreateUserRole(tx, userRoles); err != nil {
			return err
		}

		if err := permission.CreateCasbinRule(tx, casbinRules); err != nil {
			panic(err)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	enforcer := deps.CasbinEnforcer()
	if err := enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}
