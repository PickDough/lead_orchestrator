package assignLeadService

import (
	"leadOrchestrator/src/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGetClientsForLeadService struct {
	mock.Mock
}

func (m *MockGetClientsForLeadService) GetClientsForLead() ([]*domain.Client, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Client), args.Error(1)
}

type MockLeadsStorage struct {
	mock.Mock
}

func (m *MockLeadsStorage) CreateLead(clientId int64, state string) error {
	args := m.Called(clientId, state)
	return args.Error(0)
}

func TestAssignLead_NoClients(t *testing.T) {
	mockGetClients := new(MockGetClientsForLeadService)
	mockLeadsStorage := new(MockLeadsStorage)

	mockGetClients.On("GetClientsForLead").Return([]*domain.Client{}, nil)
	mockLeadsStorage.On("CreateLead", int64(0), "assigned").Return(nil)

	service := NewAssignLeadByPriorityAndMaxCapacityService(mockGetClients, mockLeadsStorage, time.Now)
	client, err := service.AssignLead()

	mockGetClients.AssertExpectations(t)
	mockLeadsStorage.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Nil(t, client)
}

func TestAssignLead_WithClients(t *testing.T) {
	mockGetClients := new(MockGetClientsForLeadService)
	mockLeadsStorage := new(MockLeadsStorage)

	clients := []*domain.Client{
		{
			Id:                1,
			Name:              "Before now",
			WorkingHoursStart: time.Duration(9) * time.Hour,
			WorkingHoursEnd:   time.Duration(11) * time.Hour,
			LeadCapacity:      10,
			Priority:          1,
		},
		{
			Id:                2,
			Name:              "After now",
			WorkingHoursStart: time.Duration(13) * time.Hour,
			WorkingHoursEnd:   time.Duration(19) * time.Hour,
			LeadCapacity:      10,
			Priority:          1,
		},
		{
			Id:                3,
			Name:              "Low priority",
			WorkingHoursStart: time.Duration(10) * time.Hour,
			WorkingHoursEnd:   time.Duration(12) * time.Hour,
			LeadCapacity:      5,
			Priority:          3,
		},
		{
			Id:                4,
			Name:              "Low capacity",
			WorkingHoursStart: time.Duration(10) * time.Hour,
			WorkingHoursEnd:   time.Duration(12) * time.Hour,
			LeadCapacity:      2,
			Priority:          1,
		},
		{
			Id:                5,
			Name:              "Best client",
			WorkingHoursStart: time.Duration(10) * time.Hour,
			WorkingHoursEnd:   time.Duration(12) * time.Hour,
			LeadCapacity:      5,
			Priority:          1,
		},
	}

	mockGetClients.On("GetClientsForLead").Return(clients, nil)
	mockLeadsStorage.On("CreateLead", clients[4].Id, "assigned").Return(nil)

	service := &assignLeadByPriorityAndMaxCapacityService{
		getClientsForLeadService: mockGetClients,
		leadsStorage:             mockLeadsStorage,
		now: func() time.Time {
			return time.Date(1970, 1, 1, 12, 0, 0, 0, time.UTC)
		},
	}
	client, err := service.AssignLead()

	mockGetClients.AssertExpectations(t)
	mockLeadsStorage.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, clients[4], client)
}
