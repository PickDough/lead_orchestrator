package app

import (
	"fmt"
	"leadOrchestrator/src/api/client/createClient"
	"leadOrchestrator/src/api/client/getClient"
	"leadOrchestrator/src/api/client/getClients"
	"leadOrchestrator/src/api/lead/assignLead"
	"leadOrchestrator/src/service/assignLeadService"
	"leadOrchestrator/src/service/createClientService"
	"leadOrchestrator/src/storage/clientStorage"
	"leadOrchestrator/src/storage/leadStorage"
	"time"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	DbConnectionString string
	MigrationsPath     string
	Port               string
	Strategy           string
	Now                func() time.Time
}

func NewApp(cfg Config) *fiber.App {
	db := connectDb(cfg)

	app := fiber.New()

	storage := clientStorage.NewStorage(db)

	createClientService := createClientService.NewCreateClientService(storage)
	createClientHandler := createClient.NewCreateClientHandler(createClientService)
	app.Post("/clients", createClientHandler.Create)

	getClientsHandler := getClients.NewGetClientsHandler(storage)
	app.Get("/clients", getClientsHandler.GetClients)

	getClientHandler := getClient.NewGetClientHandler(storage)
	app.Get("/clients/:id", getClientHandler.GetClient)

	factory := assignLeadService.NewServiceFactory(storage, leadStorage.NewStorage(db))
	serviceStrategy := factory.CreateStrategy(cfg.Strategy, cfg.Now)
	if serviceStrategy == nil {
		log.Fatalf("Strategy %s not found", cfg.Strategy)
	}
	assignLeadHandler := assignLead.NewAssignLeadHandler(serviceStrategy)
	app.Post("/leads/assign", assignLeadHandler.AssignLead)

	return app
}

func connectDb(cfg Config) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", cfg.DbConnectionString)
	if err != nil {
		log.Fatalln(err)
	}

	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		"leads", driver)
	if err != nil {
		log.Fatalln(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalln(err)
	}

	return db
}
