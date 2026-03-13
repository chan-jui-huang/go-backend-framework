package domain

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	pkgUser "github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/user"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fake"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/argon2"
	"gorm.io/gorm"
)

type UserFixture struct {
	db *gorm.DB
}

func NewUserFixture(db *gorm.DB) *UserFixture {
	return &UserFixture{
		db: db,
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

func (uo *UserFixture) GetByEmail(email string) *model.User {
	userModel, err := pkgUser.Get(uo.db, "email = ?", email)
	if err != nil {
		panic(err)
	}

	return userModel
}
