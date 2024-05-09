package clientStorage

import (
	"errors"
	"leadOrchestrator/src/domain"
)

func (s *Storage) GetClientsForLead() ([]*domain.Client, error) {
	query := `SELECT clients.* FROM clients
	LEFT JOIN leads ON clients.id = leads.client_id
	WHERE leads.client_id is null or leads.state = 'assigned'
	GROUP BY clients.id
	HAVING COUNT(leads.id) < clients.lead_capacity`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, errors.New("failed to retrieve clients")
	}

	clients := []*domain.Client{}
	for rows.Next() {
		client, err := mapClients(rows)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}
