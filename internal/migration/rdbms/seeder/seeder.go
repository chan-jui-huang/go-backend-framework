package seeder

import (
	"fmt"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"gorm.io/gorm"
)

type runSeederFunc func(tx *gorm.DB) error

type SeederExecutor struct {
	order          []string
	runSeederFuncs map[string]runSeederFunc
}

func NewSeederExecutor() *SeederExecutor {
	order := []string{
		"httpApi",
		"user",
	}
	runSeederFuncs := map[string]runSeederFunc{
		"httpApi": runHttpApiSeeder,
		"user":    runUserSeeder,
	}

	return &SeederExecutor{
		order:          order,
		runSeederFuncs: runSeederFuncs,
	}
}

func (se *SeederExecutor) ShowSeeders() {
	for _, seeder := range se.order {
		fmt.Println(seeder)
	}
}

func (se *SeederExecutor) Run(seeders []string) {
	if len(seeders) == 1 && seeders[0] == "" {
		seeders = se.order
	}

	database := deps.Database()
	err := database.Transaction(func(tx *gorm.DB) error {
		for _, seeder := range seeders {
			if fn, ok := se.runSeederFuncs[seeder]; ok {
				if err := fn(tx); err != nil {
					fmt.Printf("[%s] execute failed\n", seeder)
					return err
				}
			} else {
				fmt.Printf("[%s] does not exist\n", seeder)
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
