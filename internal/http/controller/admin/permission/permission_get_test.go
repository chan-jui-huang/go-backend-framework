package permission_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PermissionGetTestSuite struct {
	suite.Suite
	permission *model.Permission
}

func (suite *PermissionGetTestSuite) SetupTest() {
	test.Setup(suite.T())
	test.GetRuntime().Rdbms.Run()
	test.GetRuntime().Users.Register(fake.Admin())

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

	suite.permission = permissionModel
}

func (suite *PermissionGetTestSuite) Test() {
	test.GetRuntime().Permissions.AddPermissions()
	test.GetRuntime().Permissions.GrantAdminToAdminUser()
	accessToken := test.GetRuntime().Users.Login(fake.Admin().Email, fake.Admin().Password)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/admin/permission/%d", suite.permission.Id), nil)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.GetRuntime().HTTP.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	decoder := deps.MapstructureDecoder()
	data := &permission.PermissionGetData{}
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.Id)
	suite.NotEmpty(data.Name)
	suite.NotEmpty(data.CreatedAt)
	suite.NotEmpty(data.UpdatedAt)
	suite.NotEqual(0, len(data.HttpApis))
	suite.NotEmpty(data.HttpApis[0].Method)
	suite.NotEmpty(data.HttpApis[0].Path)
}

func (suite *PermissionGetTestSuite) TestWrongAccessToken() {
	test.GetRuntime().Permissions.AddPermissions()
	test.GetRuntime().Permissions.GrantAdminToAdminUser()
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/admin/permission/%d", suite.permission.Id), nil)
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

func (suite *PermissionGetTestSuite) TestAuthorizationFailed() {
	accessToken := test.GetRuntime().Users.Login(fake.Admin().Email, fake.Admin().Password)
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/admin/permission/%d", suite.permission.Id), nil)
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

func (suite *PermissionGetTestSuite) TearDownTest() {
	test.GetRuntime().Rdbms.Reset()
	test.Shutdown()
}

func TestPermissionGetTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionGetTestSuite))
}
