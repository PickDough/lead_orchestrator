package clientStorage

import (
	"database/sql"
	"errors"
	"leadOrchestrator/src/domain"
	"time"
)

func (s *Storage) GetClient(id int64) (*domain.Client, error) {
	query := `SELECT * FROM clients WHERE id = ? LIMIT 1`

	row := s.db.QueryRow(query, id)
	client, err := mapClient(row)
	// TODO: add logging
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, errors.New("failed to retrieve client")
	}

	return client, nil
}

func mapClient(row *sql.Row) (*domain.Client, error) {
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
