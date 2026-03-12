package permission_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fake"
	"github.com/stretchr/testify/suite"
)

type PermissionReloadTestSuite struct {
	suite.Suite
	runtime *test.Runtime
}

func (suite *PermissionReloadTestSuite) SetupTest() {
	suite.runtime = test.NewRdbmsRuntime(suite.T())
	suite.runtime.Users.Register(fake.Admin())
}

func (suite *PermissionReloadTestSuite) Test() {
	suite.runtime.Permissions.AddPermissions()
	suite.runtime.Permissions.GrantAdminToAdminUser()
	accessToken := suite.runtime.Users.Login(fake.Admin().Email, fake.Admin().Password)

	req := httptest.NewRequest("POST", "/api/admin/permission/reload", nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	suite.Equal(http.StatusNoContent, resp.Code)
}

func (suite *PermissionReloadTestSuite) TestWrongAccessToken() {
	suite.runtime.Permissions.AddPermissions()
	suite.runtime.Permissions.GrantAdminToAdminUser()
	req := httptest.NewRequest("POST", "/api/admin/permission/reload", nil)
	test.AddCsrfToken(req)
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

func (suite *PermissionReloadTestSuite) TestCsrfMismatch() {
	suite.runtime.Permissions.AddPermissions()
	suite.runtime.Permissions.GrantAdminToAdminUser()
	req := httptest.NewRequest("POST", "/api/admin/permission/reload", nil)
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

func (suite *PermissionReloadTestSuite) TestAuthorizationFailed() {
	accessToken := suite.runtime.Users.Login(fake.Admin().Email, fake.Admin().Password)
	req := httptest.NewRequest("POST", "/api/admin/permission/reload", nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
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

func (suite *PermissionReloadTestSuite) TearDownTest() {
	suite.runtime.Close()
}

func TestPermissionReloadTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionReloadTestSuite))
}
