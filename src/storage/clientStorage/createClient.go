package clientStorage

import (
	"errors"
	"leadOrchestrator/src/domain"
)

func (s *Storage) CreateClient(client *domain.Client) (*domain.Client, error) {
	query :=
		`INSERT INTO clients (name, working_hours_start, working_hours_end, lead_capacity, priority) 
			VALUES (?, ?, ?, ?, ?)`
	res, err := s.db.Exec(
		query,
		client.Name,
		client.WorkingHoursStart.String(),
		client.WorkingHoursEnd.String(),
		client.LeadCapacity,
		client.Priority,
	)
	// TODO: add logging
	if err != nil {
		return nil, errors.New("failed to save client")
	}

	client.Id, err = res.LastInsertId()
	// TODO: add logging
	if err != nil {
		return nil, errors.New("failed to save client")
	}

	return client, nil
}
