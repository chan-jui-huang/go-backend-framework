package database

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"gorm.io/gorm"
)

func NewTx(associations ...string) *gorm.DB {
	tx := deps.Database()
	for _, association := range associations {
		tx = tx.Preload(association)
	}

	return tx
}

func NewTxByTable(table string, associations ...string) *gorm.DB {
	database := deps.Database()
	tx := database.Table(table)

	for _, association := range associations {
		tx = tx.Preload(association)
	}

	return tx
}
