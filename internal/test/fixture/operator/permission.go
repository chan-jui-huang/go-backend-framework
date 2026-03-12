package operator

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/user"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fake"
	"gorm.io/gorm"
)

type PermissionFixture struct {
	enforcer *casbin.SyncedCachedEnforcer
	db       *gorm.DB
}

func NewPermissionFixture(enforcer *casbin.SyncedCachedEnforcer, db *gorm.DB) *PermissionFixture {
	return &PermissionFixture{
		enforcer: enforcer,
		db:       db,
	}
}

func (ps *PermissionFixture) AddPermissions() {
	preset := fake.AdminPermissionPreset()
	role := &model.Role{Name: preset.RoleName}

	err := ps.db.Transaction(func(tx *gorm.DB) error {
		if err := permission.Create(tx, preset.Permissions); err != nil {
			return err
		}

		if err := permission.CreateRole(tx, role); err != nil {
			return err
		}

		rolePermissions := make([]model.RolePermission, len(preset.Permissions))
		for i := 0; i < len(rolePermissions); i++ {
			rolePermissions[i].RoleId = role.Id
			rolePermissions[i].PermissionId = preset.Permissions[i].Id
		}
		if err := permission.CreateRolePermission(tx, rolePermissions); err != nil {
			return err
		}

		if err := permission.CreateCasbinRule(tx, preset.CasbinRules); err != nil {
			panic(err)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	if err := ps.enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}

func (ps *PermissionFixture) GrantRoleToUser(userId uint, roleName string) {
	role, err := permission.GetRole(ps.tx("Permissions"), "name = ?", roleName)
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

	err = ps.db.Transaction(func(tx *gorm.DB) error {
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

	if err := ps.enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}

func (ps *PermissionFixture) GrantAdminToAdminUser() {
	adminUser := fake.Admin()
	preset := fake.AdminPermissionPreset()
	u, err := user.Get(ps.tx(), "email = ?", adminUser.Email)
	if err != nil {
		panic(err)
	}

	ps.GrantRoleToUser(u.Id, preset.RoleName)
}

func (ps *PermissionFixture) tx(associations ...string) *gorm.DB {
	tx := ps.db
	for _, association := range associations {
		tx = tx.Preload(association)
	}

	return tx
}
