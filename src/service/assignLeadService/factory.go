package assignLeadService

import (
	"leadOrchestrator/src/domain"
	"time"
)

type LeadStrategyFactory struct {
	getClientsForLeadService GetClientsForLeadService
	leadsStorage             LeadsStorage
}

type AssignLeadService interface {
	AssignLead() (*domain.Client, error)
}

func NewServiceFactory(
	getClientsForLeadService GetClientsForLeadService,
	leadsStorage LeadsStorage,
) *LeadStrategyFactory {
	return &LeadStrategyFactory{
		getClientsForLeadService: getClientsForLeadService,
		leadsStorage:             leadsStorage,
	}
}

func (sf *LeadStrategyFactory) CreateStrategy(serviceType string, now func() time.Time) AssignLeadService {
	switch serviceType {
	case "ByPriorityAndMaxCapacity":
		return NewAssignLeadByPriorityAndMaxCapacityService(sf.getClientsForLeadService, sf.leadsStorage, now)
	default:
		return nil
	}
}
