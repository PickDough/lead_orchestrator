package createClientService

import (
	"fmt"
	"leadOrchestrator/src/domain"
	"leadOrchestrator/src/domain/errors"
	"leadOrchestrator/src/model"
	"strings"
	"time"
)

type saveClient interface {
	CreateClient(client *domain.Client) (*domain.Client, error)
}

type createClientService struct {
	saveClient saveClient
}

func NewCreateClientService(saveClient saveClient) *createClientService {
	return &createClientService{saveClient: saveClient}
}

func (ccs *createClientService) Create(model *model.CreateClientModel) (*domain.Client, error) {
	hours := strings.Split(model.WorkingHours, "-")
	h1 := strings.Split(hours[0], ":")
	d1, err := time.ParseDuration(fmt.Sprintf("%sh%sm", h1[0], h1[1]))
	if err != nil {
		return nil, err
	}
	h2 := strings.Split(hours[1], ":")
	d2, err := time.ParseDuration(fmt.Sprintf("%sh%sm", h2[0], h2[1]))
	if err != nil {
		return nil, err
	}

	if d1 >= d2 {
		return nil, &errors.WorkingHoursEndBeforeStartError{}
	}

	client := &domain.Client{
		Name:              model.Name,
		WorkingHoursStart: d1,
		WorkingHoursEnd:   d2,
		LeadCapacity:      model.LeadCapacity,
		Priority:          model.Priority,
	}

	client, err = ccs.saveClient.CreateClient(client)
	if err != nil {
		return nil, err
	}

	return client, nil
}
