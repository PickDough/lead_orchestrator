package main

import (
	"fmt"
	"leadOrchestrator/src/api"
	"leadOrchestrator/src/service"
	"leadOrchestrator/src/storage"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type config struct {
	DbConnectionString string `env:"DB_CONNECTION_STRING"`
	Port               string `env:"PORT" envDefault:"3000"`
}

type App interface {
}

func main() {
	cfg := readConfig()

	db := connectDb(cfg)

	app := fiber.New()

	storage := storage.NewStorage(db)

	createClientService := service.NewCreateClientService(storage)
	createClientHandler := api.NewCreateClientHandler(createClientService)
	app.Post("/clients", createClientHandler.Create)

	app.Listen(fmt.Sprintf(":%s", cfg.Port))
}

func connectDb(cfg config) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", cfg.DbConnectionString)
	if err != nil {
		log.Fatalln(err)
	}

	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"leads", driver)
	if err != nil {
		log.Fatalln(err)
	}

	m.Up()

	return db
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
