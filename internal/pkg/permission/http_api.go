package permission

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func CreateHttpApi(tx *gorm.DB, value any) error {
	if err := tx.Table("http_apis").Create(value).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetHttpApis(tx *gorm.DB, query any, args ...any) ([]model.HttpApi, error) {
	httpApis := []model.HttpApi{}
	err := tx.Table("http_apis").
		Where(query, args...).
		Find(&httpApis).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return httpApis, nil
}
