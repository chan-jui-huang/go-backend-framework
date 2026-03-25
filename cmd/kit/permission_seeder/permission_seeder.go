package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/user"
	appregistrar "github.com/chan-jui-huang/go-backend-framework/v2/internal/registrar"
	booter "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
)

type permissionSeeder struct {
	permissions []model.Permission
	casbinRules []gormadapter.CasbinRule
	roles       []model.Role
}

func (ps *permissionSeeder) addRoles() {
	currentRoles, err := permission.GetRoles(database.NewTx(), "")
	if err != nil {
		panic(err)
	}

	appendedRoles := []model.Role{}
	for _, role := range ps.roles {
		doesAppend := true
		for _, currentRole := range currentRoles {
			if role.Name == currentRole.Name {
				doesAppend = false
				break
			}
		}
		if doesAppend {
			appendedRoles = append(appendedRoles, role)
		}
	}

	if len(appendedRoles) == 0 {
		return
	}
	if err := permission.CreateRole(database.NewTx(), appendedRoles); err != nil {
		panic(err)
	}
}

func (ps *permissionSeeder) appendPermissionsToRoles() {
	if len(ps.permissions) == 0 {
		return
	}

	roleNames := lo.Map(ps.roles, func(role model.Role, _ int) string {
		return role.Name
	})

	var err error
	ps.roles, err = permission.GetRoles(database.NewTx(), "name IN ?", roleNames)
	if err != nil {
		panic(err)
	}

	u, err := user.Get(database.NewTx(), "email = ?", "admin@admin.com")
	if err != nil {
		panic(err)
	}

	userRoles := []model.UserRole{}
	for _, role := range ps.roles {
		userRole := model.UserRole{
			UserId: u.Id,
			RoleId: role.Id,
		}
		userRoles = append(userRoles, userRole)
	}

	for _, permission := range ps.permissions {
		casbinRule := gormadapter.CasbinRule{
			Ptype: "g",
			V0:    fmt.Sprintf("%d", u.Id),
			V1:    permission.Name,
		}
		ps.casbinRules = append(ps.casbinRules, casbinRule)
	}

	err = database.NewTx().Transaction(func(tx *gorm.DB) error {
		if err := permission.Create(tx, ps.permissions); err != nil {
			return err
		}

		rolePermissions := []model.RolePermission{}
		for i := 0; i < len(ps.roles); i++ {
			for j := 0; j < len(ps.permissions); j++ {
				rolePermission := model.RolePermission{
					RoleId:       ps.roles[i].Id,
					PermissionId: ps.permissions[j].Id,
				}
				rolePermissions = append(rolePermissions, rolePermission)
			}
		}

		if err := permission.CreateRolePermission(tx, rolePermissions); err != nil {
			return err
		}

		if err := permission.CreateUserRole(tx, userRoles); err != nil {
			return err
		}

		if err := permission.CreateCasbinRule(tx, ps.casbinRules); err != nil {
			return err
		}

		enforcer := deps.CasbinEnforcer()
		if err := enforcer.LoadPolicy(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	fxApp := fx.New(
		fx.Supply(booter.NewDefaultConfig()),
		fx.Provide(
			appregistrar.NewConfigLoader,
			appregistrar.NewDatabaseConfig,
			appregistrar.NewDatabase,
			appregistrar.NewLoggerConfigs,
			appregistrar.NewLoggers,
			appregistrar.NewCasbinEnforcer,
		),
		fx.Invoke(
			appregistrar.RegisterConfigDependencies,
			appregistrar.RegisterServiceDependencies,
		),
	)
	if err := fxApp.Err(); err != nil {
		panic(err)
	}

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

	currentPermissions, err := permission.GetMany(database.NewTx(), "")
	if err != nil {
		panic(err)
	}
	currnetCasbinRules, err := permission.GetCasbinRules(database.NewTx(), "ptype = ?", "p")
	if err != nil {
		panic(err)
	}

	appendedPermissions := []model.Permission{}
	for _, permission := range permissions {
		doesAppend := true
		for _, currentPermission := range currentPermissions {
			if permission.Name == currentPermission.Name {
				doesAppend = false
				break
			}
		}
		if doesAppend {
			appendedPermissions = append(appendedPermissions, permission)
		}
	}

	appendedCasbinRules := []gormadapter.CasbinRule{}
	for _, casbinRule := range casbinRules {
		doesAppend := true
		for _, currnetCasbinRule := range currnetCasbinRules {
			if casbinRule.V0 == currnetCasbinRule.V0 && casbinRule.V1 == currnetCasbinRule.V1 && casbinRule.V2 == currnetCasbinRule.V2 {
				doesAppend = false
				break
			}
		}
		if doesAppend {
			appendedCasbinRules = append(appendedCasbinRules, casbinRule)
		}
	}

	permissionSeeder := &permissionSeeder{
		permissions: appendedPermissions,
		casbinRules: appendedCasbinRules,
		roles: []model.Role{
			{Name: "admin"},
		},
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "add-roles",
				Usage: "add roles",
				Action: func(cCtx *cli.Context) error {
					permissionSeeder.addRoles()
					return nil
				},
			},
			{
				Name:  "append-permissions-to-roles",
				Usage: "append permissions to roles",
				Action: func(cCtx *cli.Context) error {
					permissionSeeder.appendPermissionsToRoles()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
