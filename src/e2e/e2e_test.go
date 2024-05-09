package e2e

import (
	"bytes"
	"encoding/json"
	"leadOrchestrator/src/app"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setup() *fiber.App {
	cfg := app.Config{
		DbConnectionString: ":memory:",
		MigrationsPath:     "./migrations",
		Port:               "3000",
		Strategy:           "ByPriorityAndMaxCapacity",
		Now: func() time.Time {
			return time.Date(1970, 1, 1, 12, 0, 0, 0, time.UTC)
		},
	}

	app := app.NewApp(cfg)

	return app
}

func TestCreateClient(t *testing.T) {
	app := setup()
	testCases := []struct {
		name               string
		clientData         map[string]any
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "Success",
			clientData: map[string]any{
				"name":         "Test Client",
				"workingHours": "02:00-05:00",
				"leadCapacity": 10,
				"priority":     1,
			},
			expectedStatusCode: 200,
			expectedBody: `{
				"id":3,
				"name":"Test Client",
				"workingHours":"02:00-05:00",
				"leadCapacity":10,
				"priority":1
			}`,
		},
		{
			name:               "Required fields missing",
			clientData:         map[string]any{},
			expectedStatusCode: 400,
			expectedBody:       `{"error":"Key: 'CreateClientModel.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'CreateClientModel.WorkingHours' Error:Field validation for 'WorkingHours' failed on the 'required' tag\nKey: 'CreateClientModel.LeadCapacity' Error:Field validation for 'LeadCapacity' failed on the 'required' tag"}`,
		},
		{
			name: "Incorrect working hours format",
			clientData: map[string]any{
				"name":         "Test Client",
				"workingHours": "25:00-03:00",
				"leadCapacity": 10,
				"priority":     1,
			},
			expectedStatusCode: 400,
			expectedBody:       `{"error": "working hours must be in the format HH:MM-HH:MM"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.clientData)
			req, _ := http.NewRequest("POST", "/clients", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, 10)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			respBody := make([]byte, resp.ContentLength)
			_, _ = resp.Body.Read(respBody)
			assert.JSONEq(t, tc.expectedBody, string(respBody))
		})
	}
}

func TestGetClient(t *testing.T) {
	app := setup()

	testCases := []struct {
		name               string
		clientId           string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Get Client with ID 1",
			clientId:           "1",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"id":1,"name":"Client 1","workingHours":"09:00-13:00","leadCapacity":5,"priority":1}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/clients/"+tc.clientId, nil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, 10)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			respBody := make([]byte, resp.ContentLength)
			_, _ = resp.Body.Read(respBody)
			assert.JSONEq(t, tc.expectedBody, string(respBody))
		})
	}
}

func TestGetClients(t *testing.T) {
	app := setup()

	testCases := []struct {
		name               string
		lastId             string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Get Client with lastId 0",
			lastId:             "0",
			expectedStatusCode: http.StatusOK,
			expectedBody: `[
				{"id":1,"name":"Client 1","workingHours":"09:00-13:00","leadCapacity":5,"priority":1},
				{"id":2,"name":"Client 2","workingHours":"13:00-19:00","leadCapacity":10,"priority":1}
				]`,
		},
		{
			name:               "Get Client with lastId 1",
			lastId:             "1",
			expectedStatusCode: http.StatusOK,
			expectedBody: `[
				{"id":2,"name":"Client 2","workingHours":"13:00-19:00","leadCapacity":10,"priority":1}
				]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/clients?lastId="+tc.lastId, nil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, 10)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			respBody := make([]byte, resp.ContentLength)
			_, _ = resp.Body.Read(respBody)
			assert.JSONEq(t, tc.expectedBody, string(respBody))
		})
	}
}

func TestAssignLead(t *testing.T) {
	app := setup()

	testCases := []struct {
		name               string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Assign lead to client 1",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"id":1,"name":"Client 1","workingHours":"09:00-13:00","leadCapacity":5,"priority":1}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/leads/assign", nil)
			resp, err := app.Test(req, 10)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			respBody := make([]byte, resp.ContentLength)
			_, _ = resp.Body.Read(respBody)
			assert.JSONEq(t, tc.expectedBody, string(respBody))
		})
	}
}
