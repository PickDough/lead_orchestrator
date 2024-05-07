-- Create the clients table
CREATE TABLE clients (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    working_hours_start TEXT NOT NULL,
    working_hours_end TEXT NOT NULL,
    lead_capacity INTEGER NOT NULL,
    priority INTEGER DEFAULT 0
);