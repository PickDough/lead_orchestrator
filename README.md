# Lead Orchestrator

Lead Orchestrator is a Go application that provides APIs for managing clients and leads.

## Setup

To run the application, you need to provide a configuration with the following fields:

- `DbConnectionString`: The connection string for the SQLite database.
- `MigrationsPath`: The path to the database migrations.
- `Port`: The port on which the application should listen.
- `Strategy`: The strategy to use for assigning leads.
- `Now`: A function that returns the current time.

First three are read from `env`.

The strategy is passed through cli args.

## API Endpoints

| Endpoint        | Method | Description                           | Parameters                                                                    |
| --------------- | ------ | ------------------------------------- | ----------------------------------------------------------------------------- |
| `/clients`      | POST   | Create a new client.                  | Body: `name`, `workingHours`, `working_hours_end`, `leadCapacity`, `priority` |
| `/clients`      | GET    | Get a list of all clients.            | Query param `?lastId`                                                         |
| `/clients/:id`  | GET    | Get the details of a specific client. | `id`                                                                          |
| `/leads/assign` | POST   | Assign a lead to a client.            | None                                                                          |

## Database

The application uses a SQLite database with two tables: `clients` and `leads`. The `clients` table has columns for `id`, `name`, `working_hours_start`, `working_hours_end`, `lead_capacity`, and `priority`. The `leads` table has columns for `id`, `client_id`, `state`, and `time_created`.

## Lead Assignment Strategy

The application supports different strategies for assigning leads to clients, although only one is provided as implementation. The strategy is specified in the configuration. If the specified strategy is not found, the application will log an error and exit.

## Tests
Both of the services have been accompanied by unit tests to ensure they comply with business rules.

In addition, this repository contains `E2E` test suite that checks all api methods.