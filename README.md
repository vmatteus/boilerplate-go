# Go Boilerplate

A clean and well-structured Go boilerplate following Clean Architecture principles, with advanced logging and dependency injection capabilities.

## 🚀 Features

- 🏗️ **Clean Architecture** - Organized in layers (Domain, Application, Infrastructure, Presentation)
- 🚀 **Gin Framework** - Fast HTTP web framework
- 🗄️ **GORM** - Feature-rich ORM with support for PostgreSQL and SQLite
- ⚙️ **Viper Configuration** - Flexible configuration management
- 📝 **Advanced Logging** - Structured logging with OpenTelemetry integration
- 🔍 **OpenTelemetry Ready** - Full observability with tracing and metrics
- 🧩 **Dependency Injection** - Clean DI with Uber FX
- 🐳 **Docker Support** - Ready-to-use Docker setup
- 🔄 **Hot Reload** - Development setup with Air
- 🧪 **Testing Ready** - Structured for easy testing with mocks
- 📊 **Health Checks** - Built-in health check endpoints
- 🔌 **Multiple Log Providers** - stdout, file, elasticsearch, logstash support

## Project Structure

```
├── cmd/                    # Application entry points
│   ├── main.go            # Application entry point
│   └── examples/          # Usage examples and demos
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── database/          # Database connection and migrations
│   ├── fx/                # Dependency injection configuration
│   ├── logger/            # Advanced logging with OpenTelemetry
│   │   ├── logger.go      # Main logger implementation
│   │   ├── stdout_logger.go    # Console logging
│   │   ├── file_logger.go      # File logging
│   │   ├── elasticsearch_logger.go # Elasticsearch integration
│   │   └── logstash_logger.go   # Logstash TCP logging
│   ├── middleware/        # HTTP middleware
│   ├── server/            # HTTP server setup
│   ├── telemetry/         # OpenTelemetry configuration
│   └── [modules]/         # Feature modules (domain-driven)
│       ├── domain/        # Business entities and interfaces
│       ├── application/   # Use cases and business logic
│       ├── infrastructure/# External concerns (repositories, etc.)
│       ├── presentation/  # HTTP handlers and DTOs
│       └── examples/      # Usage examples for the module
├── pkg/                   # Public library code
├── tests/                 # Test files
├── data/                  # Database files (SQLite)
├── logs/                  # Log files (when using file provider)
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

4. **See logger examples:**
   ```bash
   # Basic logging demo
   go run ./cmd/examples
   
   # With debug level
   APP_LOGGER_LEVEL=debug go run ./cmd/examples
   
   # With JSON format
   APP_LOGGER_FORMAT=json go run ./cmd/examples
   ```

5. **Access the API:**
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

## 📝 Advanced Logging

The application features a sophisticated logging system with OpenTelemetry integration:

### Key Features
- **🔍 OpenTelemetry Integration** - Automatic trace_id and span_id in logs
- **📊 Structured Logging** - JSON and console formats
- **🔌 Multiple Providers** - stdout, file, elasticsearch, logstash
- **🎯 Context-Aware** - Automatic correlation with traces
- **⚡ Performance Metrics** - Built-in timing and metrics

### Logger Configuration

```yaml
logger:
  level: "info"                    # debug, info, warn, error, fatal
  format: "console"                # console, json
  provider: "stdout"               # stdout, file, elasticsearch, logstash
  
  # File provider
  filepath: "./logs/app.log"
  
  # Elasticsearch provider  
  url: "http://localhost:9200"
  index: "boilerplate-go-logs"
  username: "elastic_user"
  password: "elastic_pass"
  api_key: "your_api_key"
  
  # Logstash provider
  url: "localhost:5044"            # TCP endpoint
```

### Usage Examples

```go
// Basic logging with context
logger.LogInfo(ctx, "User created successfully", map[string]interface{}{
    "user_id": 12345,
    "email": "user@example.com",
    "duration": "150ms",
})

// Error logging with error object
logger.LogError(ctx, "Database operation failed", err, map[string]interface{}{
    "operation": "user_create",
    "table": "users",
})

// With OpenTelemetry tracing
ctx, span := otel.Tracer("user-service").Start(ctx, "CreateUser")
defer span.End()

// Logs automatically include trace_id and span_id
logger.LogInfo(ctx, "Processing user", map[string]interface{}{
    "user_id": userID,
    "step": "validation",
})
```

### Log Output Examples

**Console Format:**
```
2024-01-15T10:30:45Z INF User created successfully user_id=12345 email=user@example.com trace_id=4bf92f3577b34da6
```

**JSON Format:**
```json
{
  "level": "info",
  "user_id": 12345,
  "email": "user@example.com", 
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "time": "2024-01-15T10:30:45Z",
  "message": "User created successfully"
}
```

## 🧩 Dependency Injection with FX

The application uses [Uber FX](https://uber-go.github.io/fx/) for clean dependency injection:

### Architecture Benefits
- **🔧 Automatic Wiring** - Dependencies resolved automatically
- **🧪 Easy Testing** - Simple mocking and injection
- **📦 Modular Design** - Components organized in modules
- **🚀 Lifecycle Management** - Proper startup/shutdown handling

### FX Modules Structure

```go
// Application modules
var AppModule = fx.Module("app",
    ConfigModule,     // Configuration
    LoggerModule,     // Advanced logging
    TelemetryModule,  // OpenTelemetry setup
    DatabaseModule,   // Database connection
    UserModule,       // User domain
    ServerModule,     // HTTP server
)
```

### Adding New Modules

1. **Create module definition:**
   ```go
   var NewFeatureModule = fx.Module("new-feature",
       fx.Provide(NewFeatureRepository),
       fx.Provide(NewFeatureService),
       fx.Provide(NewFeatureController),
   )
   ```

2. **Add to AppModule:**
   ```go
   var AppModule = fx.Module("app",
       // ... existing modules
       NewFeatureModule,
   )
   ```

3. **Dependencies are automatically injected!**

## Configuration

Configuration can be managed through:
- `config.yaml` - Main configuration file
- Environment variables (prefixed with `APP_`)
- Command line flags (to be implemented)

### Configuration Hierarchy
1. Environment variables (highest priority)
2. Configuration file
3. Default values (lowest priority)

### Complete Configuration Example

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"                    # debug, release, test

database:
  driver: "postgres"               # postgres, sqlite
  postgres:
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "password"
    database: "boilerplate"
    sslmode: "disable"

logger:
  level: "info"
  format: "console"
  provider: "stdout"
  filepath: "./logs/app.log"
  url: "http://localhost:9200"

telemetry:
  enabled: true
  tracing_enabled: true
  metrics_enabled: true
  endpoint: "http://localhost:4317"
  
application:
  name: "boilerplate-go"
  version: "1.0.0"
  environment: "development"
```

### Environment Variable Examples

```bash
# Logger configuration
export APP_LOGGER_LEVEL=debug
export APP_LOGGER_FORMAT=json
export APP_LOGGER_PROVIDER=elasticsearch
export APP_LOGGER_URL=http://elasticsearch:9200

# Database configuration
export APP_DATABASE_POSTGRES_HOST=db.example.com
export APP_DATABASE_POSTGRES_PASSWORD=secret

# Telemetry configuration
export APP_TELEMETRY_ENABLED=true
export APP_TELEMETRY_ENDPOINT=http://jaeger:4317
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

# Logger examples
make run-examples      # Run logger examples
make run-debug         # Run with debug logging
make run-json          # Run with JSON logging
```

## 🔍 Observability

### OpenTelemetry Integration

The application is fully instrumented with OpenTelemetry for:

- **📈 Distributed Tracing** - Request flow across services
- **📊 Metrics Collection** - Performance and business metrics  
- **📝 Correlated Logging** - Logs linked to traces
- **🎯 Error Tracking** - Detailed error context

### Telemetry Configuration

```yaml
telemetry:
  enabled: true
  tracing_enabled: true
  metrics_enabled: true
  host_metrics_enabled: true
  runtime_metrics_enabled: true
  endpoint: "http://localhost:4317"
  headers: "authorization=Bearer token"
  attributes: "service.name=boilerplate-go,service.version=1.0.0"
```

## Adding New Features

When adding new features, follow the Clean Architecture pattern with FX integration:

1. **Create a new module directory:**
   ```
   internal/orders/
   ├── domain/
   │   ├── order.go              # Entity
   │   └── order_repository.go   # Repository interface
   ├── application/
   │   └── order_service.go      # Business logic
   ├── infrastructure/
   │   └── gorm_order_repository.go  # Repository implementation
   └── presentation/
       └── order_controller.go   # HTTP handlers
   ```

2. **Create FX module:**
   ```go
   var OrderModule = fx.Module("order",
       fx.Provide(infrastructure.NewGormOrderRepository),
       fx.Provide(application.NewOrderService),
       fx.Provide(presentation.NewOrderController),
   )
   ```

3. **Add logging examples:**
   ```go
   func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {
       ctx, span := otel.Tracer("order-service").Start(ctx, "CreateOrder")
       defer span.End()
       
       s.logger.LogInfo(ctx, "Creating order", map[string]interface{}{
           "order_id": order.ID,
           "customer_id": order.CustomerID,
       })
       
       // ... business logic
   }
   ```

4. **Update main module and routes are automatically wired!**

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/config
go test ./internal/logger
go test ./internal/user/application

# Test with different log levels
APP_LOGGER_LEVEL=error go test ./...
```

### Testing with Mocks

The DI architecture makes testing simple:

```go
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    testLogger := createTestLogger()
    service := application.NewUserService(mockRepo, testLogger)
    
    // Test automatically includes logging
}
```

## Production Deployment

1. **Build Docker image:**
   ```bash
   make docker-build
   ```

2. **Set production configuration:**
   ```bash
   # Application
   export APP_SERVER_MODE=release
   export APP_APPLICATION_ENVIRONMENT=production
   
   # Logging
   export APP_LOGGER_LEVEL=info
   export APP_LOGGER_FORMAT=json
   export APP_LOGGER_PROVIDER=elasticsearch
   export APP_LOGGER_URL=https://elasticsearch.company.com
   
   # Telemetry
   export APP_TELEMETRY_ENABLED=true
   export APP_TELEMETRY_ENDPOINT=https://jaeger.company.com:4317
   
   # Database
   export APP_DATABASE_POSTGRES_HOST=prod-db.company.com
   export APP_DATABASE_POSTGRES_PASSWORD=secret
   ```

3. **Deploy with your preferred method**

## 📚 Documentation

- **Logger Examples**: `internal/user/examples/logger_examples.go`
- **FX Configuration**: `internal/fx/fx.go`
- **OpenTelemetry Setup**: `internal/telemetry/telemetry.go`
- **Architecture Patterns**: Follow the existing user module structure

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes following the established patterns
4. Add proper logging with OpenTelemetry traces
5. Include tests with proper DI mocking
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
