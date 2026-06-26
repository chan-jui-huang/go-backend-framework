package database

import (
	"gorm.io/gorm"
)

func NewTx(database *gorm.DB, associations ...string) *gorm.DB {
	tx := database
	for _, association := range associations {
		tx = tx.Preload(association)
	}

	return tx
}

func NewTxByTable(database *gorm.DB, table string, associations ...string) *gorm.DB {
	tx := database.Table(table)

	for _, association := range associations {
		tx = tx.Preload(association)
	}

	return tx
}
