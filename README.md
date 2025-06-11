# Go Boilerplate

A clean and well-structured Go boilerplate following Clean Architecture principles, inspired by industry best practices.

## 🚀 Features

- 🏗️ **Clean Architecture** - Organized in layers (Domain, Application, Infrastructure, Presentation)
- 🚀 **Gin Framework** - Fast HTTP web framework
- 🗄️ **GORM** - Feature-rich ORM with support for PostgreSQL and SQLite
- ⚙️ **Viper Configuration** - Flexible configuration management
- 📝 **Structured Logging** - Using zerolog with structured output
- 🐳 **Docker Support** - Ready-to-use Docker setup
- 🔄 **Hot Reload** - Development setup with Air
- 🧪 **Testing Ready** - Structured for easy testing
- 📊 **Health Checks** - Built-in health check endpoints

## Project Structure

```
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── database/          # Database connection and migrations
│   ├── middleware/        # HTTP middleware
│   ├── server/            # HTTP server setup
│   └── [modules]/         # Feature modules (domain-driven)
│       ├── domain/        # Business entities and interfaces
│       ├── application/   # Use cases and business logic
│       ├── infrastructure/# External concerns (repositories, etc.)
│       └── presentation/  # HTTP handlers and DTOs
├── pkg/                   # Public library code
├── tests/                 # Test files
├── data/                  # Database files (SQLite)
├── config.yaml           # Configuration file
├── docker-compose.yml    # Docker Compose setup
└── Makefile              # Build and development commands
```

## Quick Start

### Prerequisites

- Go 1.22 or higher
- Make (optional, for convenience commands)
- Docker & Docker Compose (optional)

### Local Development

1. **Clone and setup:**
   ```bash
   git clone <your-repo>
   cd boilerplate-go
   cp config.example.yaml config.yaml
   cp .env.example .env
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run the application:**
   ```bash
   make run
   # or
   go run ./cmd
   ```

4. **Access the API:**
   - Health check: http://localhost:8080/health
   - Welcome endpoint: http://localhost:8080/api/v1/

### Using Docker

1. **Start with Docker Compose:**
   ```bash
   make docker-up
   ```

2. **Stop services:**
   ```bash
   make docker-down
   ```

### Development with Hot Reload

1. **Install Air:**
   ```bash
   make install-dev-tools
   ```

2. **Start development server:**
   ```bash
   make dev
   ```

## Configuration

Configuration can be managed through:
- `config.yaml` - Main configuration file
- Environment variables (prefixed with `APP_`)
- Command line flags (to be implemented)

### Configuration Hierarchy
1. Environment variables (highest priority)
2. Configuration file
3. Default values (lowest priority)

### Database Configuration

**SQLite (default):**
```yaml
database:
  driver: "sqlite"
  sqlite:
    path: "./data/app.db"
```

**PostgreSQL:**
```yaml
database:
  driver: "postgres"
  postgres:
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "password"
    database: "boilerplate"
    sslmode: "disable"
```

## Available Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run the application
make test              # Run tests
make test-coverage     # Run tests with coverage
make dev               # Run with hot reload
make docker-up         # Start with Docker Compose
make docker-down       # Stop Docker services
make clean             # Clean build artifacts
make fmt               # Format code
make lint              # Run linter
```

## Adding New Features

When adding new features, follow the Clean Architecture pattern:

1. **Create a new module directory:**
   ```
   internal/users/
   ├── domain/
   │   ├── user.go              # Entity
   │   └── user_repository.go   # Repository interface
   ├── application/
   │   └── user_service.go      # Business logic
   ├── infrastructure/
   │   └── gorm_user_repository.go  # Repository implementation
   └── presentation/
       └── user_controller.go   # HTTP handlers
   ```

2. **Register routes in server/server.go**
3. **Add migrations in database/migrations.go**
4. **Update configuration if needed**

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/config
```

## Production Deployment

1. **Build Docker image:**
   ```bash
   make docker-build
   ```

2. **Set production configuration:**
   - Use environment variables for sensitive data
   - Set `APP_SERVER_MODE=release`
   - Configure proper database connection
   - Set appropriate log levels

3. **Deploy with your preferred method**

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes following the established patterns
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
