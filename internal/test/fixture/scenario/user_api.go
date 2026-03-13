package scenario

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/user"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fake"
	domainfixture "github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fixture/domain"
	httpfixture "github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fixture/http"
	"github.com/mitchellh/mapstructure"
)

type UserAPI struct {
	http  *httpfixture.Handler
	users *domainfixture.UserFixture
}

func NewUserAPI(httpHandler *httpfixture.Handler, users *domainfixture.UserFixture) *UserAPI {
	return &UserAPI{
		http:  httpHandler,
		users: users,
	}
}

func (api *UserAPI) Login(email string, password string) string {
	userLoginRequest := user.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBody))
	api.http.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	api.http.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}
	data := &user.TokenData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	return data.AccessToken
}

func (api *UserAPI) CreateAccessToken() string {
	userInput := fake.User()

	return api.Login(userInput.Email, userInput.Password)
}
