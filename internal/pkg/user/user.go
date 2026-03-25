package user

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func Create(tx *gorm.DB, user any) error {
	if err := tx.Table("users").Create(user).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Get(tx *gorm.DB, query any, args ...any) (*model.User, error) {
	user := &model.User{}
	err := tx.Table("users").Where(query, args...).
		First(user).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return user, nil
}

func Update(tx *gorm.DB, id uint, values map[string]any) (int64, error) {
	db := tx.Model(&model.User{}).
		Where("id = ?", id).
		Updates(values)
	if err := db.Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return db.RowsAffected, nil
}
