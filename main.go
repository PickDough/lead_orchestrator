package main

import (
	"flag"
	"fmt"
	"leadOrchestrator/src/app"
	"time"

	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type config struct {
	DbConnectionString string `env:"DB_CONNECTION_STRING"`
	Port               string `env:"PORT" envDefault:"3000"`
}

type App interface {
}

var strategy = flag.String("ls", "ByPriorityAndMaxCapacity", "the strategy to use to assign leads")

func main() {
	flag.Parse()

	cfg := readConfig()

	app := app.NewApp(app.Config{
		DbConnectionString: cfg.DbConnectionString,
		MigrationsPath:     "migrations",
		Port:               cfg.Port,
		Strategy:           *strategy,
		Now:                time.Now,
	})

	app.Listen(fmt.Sprintf(":%s", cfg.Port))
}

func readConfig() config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	cfg, err := env.ParseAs[config]()
	if err != nil {
		log.Fatalln(err)
	}
	return cfg
}
