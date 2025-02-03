# Pagination App

## Project Structure
This project follows a structured approach to maintain scalability and separation of concerns.

```
├── bin/                 # Contains built binary files
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Dockerfile for building the application
├── env-example          # Example environment variable file
├── go.mod               # Go module file
├── go.sum               # Go dependency lock file
├── internal/            # Contains all business logic (Core of the application)
│   ├── api/             # API handlers and request handling logic
│   │   ├── api.go
│   │   └── test.go
│   ├── domain/          # Business entities, domain models, and business rules
├── main.go              # Application entry point
└── pkg/                 # Third-party dependencies and configuration
    └── conf.go          # Application configuration logic
```

### Internal Package
- The **`internal/`** directory contains the business logic and core application functionality.
- **`internal/api/`**: Handles API routes and request processing.
- **`internal/domain/`**: Contains domain models and business rules.

### Pkg Package
- The **`pkg/`** directory is used for third-party integrations and configuration management.
- **`pkg/conf.go`**: Handles application configuration.

### Bin Directory
- The **`bin/`** directory stores all compiled binaries after building the application.

## How to Build and Run the Application

### 1. Configure Environment Variables
Copy the `.env-example` file and rename it to `.env`. Update it with your required environment variables.

```sh
cp env-example .env
```

### 2. Build and Run Using Docker Compose
Run the following command to build and start the application:

```sh
docker compose up --build
```

This will:
- Build the Go application using a multi-stage Docker build.
- Start the containerized application with the necessary environment configurations.

Once the application is running, it will be accessible via the exposed `SERVER_PORT` defined in your `.env` file.

## Notes
- Ensure that Docker and Docker Compose are installed on your system before running the commands.
- The application uses PostgreSQL as a database (if configured in `.env`). Ensure that your database is running and accessible.

---

**Happy Coding! 🚀**
