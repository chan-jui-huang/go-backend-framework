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
	"github.com/stretchr/testify/suite"
)

type UserUpdatePasswordTestSuite struct {
	suite.Suite
}

func (suite *UserUpdatePasswordTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.UserService.Register()
}

func (suite *UserUpdatePasswordTestSuite) Test() {
	accessToken := test.UserService.Login()
	reqBody := user.UserUpdatePasswordRequest{
		CurrentPassword: "abcABC123",
		Password:        "abcABC000",
		ConfirmPassword: "abcABC000",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("PUT", "/api/user/password", bytes.NewReader(reqBodyBytes))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	suite.Equal(http.StatusNoContent, resp.Code)
}

func (suite *UserUpdatePasswordTestSuite) TestWrongAccessToken() {
	req := httptest.NewRequest("PUT", "/api/user/password", nil)
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

func (suite *UserUpdatePasswordTestSuite) TestCsrfMismatch() {
	req := httptest.NewRequest("PUT", "/api/user/password", nil)
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

func (suite *UserUpdatePasswordTestSuite) TestRequestValidationFailed() {
	accessToken := test.UserService.Login()
	cases := []struct {
		reqBody  string
		expected map[string]any
	}{
		{
			reqBody: `{}`,
			expected: map[string]any{
				"current_password": "required",
				"password":         "required",
				"confirm_password": "required",
			},
		},
		{
			reqBody: `{"current_password":"abcABC123","password":"Abc12","confirm_password":"Abc12"}`,
			expected: map[string]any{
				"password": "gte",
			},
		},
		{
			reqBody: `{"current_password":"abcABC123","password":"ABCDEFG1","confirm_password":"ABCDEFG1"}`,
			expected: map[string]any{
				"password": "containsany",
			},
		},
		{
			reqBody: `{"current_password":"abcABC123","password":"abcdefg1","confirm_password":"abcdefg1"}`,
			expected: map[string]any{
				"password": "containsany",
			},
		},
		{
			reqBody: `{"current_password":"abcABC123","password":"abcABCdef","confirm_password":"abcABCdef"}`,
			expected: map[string]any{
				"password": "containsany",
			},
		},
		{
			reqBody: `{"current_password":"abcABC123","password":"abcABC123","confirm_password":"abcABC124"}`,
			expected: map[string]any{
				"confirm_password": "eqfield",
			},
		},
	}

	for _, c := range cases {
		req := httptest.NewRequest("PUT", "/api/user/password", bytes.NewReader([]byte(c.reqBody)))
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
		suite.Equal(c.expected, respBody.Context)
	}
}

func (suite *UserUpdatePasswordTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestUserUpdatePasswordTestSuite(t *testing.T) {
	suite.Run(t, new(UserUpdatePasswordTestSuite))
}
