package seeder

import (
	"fmt"

	"gorm.io/gorm"
)

type runSeederFunc func(tx *gorm.DB) error

type SeederExecutor struct {
	database       *gorm.DB
	order          []string
	runSeederFuncs map[string]runSeederFunc
}

func NewSeederExecutor(database *gorm.DB, httpApiSeeder *HttpApiSeeder) *SeederExecutor {
	order := []string{
		"httpApi",
		"user",
	}
	runSeederFuncs := map[string]runSeederFunc{
		"httpApi": httpApiSeeder.Run,
		"user":    runUserSeeder,
	}

	return &SeederExecutor{
		database:       database,
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

	err := se.database.Transaction(func(tx *gorm.DB) error {
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
