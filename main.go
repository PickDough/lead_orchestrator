package main

import (
	"fmt"
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

func main() {
	cfg := readConfig()

	_ = connectDb(cfg)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

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
