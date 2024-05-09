package getClients

import (
	"errors"
	"leadOrchestrator/src/domain"
	domainErrors "leadOrchestrator/src/domain/errors"
	"leadOrchestrator/src/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type getClients struct {
	provider ClientProvider
}

func NewGetClientsHandler(service ClientProvider) *getClients {
	return &getClients{provider: service}
}

type ClientProvider interface {
	GetClients(lastId int64) ([]*domain.Client, error)
}

func (gc *getClients) GetClients(c *fiber.Ctx) error {
	lastIdParam := c.Query("lastId")
	lastId, err := strconv.ParseInt(lastIdParam, 10, 64)
	if err != nil {
		if lastIdParam == "" {
			lastId = 0
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "lastId must be an integer",
			})
		}
	}
	clients, err := gc.provider.GetClients(lastId)
	if err != nil {
		var domainError *domainErrors.DomainError
		if errors.As(err, &domainError) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
		}
	}

	displayClients := []*model.DisplayClient{}
	for _, client := range clients {
		displayClient := model.MapToDisplayClient(client)
		displayClients = append(displayClients, displayClient)
	}

	return c.Status(fiber.StatusOK).JSON(displayClients)
}
