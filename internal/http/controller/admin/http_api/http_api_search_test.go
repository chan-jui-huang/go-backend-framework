package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	httpapi "github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/admin/http_api"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/chan-jui-huang/go-backend-package/pkg/pagination"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/suite"
)

type HttpApiSearchTestSuite struct {
	suite.Suite
}

func (suite *HttpApiSearchTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.AdminService.Register()
}

func (suite *HttpApiSearchTestSuite) Test() {
	test.PermissionService.AddPermissions()
	test.PermissionService.GrantAdminToAdminUser()
	accessToken := test.AdminService.Login()

	httpApi := &model.HttpApi{
		Method: "GET",
		Path:   "/api/test-api",
	}
	if err := permission.CreateHttpApi(database.NewTx(), httpApi); err != nil {
		panic(err)
	}

	searchRequest := httpapi.HttpApiSearchRequest{
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

	req := httptest.NewRequest("GET", "/api/admin/http-api?"+queryString.Encode(), nil)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	decoder := service.Registry.Get("mapstructureDecoder").(func(any, any) error)
	data := &httpapi.HttpApiSearchData{}
	if err := decoder(respBody.Data, &data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code, resp)
	suite.NotEmpty(data.Total)
	suite.NotEmpty(data.LastPage)
	suite.NotEmpty(data.HttpApis[0].Id)
	suite.NotEmpty(data.HttpApis[0].Method)
	suite.NotEmpty(data.HttpApis[0].Path)
	suite.NotEmpty(data.HttpApis[0].CreatedAt)
	suite.NotEmpty(data.HttpApis[0].UpdatedAt)
}

func (suite *HttpApiSearchTestSuite) TestRequestValidationFailed() {
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
		req := httptest.NewRequest("GET", "/api/admin/http-api"+c.query, nil)
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

func (suite *HttpApiSearchTestSuite) TestWrongAccessToken() {
	req := httptest.NewRequest("GET", "/api/admin/http-api", nil)
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

func (suite *HttpApiSearchTestSuite) TestAuthorizationFailed() {
	accessToken := test.AdminService.Login()
	req := httptest.NewRequest("GET", "/api/admin/http-api", nil)
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

func (suite *HttpApiSearchTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestHttpApiSearchTestSuite(t *testing.T) {
	suite.Run(t, new(HttpApiSearchTestSuite))
}
