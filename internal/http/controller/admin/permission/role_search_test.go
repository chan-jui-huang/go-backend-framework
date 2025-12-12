package permission_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/admin/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	pkgPermission "github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/chan-jui-huang/go-backend-package/pkg/pagination"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RoleSearchTestSuite struct {
	suite.Suite
}

func (suite *RoleSearchTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.AdminService.Register()
}

func (suite *RoleSearchTestSuite) Test() {
	role := &model.Role{Name: "role1"}
	permissionModel := &model.Permission{Name: "permission1"}

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

		return nil
	})
	if err != nil {
		panic(err)
	}

	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	accessToken := test.AdminService.Login()

	searchRequest := permission.RoleSearchRequest{
		Name:     role.Name,
		IsPublic: role.IsPublic,
		PaginationRequest: pagination.PaginationRequest{
			Page:    1,
			PerPage: 10,
		},
	}
	queryString := url.Values{}
	encoder := schema.NewEncoder()
	if err := encoder.Encode(searchRequest, queryString); err != nil {
		panic(err)
	}

	req := httptest.NewRequest("GET", "/api/admin/role?"+queryString.Encode(), nil)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	decoder := service.Registry.Get("mapstructureDecoder").(func(any, any) error)
	data := &permission.RoleSearchData{}
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.Total)
	suite.NotEmpty(data.LastPage)
	suite.NotEmpty(data.Roles[0].Id)
	suite.Equal(role.Name, data.Roles[0].Name)
	suite.Equal(role.IsPublic, data.Roles[0].IsPublic)
	suite.NotEmpty(data.Roles[0].CreatedAt)
	suite.NotEmpty(data.Roles[0].UpdatedAt)
	suite.NotEmpty(data.Roles[0].Permissions[0].Id)
	suite.NotEmpty(data.Roles[0].Permissions[0].Name)
	suite.NotEmpty(data.Roles[0].Permissions[0].CreatedAt)
	suite.NotEmpty(data.Roles[0].Permissions[0].UpdatedAt)
}

func (suite *RoleSearchTestSuite) TestRequestValidationFailed() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	accessToken := test.AdminService.Login()

	cases := []struct {
		query    string
		expected map[string]any
	}{
		{
			query: "",
			expected: map[string]any{
				"page":     "required",
				"per_page": "required",
			},
		},
		{
			query: "?page=-1&per_page=10",
			expected: map[string]any{
				"page": "gt",
			},
		},
		{
			query: "?page=1&per_page=1",
			expected: map[string]any{
				"per_page": "min",
			},
		},
	}

	for _, c := range cases {
		req := httptest.NewRequest("GET", "/api/admin/role"+c.query, nil)
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

func (suite *RoleSearchTestSuite) TestWrongAccessToken() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	req := httptest.NewRequest("GET", "/api/admin/role", nil)
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

func (suite *RoleSearchTestSuite) TestAuthorizationFailed() {
	accessToken := test.AdminService.Login()
	req := httptest.NewRequest("GET", "/api/admin/role", nil)
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

func (suite *RoleSearchTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestRoleSearchTestSuite(t *testing.T) {
	suite.Run(t, new(RoleSearchTestSuite))
}
