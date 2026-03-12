package operator

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/user"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	pkgUser "github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/user"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fake"
	httpfixture "github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fixture/http"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/argon2"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type UserFixture struct {
	http *httpfixture.Handler
	db   *gorm.DB
}

func NewUserFixture(httpHandler *httpfixture.Handler, db *gorm.DB) *UserFixture {
	return &UserFixture{
		http: httpHandler,
		db:   db,
	}
}

func (uo *UserFixture) Register(input fake.UserInput) *model.User {
	userModel := &model.User{
		Name:  input.Name,
		Email: input.Email,
	}
	userModel.Password = argon2.MakeArgon2IdHash(input.Password)

	err := pkgUser.Create(uo.db, userModel)
	if err != nil {
		panic(err)
	}

	return userModel
}

func (uo *UserFixture) Login(email string, password string) string {
	userLoginRequest := user.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBody))
	uo.http.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	uo.http.ServeHTTP(resp, req)

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
