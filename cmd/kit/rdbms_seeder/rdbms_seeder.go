package main

import (
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"

	appregistrar "github.com/chan-jui-huang/go-backend-framework/v2/cmd/app/registrar"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/migration/rdbms/seeder"
	booter "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
)

func main() {
	fxApp := fx.New(
		fx.Supply(booter.NewDefaultConfig()),
		fx.Provide(
			appregistrar.NewConfigLoader,
			appregistrar.NewDatabaseConfig,
			appregistrar.NewDatabase,
		),
		fx.Invoke(
			appregistrar.RegisterConfigDependencies,
			appregistrar.RegisterServiceDependencies,
		),
	)
	if err := fxApp.Err(); err != nil {
		panic(err)
	}

	seederExecutor := seeder.NewSeederExecutor()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "show",
				Usage: "show all seeders",
				Action: func(cCtx *cli.Context) error {
					seederExecutor.ShowSeeders()
					return nil
				},
			},
			{
				Name:  "run",
				Usage: "Run the seeders. EX: database_seeder run seeder1,seeder2 (run specific seeders). database_seeder run (run all seeders).",
				Action: func(cCtx *cli.Context) error {
					seederExecutor.Run(strings.Split(cCtx.Args().First(), ","))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
