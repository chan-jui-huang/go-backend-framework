package permission_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/admin/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	pkgPermission "github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fake"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/pagination"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PermissionSearchTestSuite struct {
	suite.Suite
}

func (suite *PermissionSearchTestSuite) SetupTest() {
	test.Setup(suite.T())
	test.GetRuntime().Rdbms.Run()
	test.GetRuntime().Users.Register(fake.Admin())
}

func (suite *PermissionSearchTestSuite) Test() {
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

	test.GetRuntime().Permissions.AddPermissions()
	test.GetRuntime().Permissions.GrantAdminToAdminUser()
	accessToken := test.GetRuntime().Users.Login(fake.Admin().Email, fake.Admin().Password)

	searchRequest := permission.PermissionSearchRequest{
		Name: permissionModel.Name,
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

	req := httptest.NewRequest("GET", "/api/admin/permission?"+queryString.Encode(), nil)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.GetRuntime().HTTP.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	decoder := deps.MapstructureDecoder()
	data := &permission.PermissionSearchData{}
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.Total)
	suite.NotEmpty(data.LastPage)
	suite.NotEmpty(data.Permissions[0].Id)
	suite.NotEmpty(data.Permissions[0].Name)
	suite.NotEmpty(data.Permissions[0].CreatedAt)
	suite.NotEmpty(data.Permissions[0].UpdatedAt)
}

func (suite *PermissionSearchTestSuite) TestRequestValidationFailed() {
	test.GetRuntime().Permissions.AddPermissions()
	test.GetRuntime().Permissions.GrantAdminToAdminUser()
	accessToken := test.GetRuntime().Users.Login(fake.Admin().Email, fake.Admin().Password)

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
		req := httptest.NewRequest("GET", "/api/admin/permission"+c.query, nil)
		test.AddBearerToken(req, accessToken)
		resp := httptest.NewRecorder()
		test.GetRuntime().HTTP.ServeHTTP(resp, req)

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

func (suite *PermissionSearchTestSuite) TestWrongAccessToken() {
	test.GetRuntime().Permissions.AddPermissions()
	test.GetRuntime().Permissions.GrantAdminToAdminUser()
	req := httptest.NewRequest("GET", "/api/admin/permission", nil)
	resp := httptest.NewRecorder()
	test.GetRuntime().HTTP.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Equal(response.Unauthorized, respBody.Message)
	suite.Equal(response.MessageToCode[response.Unauthorized], respBody.Code)
}

func (suite *PermissionSearchTestSuite) TestAuthorizationFailed() {
	accessToken := test.GetRuntime().Users.Login(fake.Admin().Email, fake.Admin().Password)
	req := httptest.NewRequest("GET", "/api/admin/permission", nil)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.GetRuntime().HTTP.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusForbidden, resp.Code)
	suite.Equal(response.Forbidden, respBody.Message)
	suite.Equal(response.MessageToCode[response.Forbidden], respBody.Code)
}

func (suite *PermissionSearchTestSuite) TearDownTest() {
	test.GetRuntime().Rdbms.Reset()
	test.Shutdown()
}

func TestPermissionSearchTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionSearchTestSuite))
}
