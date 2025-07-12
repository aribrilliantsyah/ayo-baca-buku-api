# Ayo Baca Buku API

API service for the Ayo Baca Buku application, managing users, books, and reading activities.

## Features

*   User Registration and Login
*   CRUD operations for Users
*   (Placeholder for future features like Book Management, Reading Activity Tracking)

## Technology Stack

*   **Go (Golang)**
*   **Fiber** (Web Framework)
*   **GORM** (ORM for database interaction)
*   **PostgreSQL** (Assumed database, please verify and update if different)
*   **Swagger** (API Documentation)
*   **Zap** (Logging)

## Setup & Installation

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd <repository-name>
    ```

2.  **Install Go:**
    Ensure you have Go installed (version 1.18 or higher is recommended).

3.  **Environment Configuration:**
    Copy the example environment file and update it with your local configuration:
    ```bash
    cp example.env .env
    ```
    Edit `.env` to set your database credentials, JWT secret, and other necessary configurations.

4.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

## Database

This project uses GORM for database interactions. Ensure your database server (e.g., PostgreSQL) is running and configured in your `.env` file.

*   **Migrations:** Database migrations are applied automatically when the application starts. These will set up the necessary tables.
*   **Seeders:** Basic seed data (e.g., an initial admin user) is also inserted automatically on startup if the database is empty or specific tables are empty.

## Running the Application

To run the application:

```bash
go run cmd/main.go
```

The server will start, typically on `http://localhost:3000` (or as configured by the `PORT` environment variable if set).

## API Endpoints & Documentation

The API provides various endpoints for managing application resources. Comprehensive API documentation is available via Swagger:

*   **Swagger UI:** [http://localhost:3000/docs/](http://localhost:3000/docs/)
*   **Scalar UI:** [http://localhost:3000/scalar](http://localhost:3000/scalar) (Alternative API documentation interface)

## Directory Structure

```
.
├── app/                  # Core application logic
│   ├── config/           # Application configuration (e.g., loading .env)
│   ├── controllers/      # HTTP request handlers (business logic)
│   ├── database/         # Database connection, migrations, seeders
│   ├── models/           # GORM models and request/response structs
│   ├── routes/           # API route definitions
│   └── util/             # Utility packages (JWT, logger, validation, etc.)
├── cmd/                  # Main application entry point (main.go)
├── docs/                 # Swagger API documentation files (generated)
├── logs/                 # Application log files
├── .env                  # Local environment configuration (ignored by Git)
├── example.env           # Example environment file
├── go.mod                # Go module definitions
├── go.sum                # Go module checksums
└── README.md             # This file
```

## Environment Variables

Key environment variables to configure in your `.env` file (refer to `example.env` for a full list):

*   `DB_HOST`: Database host
*   `DB_PORT`: Database port
*   `DB_USER`: Database username
*   `DB_PASSWORD`: Database password
*   `DB_NAME`: Database name
*   `DB_SSL_MODE`: (e.g., `disable`, `require`)
*   `JWT_SECRET_KEY`: Secret key for signing JWT tokens
*   `LOG_LEVEL`: Logging level (e.g., `debug`, `info`, `warn`, `error`)
*   `PORT`: (Optional) Port for the application to listen on (defaults to 3000 if not set)

## Contributing

Contributions are welcome! Please follow the coding standards and guidelines. For AI-assisted development, please refer to the `AGENTS.md` file in the repository for specific instructions and conventions to maintain consistency.
---

*This README provides a general overview. Please update it with more specific details as the project evolves.*
