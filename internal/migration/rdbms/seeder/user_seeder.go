package seeder

import (
	"fmt"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/argon2"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/random"
	"gorm.io/gorm"
)

func runUserSeeder(tx *gorm.DB) error {
	password := random.RandomString(16)
	user := &model.User{
		Name:     "admin",
		Password: argon2.MakeArgon2IdHash(password),
		Email:    "admin@admin.com",
	}

	if err := tx.Table("users").Create(user).Error; err != nil {
		return err
	}

	fmt.Printf("admin user password is [%s]\n", password)

	return nil
}
