package api

import (
	"leadOrchestrator/src/domain"
	"leadOrchestrator/src/model"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()
var r = regexp.MustCompile("(([01][0-9])|(2[1-3])):([0-5][0-9])-(([01][0-9])|(2[1-3])):([0-5][0-9])")

type createClient struct {
	service CreateClientService
}

func NewCreateClientHandler(service CreateClientService) *createClient {
	return &createClient{service: service}
}

type CreateClientService interface {
	Create(model *model.CreateClientModel) (*domain.Client, error)
}

func (cc *createClient) Create(c *fiber.Ctx) error {
	req := &model.CreateClientModel{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if ok := r.Match([]byte(req.WorkingHours)); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "working hours must be in the format HH:MM-HH:MM",
		})
	}

	client, err := cc.service.Create(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	displayClient := model.MapToDisplayClient(client)

	return c.Status(fiber.StatusCreated).JSON(displayClient)
}
