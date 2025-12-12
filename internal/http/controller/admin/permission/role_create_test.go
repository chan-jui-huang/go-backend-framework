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
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RoleCreateTestSuite struct {
	suite.Suite
}

func (suite *RoleCreateTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.AdminService.Register()
}

func (suite *RoleCreateTestSuite) Test() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	accessToken := test.AdminService.Login()

	permissionModel := &model.Permission{Name: "permission1"}
	casbinRules := []gormadapter.CasbinRule{
		{
			Ptype: "p",
			V0:    permissionModel.Name,
			V1:    "/api/test-api-1",
			V2:    "GET",
		},
	}

	err := database.NewTx().Transaction(func(tx *gorm.DB) error {
		if err := pkgPermission.Create(tx, permissionModel); err != nil {
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

	reqBody := permission.RoleCreateRequest{
		Name:          "role1",
		PermissionIds: []uint{permissionModel.Id},
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/admin/role", bytes.NewReader(reqBodyBytes))
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
	suite.Equal(len(reqBody.PermissionIds), len(data.Permissions))
	suite.NotEmpty(data.Permissions[0].Id)
	suite.NotEmpty(data.Permissions[0].Name)
	suite.NotEmpty(data.Permissions[0].CreatedAt)
	suite.NotEmpty(data.Permissions[0].UpdatedAt)
}

func (suite *RoleCreateTestSuite) TestRequestValidationFailed() {
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
		req := httptest.NewRequest("POST", "/api/admin/role", bytes.NewReader([]byte(c.reqBody)))
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

func (suite *RoleCreateTestSuite) TestWrongAccessToken() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	req := httptest.NewRequest("POST", "/api/admin/role", nil)
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

func (suite *RoleCreateTestSuite) TestCsrfMismatch() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	req := httptest.NewRequest("POST", "/api/admin/role", nil)
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

func (suite *RoleCreateTestSuite) TestAuthorizationFailed() {
	accessToken := test.AdminService.Login()
	req := httptest.NewRequest("POST", "/api/admin/role", nil)
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

func (suite *RoleCreateTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestRoleCreateTestSuite(t *testing.T) {
	suite.Run(t, new(RoleCreateTestSuite))
}
