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
	runtime *test.Runtime
	userId  uint
}

func (suite *UserMeTestSuite) SetupSuite() {
	suite.runtime = test.NewRdbmsRuntime(suite.T())
	createdUser := suite.runtime.Users.Register(fake.User())
	suite.userId = createdUser.Id
}

func (suite *UserMeTestSuite) Test() {
	suite.runtime.Permissions.AddPermissions()
	suite.runtime.Permissions.GrantRoleToUser(suite.userId, "admin")

	accessToken := suite.runtime.UserAPI.CreateAccessToken()
	req := httptest.NewRequest("GET", "/api/user/me", nil)
	suite.runtime.HTTP.AddCsrfToken(req)
	suite.runtime.HTTP.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

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
	suite.runtime.HTTP.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Equal(response.Unauthorized, respBody.Message)
	suite.Equal(response.MessageToCode[response.Unauthorized], respBody.Code)
}

func (suite *UserMeTestSuite) TearDownSuite() {
	suite.runtime.Close()
}

func TestUserMeTestSuite(t *testing.T) {
	suite.Run(t, new(UserMeTestSuite))
}
