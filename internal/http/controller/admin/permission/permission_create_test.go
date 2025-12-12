package permission_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/admin/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/stretchr/testify/suite"
)

type PermissionCreateTestSuite struct {
	suite.Suite
}

func (suite *PermissionCreateTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.AdminService.Register()
}

func (suite *PermissionCreateTestSuite) Test() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	accessToken := test.AdminService.Login()

	reqBody := permission.PermissionCreateRequest{
		Name: "permission1",
		HttpApis: []struct {
			Path   string "json:\"path\" binding:\"required\""
			Method string "json:\"method\" binding:\"required\""
		}{
			{
				Path:   "/api/test-api-1",
				Method: "GET",
			},
			{
				Path:   "/api/test-api-2",
				Method: "GET",
			},
		},
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/admin/permission", bytes.NewReader(reqBodyBytes))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	decoder := service.Registry.Get("mapstructureDecoder").(func(any, any) error)
	data := &permission.PermissionCreateData{}
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.Id)
	suite.Equal(reqBody.Name, data.Name)
	suite.NotEmpty(data.CreatedAt)
	suite.NotEmpty(data.UpdatedAt)
	suite.Equal(len(reqBody.HttpApis), len(data.HttpApis))
	suite.NotEmpty(data.HttpApis[0].Method)
	suite.NotEmpty(data.HttpApis[0].Path)
}

func (suite *PermissionCreateTestSuite) TestRequestValidationFailed() {
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
				"name":      "required",
				"http_apis": "required",
			},
		},
		{
			reqBody: `{"name":"p1","http_apis":[]}`,
			expected: map[string]any{
				"http_apis": "min",
			},
		},
		{
			reqBody: `{"name":"p1","http_apis":[{"method":"GET"}]}`,
			expected: map[string]any{
				"path": "required",
			},
		},
		{
			reqBody: `{"name":"p1","http_apis":[{"path":"/api/a"}]}`,
			expected: map[string]any{
				"method": "required",
			},
		},
	}

	for _, c := range cases {
		req := httptest.NewRequest("POST", "/api/admin/permission", bytes.NewReader([]byte(c.reqBody)))
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

func (suite *PermissionCreateTestSuite) TestWrongAccessToken() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	req := httptest.NewRequest("POST", "/api/admin/permission", nil)
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

func (suite *PermissionCreateTestSuite) TestCsrfMismatch() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	req := httptest.NewRequest("POST", "/api/admin/permission", nil)
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

func (suite *PermissionCreateTestSuite) TestAuthorizationFailed() {
	accessToken := test.AdminService.Login()
	req := httptest.NewRequest("POST", "/api/admin/permission", nil)
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

func (suite *PermissionCreateTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestPermissionCreateTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionCreateTestSuite))
}
