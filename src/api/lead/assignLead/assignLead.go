package assignLead

import (
	"errors"
	domainErrors "leadOrchestrator/src/domain/errors"
	"leadOrchestrator/src/model"
	"leadOrchestrator/src/service/assignLeadService"

	"github.com/gofiber/fiber/v2"
)

type asignLeadhandler struct {
	assignLeadService assignLeadService.AssignLeadService
}

func NewAssignLeadHandler(assignLeadService assignLeadService.AssignLeadService) *asignLeadhandler {
	return &asignLeadhandler{assignLeadService: assignLeadService}
}

func (alh *asignLeadhandler) AssignLead(c *fiber.Ctx) error {
	client, err := alh.assignLeadService.AssignLead()
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
			"error": "suitable client not found",
		})
	}

	displayClient := model.MapToDisplayClient(client)

	return c.Status(fiber.StatusOK).JSON(displayClient)
}
