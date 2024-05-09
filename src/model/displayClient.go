package model

import (
	"fmt"
	"leadOrchestrator/src/domain"
	"time"
)

type DisplayClient struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	WorkingHours string `json:"workingHours"`
	LeadCapacity int    `json:"leadCapacity"`
	Priority     int    `json:"priority"`
}

func MapToDisplayClient(client *domain.Client) *DisplayClient {
	return &DisplayClient{
		Id:   client.Id,
		Name: client.Name,
		WorkingHours: fmt.Sprintf(
			"%02d:%02d-%02d:%02d",
			int(client.WorkingHoursStart.Hours()),
			int(client.WorkingHoursStart.Minutes()-
				client.WorkingHoursStart.Truncate(time.Hour).Hours()*60),
			int(client.WorkingHoursEnd.Hours()),
			int(client.WorkingHoursEnd.Minutes()-
				client.WorkingHoursEnd.Truncate(time.Hour).Hours()*60),
		),
		LeadCapacity: client.LeadCapacity,
		Priority:     client.Priority,
	}
}
