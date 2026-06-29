package permission

import (
	"github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func CreateCasbinRule(tx *gorm.DB, value any) error {
	if err := tx.Table("casbin_rules").Create(value).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetCasbinRules(tx *gorm.DB, query any, args ...any) ([]gormadapter.CasbinRule, error) {
	casbinRules := []gormadapter.CasbinRule{}
	err := tx.Table("casbin_rules").
		Where(query, args...).
		Find(&casbinRules).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return casbinRules, nil
}

func UpdateCasbinRule(tx *gorm.DB, values map[string]any, query any, args ...any) (int64, error) {
	db := tx.Table("casbin_rules").
		Where(query, args...).
		Updates(values)
	if err := db.Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return db.RowsAffected, nil
}

func DeleteCasbinRule(tx *gorm.DB, query any, args ...any) error {
	err := tx.Table("casbin_rules").
		Where(query, args...).
		Delete(&struct{}{}).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
