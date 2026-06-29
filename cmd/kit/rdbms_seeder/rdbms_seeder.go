package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/migration/rdbms/seeder"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/registrar"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
)

type RdbmsSeederRunner struct {
	seederExecutor *seeder.SeederExecutor
}

func NewRdbmsSeederRunner(seederExecutor *seeder.SeederExecutor) *RdbmsSeederRunner {
	return &RdbmsSeederRunner{seederExecutor: seederExecutor}
}

func (r *RdbmsSeederRunner) Run(args []string) error {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "show",
				Usage: "show all seeders",
				Action: func(cCtx *cli.Context) error {
					r.seederExecutor.ShowSeeders()
					return nil
				},
			},
			{
				Name:  "run",
				Usage: "Run the seeders. EX: database_seeder run seeder1,seeder2 (run specific seeders). database_seeder run (run all seeders).",
				Action: func(cCtx *cli.Context) error {
					r.seederExecutor.Run(strings.Split(cCtx.Args().First(), ","))
					return nil
				},
			},
		},
	}

	return app.Run(args)
}

func main() {
	var runner *RdbmsSeederRunner

	fxApp := fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
		fx.Supply(booter.NewDefaultConfig()),
		fx.Provide(
			registrar.NewConfigLoader,
			registrar.NewDatabaseConfig,
			registrar.NewDatabase,
			registrar.NewLoggerConfigs,
			registrar.NewLoggers,
			registrar.NewAuthenticationConfig,
			registrar.NewAuthenticator,
			registrar.NewCasbinEnforcer,
			seeder.NewHttpApiSeeder,
			seeder.NewSeederExecutor,
			NewRdbmsSeederRunner,
		),
		fx.Populate(&runner),
	)
	if err := fxApp.Err(); err != nil {
		log.Fatal(err)
	}

	startCtx, cancelStart := context.WithTimeout(context.Background(), 15*time.Second)
	startErr := fxApp.Start(startCtx)
	cancelStart()
	if startErr != nil {
		log.Fatal(startErr)
	}

	runErr := runner.Run(os.Args)

	stopCtx, cancelStop := context.WithTimeout(context.Background(), 15*time.Second)
	stopErr := fxApp.Stop(stopCtx)
	cancelStop()

	if runErr != nil {
		log.Fatal(runErr)
	}
	if stopErr != nil {
		log.Fatal(stopErr)
	}
}
