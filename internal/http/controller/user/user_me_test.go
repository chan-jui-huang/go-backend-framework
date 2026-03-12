package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/user"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fake"
	"github.com/stretchr/testify/suite"
)

type UserMeTestSuite struct {
	suite.Suite
	userId uint
}

func (suite *UserMeTestSuite) SetupSuite() {
	test.Setup(suite.T())
	test.GetRuntime().Rdbms.Run()
	createdUser := test.GetRuntime().Users.Register(fake.User())
	suite.userId = createdUser.Id
}

func (suite *UserMeTestSuite) Test() {
	test.GetRuntime().Permissions.AddPermissions()
	test.GetRuntime().Permissions.GrantRoleToUser(suite.userId, "admin")

	accessToken := test.GetRuntime().Users.Login(fake.User().Email, fake.User().Password)
	req := httptest.NewRequest("GET", "/api/user/me", nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.GetRuntime().HTTP.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &user.UserData{}
	decoder := deps.MapstructureDecoder()
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.Id)
	suite.NotEmpty(data.Name)
	suite.NotEmpty(data.Email)
	suite.NotEmpty(data.CreatedAt)
	suite.NotEmpty(data.UpdatedAt)
	suite.NotEmpty(data.Roles[0].Id)
	suite.NotEmpty(data.Roles[0].Name)
	suite.NotEmpty(data.Roles[0].CreatedAt)
	suite.NotEmpty(data.Roles[0].UpdatedAt)
	suite.NotEmpty(data.Roles[0].Permissions[0].Id)
	suite.NotEmpty(data.Roles[0].Permissions[0].Name)
	suite.NotEmpty(data.Roles[0].Permissions[0].CreatedAt)
	suite.NotEmpty(data.Roles[0].Permissions[0].UpdatedAt)
}

func (suite *UserMeTestSuite) TestWrongAccessToken() {
	req := httptest.NewRequest("GET", "/api/user/me", nil)
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.GetRuntime().HTTP.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Equal(response.Unauthorized, respBody.Message)
	suite.Equal(response.MessageToCode[response.Unauthorized], respBody.Code)
}

func (suite *UserMeTestSuite) TearDownSuite() {
	test.GetRuntime().Rdbms.Reset()
	test.Shutdown()
}

func TestUserMeTestSuite(t *testing.T) {
	suite.Run(t, new(UserMeTestSuite))
}
