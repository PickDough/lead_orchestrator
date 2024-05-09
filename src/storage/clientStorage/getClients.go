package clientStorage

import (
	"database/sql"
	"errors"
	"leadOrchestrator/src/domain"
	"time"
)

func (s *Storage) GetClients(lastId int64) ([]*domain.Client, error) {
	var rows *sql.Rows
	if lastId > 0 {
		query := `SELECT * FROM clients WHERE id > ? ORDER BY id ASC LIMIT 100`
		var err error
		rows, err = s.db.Query(query, lastId)
		if err != nil {
			return nil, errors.New("failed to retrieve clients")
		}
	} else {
		query := `SELECT * FROM clients ORDER BY id ASC LIMIT 100`
		var err error
		rows, err = s.db.Query(query)
		if err != nil {
			return nil, errors.New("failed to retrieve clients")
		}
	}
	rows.Scan()

	clients := make([]*domain.Client, 0)

	for rows.Next() {
		client, err := mapClients(rows)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func mapClients(row *sql.Rows) (*domain.Client, error) {
	client := &domain.Client{}
	start := ""
	end := ""
	err := row.Scan(&client.Id, &client.Name, &start, &end, &client.LeadCapacity, &client.Priority)
	if err != nil {
		return nil, err
	}
	client.WorkingHoursStart, _ = time.ParseDuration(start)
	client.WorkingHoursEnd, _ = time.ParseDuration(end)

	return client, nil
}
