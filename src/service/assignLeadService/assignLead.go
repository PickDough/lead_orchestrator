package assignLeadService

import (
	"leadOrchestrator/src/domain"
	"time"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type assignLeadByPriorityAndMaxCapacityService struct {
	getClientsForLeadService GetClientsForLeadService
	leadsStorage             LeadsStorage
	now                      func() time.Time
}

type GetClientsForLeadService interface {
	GetClientsForLead() ([]*domain.Client, error)
}

type LeadsStorage interface {
	CreateLead(clientId int64, state string) error
}

func NewAssignLeadByPriorityAndMaxCapacityService(
	getClientsForLeadService GetClientsForLeadService,
	leadsStorage LeadsStorage,
	now func() time.Time,
) *assignLeadByPriorityAndMaxCapacityService {
	return &assignLeadByPriorityAndMaxCapacityService{
		getClientsForLeadService: getClientsForLeadService,
		leadsStorage:             leadsStorage,
		now:                      now,
	}
}

func (al *assignLeadByPriorityAndMaxCapacityService) AssignLead() (*domain.Client, error) {
	clients, err := al.getClientsForLeadService.GetClientsForLead()
	if err != nil {
		return nil, err
	}

	// would move time logic to SQL if I had to rewrite this repository
	now := al.now()
	nowDuration := time.Duration(now.Hour())*time.Hour + time.Duration(now.Minute())*time.Minute
	clientsWithinTimeRange :=
		lo.Filter(clients, func(client *domain.Client, idx int) bool {
			return client.WorkingHoursStart <= nowDuration && nowDuration <= client.WorkingHoursEnd
		})
	slices.SortFunc(clientsWithinTimeRange, func(a, b *domain.Client) int {
		if a.Priority == b.Priority {
			return b.LeadCapacity - a.LeadCapacity
		}
		return a.Priority - b.Priority
	})
	if len(clients) == 0 {
		al.leadsStorage.CreateLead(0, "assigned")
		return nil, nil
	}

	err = al.leadsStorage.CreateLead(clientsWithinTimeRange[0].Id, "assigned")
	if err != nil {
		return nil, err
	}

	return clientsWithinTimeRange[0], nil
}
