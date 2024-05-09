package leadStorage

import "errors"

func (s *Storage) CreateLead(clientId int64, state string) error {
	var err error
	if clientId == 0 {

		query := `INSERT INTO leads (state) VALUES (?)`
		_, err = s.db.Exec(query, state)
	} else {
		query := `INSERT INTO leads (client_id, state) VALUES (?, ?)`
		_, err = s.db.Exec(query, clientId, state)
	}
	if err != nil {
		return errors.New("failed to save lead")
	}

	return nil
}
