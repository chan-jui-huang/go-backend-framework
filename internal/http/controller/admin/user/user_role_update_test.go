package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/admin/user"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	pkgPermission "github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fake"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserRoleUpdateTestSuite struct {
	suite.Suite
	runtime     *test.Runtime
	roles       []model.Role
	permissions []model.Permission
}

func (suite *UserRoleUpdateTestSuite) SetupTest() {
	suite.runtime = test.NewRdbmsRuntime(suite.T())
	suite.runtime.Users.Register(fake.Admin())

	roles := []model.Role{{Name: "role1"}, {Name: "role2"}}
	permissions := []model.Permission{{Name: "permission1"}, {Name: "permission2"}}
	userRole := &model.UserRole{UserId: 1}
	casbinRules := []gormadapter.CasbinRule{
		{
			Ptype: "p",
			V0:    permissions[0].Name,
			V1:    "/api/test-api-1",
			V2:    "GET",
		},
		{
			Ptype: "p",
			V0:    permissions[1].Name,
			V1:    "/api/test-api-2",
			V2:    "GET",
		},
		{
			Ptype: "g",
			V0:    "1",
			V1:    permissions[0].Name,
		},
	}

	err := database.NewTx().Transaction(func(tx *gorm.DB) error {
		if err := pkgPermission.Create(tx, permissions); err != nil {
			return err
		}

		if err := pkgPermission.CreateRole(tx, roles); err != nil {
			return err
		}

		rolePermissions := []model.RolePermission{
			{
				RoleId:       roles[0].Id,
				PermissionId: permissions[0].Id,
			},
			{
				RoleId:       roles[1].Id,
				PermissionId: permissions[0].Id,
			},
			{
				RoleId:       roles[1].Id,
				PermissionId: permissions[1].Id,
			},
		}
		if err := pkgPermission.CreateRolePermission(tx, rolePermissions); err != nil {
			return err
		}

		userRole.RoleId = roles[0].Id
		if err := pkgPermission.CreateUserRole(tx, userRole); err != nil {
			return err
		}

		if err := pkgPermission.CreateCasbinRule(tx, casbinRules); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	suite.roles = roles
	suite.permissions = permissions
}

func (suite *UserRoleUpdateTestSuite) Test() {
	accessToken := suite.runtime.AdminAPI.CreateAuthorizedAccessToken()

	reqBody := user.UserRoleUpdateRequest{
		UserId:  1,
		RoleIds: []uint{suite.roles[1].Id},
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("PUT", "/api/admin/user-role", bytes.NewReader(reqBodyBytes))
	suite.runtime.HTTP.AddCsrfToken(req)
	suite.runtime.HTTP.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	decoder := deps.MapstructureDecoder()
	data := &user.UserData{}
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.Id)
	suite.NotEmpty(data.Name)
	suite.NotEmpty(data.Email)
	suite.NotEmpty(data.CreatedAt)
	suite.NotEmpty(data.UpdatedAt)
	suite.Equal(len(reqBody.RoleIds), len(data.Roles))
	suite.Contains(reqBody.RoleIds, data.Roles[0].Id)
}

func (suite *UserRoleUpdateTestSuite) TestDeleteAllRoles() {
	accessToken := suite.runtime.AdminAPI.CreateAuthorizedAccessToken()

	reqBody := user.UserRoleUpdateRequest{
		UserId:  1,
		RoleIds: []uint{},
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("PUT", "/api/admin/user-role", bytes.NewReader(reqBodyBytes))
	suite.runtime.HTTP.AddCsrfToken(req)
	suite.runtime.HTTP.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	decoder := deps.MapstructureDecoder()
	data := &user.UserData{}
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.Id)
	suite.NotEmpty(data.Name)
	suite.NotEmpty(data.Email)
	suite.NotEmpty(data.CreatedAt)
	suite.NotEmpty(data.UpdatedAt)
	suite.Equal(len(reqBody.RoleIds), len(data.Roles))
}

func (suite *UserRoleUpdateTestSuite) TestPermissionIsRepeat() {
	accessToken := suite.runtime.AdminAPI.CreateAuthorizedAccessToken()

	reqBody := user.UserRoleUpdateRequest{
		UserId:  1,
		RoleIds: []uint{suite.roles[0].Id, suite.roles[1].Id},
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("PUT", "/api/admin/user-role", bytes.NewReader(reqBodyBytes))
	suite.runtime.HTTP.AddCsrfToken(req)
	suite.runtime.HTTP.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Equal(response.PermissionIsRepeat, respBody.Message)
	suite.Equal(response.MessageToCode[response.PermissionIsRepeat], respBody.Code)
}

func (suite *UserRoleUpdateTestSuite) TestRequestValidationFailed() {
	accessToken := suite.runtime.AdminAPI.CreateAuthorizedAccessToken()

	reqBodyBytes := []byte(`{}`)
	req := httptest.NewRequest("PUT", "/api/admin/user-role", bytes.NewReader(reqBodyBytes))
	suite.runtime.HTTP.AddCsrfToken(req)
	suite.runtime.HTTP.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Equal(response.RequestValidationFailed, respBody.Message)
	suite.Equal(response.MessageToCode[response.RequestValidationFailed], respBody.Code)
	suite.Equal(map[string]any{
		"user_id":  "required",
		"role_ids": "required",
	}, respBody.Context)
}

func (suite *UserRoleUpdateTestSuite) TestWrongAccessToken() {
	suite.runtime.AdminAPI.GrantAdminAccess()
	req := httptest.NewRequest("PUT", "/api/admin/user-role", nil)
	suite.runtime.HTTP.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Equal(response.Unauthorized, respBody.Message)
	suite.Equal(response.MessageToCode[response.Unauthorized], respBody.Code)
}

func (suite *UserRoleUpdateTestSuite) TestCsrfMismatch() {
	suite.runtime.AdminAPI.GrantAdminAccess()
	req := httptest.NewRequest("PUT", "/api/admin/user-role", nil)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusForbidden, resp.Code)
	suite.Equal(response.Forbidden, respBody.Message)
	suite.Equal(response.MessageToCode[response.Forbidden], respBody.Code)
}

func (suite *UserRoleUpdateTestSuite) TestAuthorizationFailed() {
	accessToken := suite.runtime.AdminAPI.CreateAccessToken()
	req := httptest.NewRequest("PUT", "/api/admin/user-role", nil)
	suite.runtime.HTTP.AddCsrfToken(req)
	suite.runtime.HTTP.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusForbidden, resp.Code)
	suite.Equal(response.Forbidden, respBody.Message)
	suite.Equal(response.MessageToCode[response.Forbidden], respBody.Code)
}

func (suite *UserRoleUpdateTestSuite) TearDownTest() {
	suite.runtime.Close()
}

func TestUserRoleUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(UserRoleUpdateTestSuite))
}
