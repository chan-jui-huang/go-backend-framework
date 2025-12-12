package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/user"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/stretchr/testify/suite"
)

type UserUpdateTestSuite struct {
	suite.Suite
}

func (suite *UserUpdateTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.UserService.Register()
}

func (suite *UserUpdateTestSuite) Test() {
	accessToken := test.UserService.Login()
	reqBody := user.UserUpdateRequest{
		Name:  "bob",
		Email: "bob@test.com",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("PUT", "/api/user", bytes.NewReader(reqBodyBytes))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &user.UserData{}
	decoder := service.Registry.Get("mapstructureDecoder").(func(any, any) error)
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.Id)
	suite.Equal(reqBody.Name, data.Name)
	suite.Equal(reqBody.Email, data.Email)
	suite.NotEmpty(data.CreatedAt)
	suite.NotEmpty(data.UpdatedAt)
}

func (suite *UserUpdateTestSuite) TestWrongAccessToken() {
	req := httptest.NewRequest("PUT", "/api/user", nil)
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Equal(response.Unauthorized, respBody.Message)
	suite.Equal(response.MessageToCode[response.Unauthorized], respBody.Code)
}

func (suite *UserUpdateTestSuite) TestCsrfMismatch() {
	req := httptest.NewRequest("PUT", "/api/user", nil)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusForbidden, resp.Code)
	suite.Equal(response.Forbidden, respBody.Message)
	suite.Equal(response.MessageToCode[response.Forbidden], respBody.Code)
}

func (suite *UserUpdateTestSuite) TestRequestValidationFailed() {
	accessToken := test.UserService.Login()
	reqBodyBytes := []byte(`{}`)
	req := httptest.NewRequest("PUT", "/api/user", bytes.NewReader(reqBodyBytes))
	test.AddBearerToken(req, accessToken)
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Equal(response.RequestValidationFailed, respBody.Message)
	suite.Equal(response.MessageToCode[response.RequestValidationFailed], respBody.Code)
	suite.Equal(map[string]any{
		"name":  "required",
		"email": "required",
	}, respBody.Context)
}

func (suite *UserUpdateTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestUserUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(UserUpdateTestSuite))
}
