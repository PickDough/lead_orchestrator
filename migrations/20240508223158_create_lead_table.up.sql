CREATE TABLE leads (
    id INTEGER PRIMARY KEY,
    client_id INTEGER,
    state TEXT CHECK(state IN ('assigned', 'completed', 'not_assigned')) NOT NULL,
    time_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(client_id) REFERENCES client(id)
);