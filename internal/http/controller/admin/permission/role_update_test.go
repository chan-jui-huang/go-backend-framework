package permission_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/admin/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	pkgPermission "github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RoleUpdateTestSuite struct {
	suite.Suite
	role       *model.Role
	permission *model.Permission
}

func (suite *RoleUpdateTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.AdminService.Register()

	role := &model.Role{Name: "role1"}
	permissionModel := &model.Permission{Name: "permission1"}
	userRole := &model.UserRole{UserId: 1}
	casbinRules := []gormadapter.CasbinRule{
		{
			Ptype: "p",
			V0:    permissionModel.Name,
			V1:    "/api/test-api-1",
			V2:    "GET",
		},
		{
			Ptype: "g",
			V0:    "1",
			V1:    permissionModel.Name,
		},
	}

	err := database.NewTx().Transaction(func(tx *gorm.DB) error {
		if err := pkgPermission.Create(tx, permissionModel); err != nil {
			return err
		}

		if err := pkgPermission.CreateRole(tx, role); err != nil {
			return err
		}

		rolePermission := &model.RolePermission{
			RoleId:       role.Id,
			PermissionId: permissionModel.Id,
		}
		if err := pkgPermission.CreateRolePermission(tx, rolePermission); err != nil {
			return err
		}

		userRole.RoleId = role.Id
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

	suite.role = role
	suite.permission = permissionModel
}

func (suite *RoleUpdateTestSuite) Test() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	accessToken := test.AdminService.Login()

	reqBody := permission.RoleUpdateRequest{
		Name:          "role2",
		IsPublic:      true,
		PermissionIds: []uint{suite.permission.Id},
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/admin/role/%d", suite.role.Id), bytes.NewReader(reqBodyBytes))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	decoder := service.Registry.Get("mapstructureDecoder").(func(any, any) error)
	data := &permission.RoleData{}
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.Id)
	suite.Equal(reqBody.Name, data.Name)
	suite.Equal(reqBody.IsPublic, data.IsPublic)
	suite.NotEmpty(data.CreatedAt)
	suite.NotEmpty(data.UpdatedAt)
	suite.NotEmpty(data.Permissions[0].Id)
	suite.NotEmpty(data.Permissions[0].Name)
	suite.NotEmpty(data.Permissions[0].CreatedAt)
	suite.NotEmpty(data.Permissions[0].UpdatedAt)
}

func (suite *RoleUpdateTestSuite) TestRequestValidationFailed() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	accessToken := test.AdminService.Login()

	cases := []struct {
		reqBody  string
		expected map[string]any
	}{
		{
			reqBody: `{}`,
			expected: map[string]any{
				"name":           "required",
				"permission_ids": "required",
			},
		},
		{
			reqBody: `{"name":"role1","permission_ids":[]}`,
			expected: map[string]any{
				"permission_ids": "min",
			},
		},
	}

	for _, c := range cases {
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/admin/role/%d", suite.role.Id), bytes.NewReader([]byte(c.reqBody)))
		test.AddCsrfToken(req)
		test.AddBearerToken(req, accessToken)
		resp := httptest.NewRecorder()
		test.HttpHandler.ServeHTTP(resp, req)

		respBody := &response.ErrorResponse{}
		if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
			panic(err)
		}

		suite.Equal(http.StatusBadRequest, resp.Code)
		suite.Equal(response.RequestValidationFailed, respBody.Message)
		suite.Equal(response.MessageToCode[response.RequestValidationFailed], respBody.Code)
		suite.Equal(c.expected, respBody.Context)
	}
}

func (suite *RoleUpdateTestSuite) TestWrongAccessToken() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/admin/role/%d", suite.role.Id), nil)
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Equal(response.Unauthorized, respBody.Message)
	suite.Equal(response.MessageToCode[response.Unauthorized], respBody.Code)
}

func (suite *RoleUpdateTestSuite) TestCsrfMismatch() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/admin/role/%d", suite.role.Id), nil)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusForbidden, resp.Code)
	suite.Equal(response.Forbidden, respBody.Message)
	suite.Equal(response.MessageToCode[response.Forbidden], respBody.Code)
}

func (suite *RoleUpdateTestSuite) TestAuthorizationFailed() {
	accessToken := test.AdminService.Login()
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/admin/role/%d", suite.role.Id), nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusForbidden, resp.Code)
	suite.Equal(response.Forbidden, respBody.Message)
	suite.Equal(response.MessageToCode[response.Forbidden], respBody.Code)
}

func (suite *RoleUpdateTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestRoleUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(RoleUpdateTestSuite))
}
