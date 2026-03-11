package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/user"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	pkgUser "github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/user"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/argon2"
	"github.com/mitchellh/mapstructure"
)

type userService struct {
	User         *model.User
	UserPassword string
}

var UserService *userService

func NewUserService() *userService {
	return &userService{
		User: &model.User{
			Name:  "john",
			Email: "john@test.com",
		},
		UserPassword: "abcABC123",
	}
}

func (us *userService) Register() {
	us.User.Password = argon2.MakeArgon2IdHash(us.UserPassword)
	err := pkgUser.Create(database.NewTx(), us.User)
	if err != nil {
		panic(err)
	}
}

func (us *userService) Login() string {
	userLoginRequest := user.UserLoginRequest{
		Email:    us.User.Email,
		Password: us.UserPassword,
	}
	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBody))
	AddCsrfToken(req)
	resp := httptest.NewRecorder()
	HttpHandler.ServeHTTP(resp, req)

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
