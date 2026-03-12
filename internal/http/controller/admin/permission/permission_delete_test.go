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
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fake"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PermissionDeleteTestSuite struct {
	suite.Suite
	runtime    *test.Runtime
	permission *model.Permission
}

func (suite *PermissionDeleteTestSuite) SetupTest() {
	suite.runtime = test.NewRdbmsRuntime(suite.T())
	suite.runtime.Users.Register(fake.Admin())

	permissionModel := &model.Permission{Name: "permission1"}
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

		if err := pkgPermission.CreateCasbinRule(tx, casbinRules); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	suite.permission = permissionModel
}

func (suite *PermissionDeleteTestSuite) Test() {
	suite.runtime.Permissions.AddPermissions()
	suite.runtime.Permissions.GrantAdminToAdminUser()
	accessToken := suite.runtime.Users.Login(fake.Admin().Email, fake.Admin().Password)

	reqBody := permission.PermissionDeleteRequest{
		Ids: []uint{suite.permission.Id},
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("DELETE", "/api/admin/permission", bytes.NewReader(reqBodyBytes))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	p, err := pkgPermission.GetCasbinRules(database.NewTx(), "ptype = ? AND v0 = ?", "p", suite.permission.Name)
	if err != nil {
		panic(err)
	}

	g, err := pkgPermission.GetCasbinRules(database.NewTx(), "ptype = ? AND v1 = ?", "g", suite.permission.Name)
	if err != nil {
		panic(err)
	}

	suite.Equal(http.StatusNoContent, resp.Code)
	suite.Equal(0, len(p))
	suite.Equal(0, len(g))
}

func (suite *PermissionDeleteTestSuite) TestWrongAccessToken() {
	suite.runtime.Permissions.AddPermissions()
	suite.runtime.Permissions.GrantAdminToAdminUser()
	req := httptest.NewRequest("DELETE", "/api/admin/permission", nil)
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

func (suite *PermissionDeleteTestSuite) TestRequestValidationFailed() {
	suite.runtime.Permissions.AddPermissions()
	suite.runtime.Permissions.GrantAdminToAdminUser()
	accessToken := suite.runtime.Users.Login(fake.Admin().Email, fake.Admin().Password)

	reqBodyBytes := []byte(`{}`)
	req := httptest.NewRequest("DELETE", "/api/admin/permission", bytes.NewReader(reqBodyBytes))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
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
		"ids": "required",
	}, respBody.Context)
}

func (suite *PermissionDeleteTestSuite) TestCsrfMismatch() {
	suite.runtime.Permissions.AddPermissions()
	suite.runtime.Permissions.GrantAdminToAdminUser()
	req := httptest.NewRequest("DELETE", "/api/admin/permission", nil)
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

func (suite *PermissionDeleteTestSuite) TestAuthorizationFailed() {
	accessToken := suite.runtime.Users.Login(fake.Admin().Email, fake.Admin().Password)
	req := httptest.NewRequest("DELETE", "/api/admin/permission", nil)
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

func (suite *PermissionDeleteTestSuite) TearDownTest() {
	suite.runtime.Close()
}

func TestPermissionDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionDeleteTestSuite))
}
