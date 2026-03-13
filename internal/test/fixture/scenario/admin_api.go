package scenario

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fake"
	domainfixture "github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fixture/domain"
)

type AdminAPI struct {
	users       *domainfixture.UserFixture
	permissions *domainfixture.PermissionFixture
	userAPI     *UserAPI
}

func NewAdminAPI(users *domainfixture.UserFixture, permissions *domainfixture.PermissionFixture, userAPI *UserAPI) *AdminAPI {
	return &AdminAPI{
		users:       users,
		permissions: permissions,
		userAPI:     userAPI,
	}
}

func (api *AdminAPI) GrantAdminAccess() {
	adminUser := api.users.GetByEmail(fake.Admin().Email)

	api.permissions.AddPermissions()
	api.permissions.GrantAdminToUser(adminUser.Id)
}

func (api *AdminAPI) CreateAccessToken() string {
	adminInput := fake.Admin()

	return api.userAPI.Login(adminInput.Email, adminInput.Password)
}

func (api *AdminAPI) CreateAuthorizedAccessToken() string {
	api.GrantAdminAccess()

	return api.CreateAccessToken()
}
