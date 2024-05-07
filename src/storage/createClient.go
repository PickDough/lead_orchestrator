package storage

import (
	"leadOrchestrator/src/domain"
)

func (s *Storage) SaveClient(client *domain.Client) (*domain.Client, error) {
	query :=
		`INSERT INTO clients (name, working_hours_start, working_hours_end, lead_capacity, priority) VALUES (?, ?, ?, ?, ?)`
	res, err := s.db.Exec(
		query,
		client.Name,
		client.WorkingHoursStart.String(),
		client.WorkingHoursEnd.String(),
		client.LeadCapacity,
		client.Priority,
	)
	if err != nil {
		return nil, err
	}

	client.Id, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return client, nil
}
