# Language Learning Portal - Backend

This is the backend service for the Language Learning Portal, built with Go and SQLite.

## Project Structure

```
backend_go/
├── cmd/
│   └── server/        # Main application entry point
├── internal/
│   ├── models/        # Database models and business logic
│   └── handlers/      # HTTP request handlers
├── db/
│   ├── migrations/    # Database schema migrations
│   └── seeds/         # Initial data for the database
└── docs/             # Documentation
```

## Prerequisites

- Go 1.21 or higher
- SQLite 3
- [Mage](https://magefile.org/) build tool

## Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd backend_go
```

2. Install dependencies:
```bash
go mod download
```

3. Initialize the database:
```bash
mage initDB
mage migrate
mage seed
```

4. Start the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`.

## Development

### Database Migrations

Migrations are handled using Mage tasks. Migration files are located in `db/migrations/`.

- Run migrations: `mage migrate`
- Reset study history: `mage resetHistory`
- Full system reset: `mage fullReset`

### Adding New Features

1. Add new models in `internal/models/`
2. Create handlers in `internal/handlers/`
3. Register routes in `cmd/server/main.go`
4. Update documentation in `docs/`

## Testing

Run tests with:
```bash
go test ./...
```

## Configuration

The application uses environment variables for configuration:
- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database path (default: words.db) 