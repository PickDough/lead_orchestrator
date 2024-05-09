package getClient

import (
	"errors"
	"leadOrchestrator/src/domain"
	domainErrors "leadOrchestrator/src/domain/errors"
	"leadOrchestrator/src/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type getClient struct {
	provider ClientProvider
}

func NewGetClientHandler(service ClientProvider) *getClient {
	return &getClient{provider: service}
}

type ClientProvider interface {
	GetClient(id int64) (*domain.Client, error)
}

func (gc *getClient) GetClient(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id must be an integer",
		})
	}
	client, err := gc.provider.GetClient(id)
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

	if client == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "client not found",
		})
	}

	displayClient := model.MapToDisplayClient(client)

	return c.Status(fiber.StatusOK).JSON(displayClient)
}
