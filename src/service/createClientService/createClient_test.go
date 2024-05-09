package createClientService

import (
	"leadOrchestrator/src/domain"
	"leadOrchestrator/src/domain/errors"
	"leadOrchestrator/src/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type saveClientMock struct {
	mock.Mock
}

func (m *saveClientMock) CreateClient(client *domain.Client) (*domain.Client, error) {
	args := m.Called(client)
	return args.Get(0).(*domain.Client), args.Error(1)
}

func TestCreateClientSuccess(t *testing.T) {
	testCases := []struct {
		name           string
		clientModel    *model.CreateClientModel
		expectedClient *domain.Client
	}{
		{
			name: "Simple time",
			clientModel: &model.CreateClientModel{
				Name:         "Test Client",
				WorkingHours: "02:00-05:00",
				LeadCapacity: 10,
				Priority:     1,
			},
			expectedClient: &domain.Client{
				Name:              "Test Client",
				WorkingHoursStart: 2 * time.Hour,
				WorkingHoursEnd:   5 * time.Hour,
				LeadCapacity:      10,
				Priority:          1,
			},
		},
		{
			name: "Time with minutes",
			clientModel: &model.CreateClientModel{
				Name:         "Test Client",
				WorkingHours: "02:34-05:57",
				LeadCapacity: 10,
				Priority:     1,
			},
			expectedClient: &domain.Client{
				Name:              "Test Client",
				WorkingHoursStart: 2*time.Hour + 34*time.Minute,
				WorkingHoursEnd:   5*time.Hour + 57*time.Minute,
				LeadCapacity:      10,
				Priority:          1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create an instance of our test object
			mockObj := new(saveClientMock)

			// Setup expectations
			mockObj.On("SaveClient", tc.expectedClient).Return(tc.expectedClient, nil)

			// Create service
			service := NewCreateClientService(mockObj)

			// Call the method we want to test
			result, err := service.Create(tc.clientModel)

			// Assert expectations
			mockObj.AssertExpectations(t)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedClient, result)
		})
	}
}

func TestCreateClientStartBeforeEnd(t *testing.T) {
	testCases := []struct {
		name        string
		clientModel *model.CreateClientModel
	}{
		{
			name: "Simple time",
			clientModel: &model.CreateClientModel{
				Name:         "Test Client",
				WorkingHours: "02:00-01:00",
				LeadCapacity: 10,
				Priority:     1,
			},
		},
		{
			name: "Edge time",
			clientModel: &model.CreateClientModel{
				Name:         "Test Client",
				WorkingHours: "23:59-00:00",
				LeadCapacity: 10,
				Priority:     1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create an instance of our test object
			mockObj := new(saveClientMock)

			// Setup expectations
			mockObj.AssertNotCalled(t, "SaveClient")

			// Create service
			service := NewCreateClientService(mockObj)

			// Call the method we want to test
			result, err := service.Create(tc.clientModel)

			// Assert expectations
			mockObj.AssertExpectations(t)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, &errors.WorkingHoursEndBeforeStartError{})
		})
	}
}
