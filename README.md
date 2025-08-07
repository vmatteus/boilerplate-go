# Go Boilerplate

Um boilerplate Go limpo e bem estruturado seguindo os princípios da Clean Architecture, com capacidades avançadas de logging e injeção de dependência.

Este projeto segue o padrão de layout recomendado pelo [golang-standards/project-layout](https://github.com/golang-standards/project-layout) para organização de projetos Go.

## 🚀 Funcionalidades

- 🏗️ **Clean Architecture** - Organizado em camadas (Domínio, Aplicação, Infraestrutura, Apresentação)
- 🚀 **Gin Framework** - Framework web HTTP rápido
- 🗄️ **GORM** - ORM rico em funcionalidades com suporte para PostgreSQL e SQLite
- ⚙️ **Configuração Viper** - Gerenciamento de configuração flexível
- 📝 **Logging Avançado** - Logging estruturado com integração OpenTelemetry
- 🔍 **OpenTelemetry Ready** - Observabilidade completa com rastreamento e métricas
- 🧩 **Injeção de Dependência** - DI limpa com Uber FX
- 🐳 **Suporte Docker** - Configuração Docker pronta para uso
- 🔄 **Hot Reload** - Configuração de desenvolvimento com Air
- 🧪 **Pronto para Testes** - Estruturado para testes fáceis com mocks
- 📊 **Health Checks** - Endpoints de verificação de saúde integrados
- 🔌 **Múltiplos Provedores de Log** - suporte para stdout, arquivo, elasticsearch, logstash

## Comandos Disponíveis

```bash
make help              # Mostra todos os comandos disponíveis
make build             # Compila a aplicação
make run               # Executa a aplicação
make test              # Executa os testes
make test-coverage     # Executa os testes com cobertura
make dev               # Executa com hot reload
make docker-up         # Inicia com Docker Compose
make docker-down       # Para os serviços Docker
make clean             # Limpa artefatos de build
make fmt               # Formata o código
make lint              # Executa o linter

# Exemplos de logger
make run-examples      # Executa exemplos de logger
make run-debug         # Executa com nível debug
make run-json          # Executa com formato JSON
```

## Estrutura do Projeto

Este projeto segue o padrão de layout recomendado pelo [golang-standards/project-layout](https://github.com/golang-standards/project-layout) para organização de projetos Go.

```text
├── cmd/                    # Pontos de entrada da aplicação
│   ├── main.go             # Ponto de entrada principal
│   └── examples/           # Exemplos de uso e demos
├── internal/               # Código privado da aplicação
│   ├── config/             # Gerenciamento de configuração
│   ├── database/           # Conexão com banco de dados e migrações
│   ├── fx/                 # Configuração de injeção de dependência
│   ├── logger/             # Logging avançado com OpenTelemetry
│   │   ├── logger.go
│   │   ├── stdout_logger.go
│   │   ├── file_logger.go
│   │   ├── elasticsearch_logger.go
│   │   └── logstash_logger.go
│   ├── middleware/         # Middleware HTTP
│   ├── server/             # Configuração do servidor HTTP
│   ├── telemetry/          # Configuração OpenTelemetry
│   └── [modules]/          # Módulos de funcionalidades (domain-driven)
│       ├── domain/         # Entidades e interfaces de negócio
│       ├── application/    # Casos de uso e lógica de negócio
│       ├── infrastructure/ # Repositórios e integrações externas
│       ├── presentation/   # Handlers HTTP e DTOs
│       └── examples/       # Exemplos de uso do módulo
├── pkg/                    # Código de biblioteca pública
├── tests/                  # Arquivos de teste
├── data/                   # Arquivos de banco de dados (SQLite)
├── logs/                   # Arquivos de log (se usar provedor de arquivo)
├── config.yaml             # Arquivo de configuração
├── docker-compose.yml      # Configuração Docker Compose
└── Makefile                # Comandos de build e desenvolvimento
```

## Início Rápido

### Pré-requisitos

- Go 1.22 ou superior
- Make (opcional, para comandos de conveniência)
- Docker & Docker Compose (opcional)

### Desenvolvimento Local

1. **Clone e configure:**
   ```bash
   git clone <seu-repo>
   cd boilerplate-go
   cp config.example.yaml config.yaml
   cp .env.example .env
   ```

2. **Instale as dependências:**
   ```bash
   go mod download
   ```

3. **Execute a aplicação:**
   ```bash
   make run
   # ou
   go run ./cmd
   ```

4. **Veja os exemplos de logger:**
   ```bash
   # Demo básico de logging
   go run ./cmd/examples
   
   # Com nível debug
   APP_LOGGER_LEVEL=debug go run ./cmd/examples
   
   # Com formato JSON
   APP_LOGGER_FORMAT=json go run ./cmd/examples
   ```

5. **Acesse a API:**
   - Health check: http://localhost:8080/health
   - Endpoint de boas-vindas: http://localhost:8080/api/v1/

### Usando Docker

1. **Inicie com Docker Compose:**
   ```bash
   make docker-up
   ```

2. **Pare os serviços:**
   ```bash
   make docker-down
   ```

### Desenvolvimento com Hot Reload

1. **Instale o Air:**
   ```bash
   make install-dev-tools
   ```

2. **Inicie o servidor de desenvolvimento:**
   ```bash
   make dev
   ```

## 📝 Logging Avançado

A aplicação possui um sistema de logging sofisticado com integração OpenTelemetry:

### Funcionalidades Principais
- **🔍 Integração OpenTelemetry** - trace_id e span_id automáticos nos logs
- **📊 Logging Estruturado** - Formatos JSON e console
- **🔌 Múltiplos Provedores** - stdout, arquivo, elasticsearch, logstash
- **🎯 Consciente do Contexto** - Correlação automática com traces
- **⚡ Métricas de Performance** - Timing e métricas integrados

### Configuração do Logger

```yaml
logger:
  level: "info"                    # debug, info, warn, error, fatal
  format: "console"                # console, json
  provider: "stdout"               # stdout, file, elasticsearch, logstash
  
  # Provedor de arquivo
  filepath: "./logs/app.log"
  
  # Provedor Elasticsearch  
  url: "http://localhost:9200"
  index: "boilerplate-go-logs"
  username: "elastic_user"
  password: "elastic_pass"
  api_key: "your_api_key"
  
  # Provedor Logstash
  url: "localhost:5044"            # Endpoint TCP
```

### Exemplos de Uso

```go
// Logging básico com contexto
logger.LogInfo(ctx, "Usuário criado com sucesso", map[string]interface{}{
    "user_id": 12345,
    "email": "user@example.com",
    "duration": "150ms",
})

// Logging de erro com objeto de erro
logger.LogError(ctx, "Operação do banco de dados falhou", err, map[string]interface{}{
    "operation": "user_create",
    "table": "users",
})

// Com rastreamento OpenTelemetry
ctx, span := otel.Tracer("user-service").Start(ctx, "CreateUser")
defer span.End()

// Logs incluem automaticamente trace_id e span_id
logger.LogInfo(ctx, "Processando usuário", map[string]interface{}{
    "user_id": userID,
    "step": "validation",
})
```

### Exemplos de Saída de Log

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

## 🧩 Injeção de Dependência com FX

A aplicação usa [Uber FX](https://uber-go.github.io/fx/) para injeção de dependência limpa:

### Benefícios da Arquitetura
- **🔧 Conexão Automática** - Dependências resolvidas automaticamente
- **🧪 Testes Facilitados** - Mock e injeção simplificados
- **📦 Design Modular** - Componentes organizados em módulos
- **🚀 Gerenciamento de Ciclo de Vida** - Controle adequado de inicialização/finalização

### Estrutura de Módulos FX

```go
// Módulos da aplicação
var AppModule = fx.Module("app",
    ConfigModule,     // Configuração
    LoggerModule,     // Logging avançado
    TelemetryModule,  // Configuração OpenTelemetry
    DatabaseModule,   // Conexão com banco de dados
    UserModule,       // Domínio de usuário
    ServerModule,     // Servidor HTTP
)
```

### Adicionando Novos Módulos

1. **Crie a definição do módulo:**
   ```go
   var NovoModuloFeature = fx.Module("nova-feature",
       fx.Provide(NewFeatureRepository),
       fx.Provide(NewFeatureService),
       fx.Provide(NewFeatureController),
   )
   ```

2. **Adicione ao AppModule:**
   ```go
   var AppModule = fx.Module("app",
       // ... módulos existentes
       NovoModuloFeature,
   )
   ```

3. **As dependências são injetadas automaticamente!**

## Configuração

A configuração pode ser gerenciada através de:
- `config.yaml` - Arquivo principal de configuração
- Variáveis de ambiente (prefixadas com `APP_`)
- Flags de linha de comando (a ser implementado)

### Hierarquia de Configuração
1. Variáveis de ambiente (prioridade mais alta)
2. Arquivo de configuração
3. Valores padrão (prioridade mais baixa)

### Exemplo Completo de Configuração

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

## 🔍 Observabilidade

### Integração OpenTelemetry

A aplicação está totalmente instrumentada com OpenTelemetry para:

- **📈 Rastreamento Distribuído** - Fluxo de requisições entre serviços
- **📊 Coleta de Métricas** - Métricas de performance e negócio
- **📝 Logging Correlacionado** - Logs vinculados a traces
- **🎯 Rastreamento de Erros** - Contexto detalhado de erros

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

## Adicionando Novas Funcionalidades

Ao adicionar novas funcionalidades, siga o padrão Clean Architecture com integração FX:

1. **Crie um novo diretório de módulo:**
   ```
   internal/orders/
   ├── domain/
   │   ├── order.go              # Entidade
   │   └── order_repository.go   # Interface do repositório
   ├── application/
   │   └── order_service.go      # Lógica de negócio
   ├── infrastructure/
   │   └── gorm_order_repository.go  # Implementação do repositório
   └── presentation/
       └── order_controller.go   # Handlers HTTP
   ```

2. **Crie o módulo FX:**
   ```go
   var OrderModule = fx.Module("order",
       fx.Provide(infrastructure.NewGormOrderRepository),
       fx.Provide(application.NewOrderService),
       fx.Provide(presentation.NewOrderController),
   )
   ```

3. **Adicione exemplos de logging:**
   ```go
   func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {
       ctx, span := otel.Tracer("order-service").Start(ctx, "CreateOrder")
       defer span.End()
       
       s.logger.LogInfo(ctx, "Criando pedido", map[string]interface{}{
           "order_id": order.ID,
           "customer_id": order.CustomerID,
       })
       
       // ... lógica de negócio
   }
   ```

4. **Atualize o módulo principal e as rotas serão automaticamente conectadas!**

## Testes

```bash
# Execute todos os testes
make test

# Execute os testes com cobertura
make test-coverage

# Execute testes de pacotes específicos
go test ./internal/config
go test ./internal/logger
go test ./internal/user/application

# Teste com diferentes níveis de log
APP_LOGGER_LEVEL=error go test ./...
```

### Testando com Mocks

A arquitetura de DI facilita os testes:

```go
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    testLogger := createTestLogger()
    service := application.NewUserService(mockRepo, testLogger)
    
    // O teste inclui automaticamente o logging
}
```

## Implantação em Produção

1. **Crie a imagem Docker:**
   ```bash
   make docker-build
   ```

2. **Defina a configuração de produção:**
   ```bash
   # Aplicação
   export APP_SERVER_MODE=release
   export APP_APPLICATION_ENVIRONMENT=production
   
   # Logging
   export APP_LOGGER_LEVEL=info
   export APP_LOGGER_FORMAT=json
   export APP_LOGGER_PROVIDER=elasticsearch
   export APP_LOGGER_URL=https://elasticsearch.company.com
   
   # Telemetria
   export APP_TELEMETRY_ENABLED=true
   export APP_TELEMETRY_ENDPOINT=https://jaeger.company.com:4317
   
   # Banco de Dados
   export APP_DATABASE_POSTGRES_HOST=prod-db.company.com
   export APP_DATABASE_POSTGRES_PASSWORD=secret
   ```

3. **Implante com o método de sua preferência**

## 📚 Documentation

- **Logger Examples**: `internal/user/examples/logger_examples.go`
- **FX Configuration**: `internal/fx/fx.go`
- **OpenTelemetry Setup**: `internal/telemetry/telemetry.go`
- **Architecture Patterns**: Follow the existing user module structure

## Contribuindo

1. Faça um fork do repositório
2. Crie um branch para a funcionalidade
3. Faça suas alterações seguindo os padrões estabelecidos
4. Adicione logging adequado com traces OpenTelemetry
5. Inclua testes com mocks adequados para DI
6. Envie um pull request
