package permission_test

import (
	"bytes"
	"encoding/json"
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
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RoleDeleteTestSuite struct {
	suite.Suite
}

func (suite *RoleDeleteTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.AdminService.Register()
}

func (suite *RoleDeleteTestSuite) Test() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	accessToken := test.AdminService.Login()

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

	reqBody := permission.PermissionDeleteRequest{
		Ids: []uint{role.Id},
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("DELETE", "/api/admin/role", bytes.NewReader(reqBodyBytes))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	suite.Equal(http.StatusNoContent, resp.Code)
}

func (suite *RoleDeleteTestSuite) TestRequestValidationFailed() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	accessToken := test.AdminService.Login()

	reqBodyBytes := []byte(`{}`)
	req := httptest.NewRequest("DELETE", "/api/admin/role", bytes.NewReader(reqBodyBytes))
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
	suite.Equal(map[string]any{
		"ids": "required",
	}, respBody.Context)
}

func (suite *RoleDeleteTestSuite) TestWrongAccessToken() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	req := httptest.NewRequest("DELETE", "/api/admin/role", nil)
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

func (suite *RoleDeleteTestSuite) TestCsrfMismatch() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	req := httptest.NewRequest("DELETE", "/api/admin/role", nil)
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

func (suite *RoleDeleteTestSuite) TestAuthorizationFailed() {
	accessToken := test.AdminService.Login()
	req := httptest.NewRequest("DELETE", "/api/admin/role", nil)
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

func (suite *RoleDeleteTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestRoleDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(RoleDeleteTestSuite))
}
