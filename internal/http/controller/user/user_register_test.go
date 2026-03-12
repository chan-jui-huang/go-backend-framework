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
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/suite"
)

type UserRegisterTestSuite struct {
	suite.Suite
	runtime *test.Runtime
}

func (suite *UserRegisterTestSuite) SetupSuite() {
	suite.runtime = test.NewRdbmsRuntime(suite.T())
}

func (suite *UserRegisterTestSuite) Test() {
	reqBody := user.UserRegisterRequest{
		Name:     "bob",
		Email:    "bob@test.com",
		Password: "abcABC123",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(reqBodyBytes))
	suite.runtime.HTTP.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &user.TokenData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.AccessToken)
}

func (suite *UserRegisterTestSuite) TestCsrfMismatch() {
	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader([]byte{}))
	resp := httptest.NewRecorder()
	suite.runtime.HTTP.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusForbidden, resp.Code)
	suite.Equal(response.Forbidden, respBody.Message)
	suite.Equal(response.MessageToCode[response.Forbidden], respBody.Code)
}

func (suite *UserRegisterTestSuite) TestRequestValidationFailed() {
	cases := []struct {
		reqBody  string
		expected map[string]any
	}{
		{
			reqBody: `{}`,
			expected: map[string]any{
				"name":     "required",
				"email":    "required",
				"password": "required",
			},
		},
		{
			reqBody: `{"name":"bob","email":"not-an-email","password":"abcABC123"}`,
			expected: map[string]any{
				"email": "email",
			},
		},
		{
			reqBody: `{"name":"bob","email":"bob@test.com","password":"Abc12"}`,
			expected: map[string]any{
				"password": "gte",
			},
		},
		{
			reqBody: `{"name":"bob","email":"bob@test.com","password":"ABCDEFG1"}`,
			expected: map[string]any{
				"password": "containsany",
			},
		},
		{
			reqBody: `{"name":"bob","email":"bob@test.com","password":"abcdefg1"}`,
			expected: map[string]any{
				"password": "containsany",
			},
		},
		{
			reqBody: `{"name":"bob","email":"bob@test.com","password":"abcABCdef"}`,
			expected: map[string]any{
				"password": "containsany",
			},
		},
	}

	for _, c := range cases {
		req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader([]byte(c.reqBody)))
		suite.runtime.HTTP.AddCsrfToken(req)
		resp := httptest.NewRecorder()
		suite.runtime.HTTP.ServeHTTP(resp, req)

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

func (suite *UserRegisterTestSuite) TearDownSuite() {
	suite.runtime.Close()
}

func TestUserRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(UserRegisterTestSuite))
}
