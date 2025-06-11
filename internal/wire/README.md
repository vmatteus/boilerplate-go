# Wire Dependency Injection

Este diretório contém a implementação da injeção de dependência usando Google Wire para o projeto.

## 📋 Estrutura dos Arquivos

- `providers.go` - Define todos os providers (funções que criam dependências)
- `wire_sets.go` - Organiza providers em conjuntos lógicos por domínio
- `wire.go` - Define as funções de inicialização usando build tags
- `wire_gen.go` - Código gerado automaticamente pelo Wire (NÃO EDITAR)

## 🔧 Como Funciona

### Providers
Cada provider é uma função que cria uma dependência específica:

```go
func ProvideConfig() (*config.Config, error) {
    return config.Load()
}

func ProvideDatabase(cfg *config.Config) (*gorm.DB, error) {
    return database.Connect(cfg.Database)
}
```

### Wire Sets
Organizam providers relacionados em grupos:

```go
var ConfigSet = wire.NewSet(ProvideConfig)
var DatabaseSet = wire.NewSet(ProvideDatabase)
var ApplicationSet = wire.NewSet(ConfigSet, DatabaseSet, UserSet, ServerSet)
```

### Funções de Inicialização
Definem o que deve ser criado e retornado:

```go
func InitializeApplication(ctx context.Context) (*Application, error) {
    panic(wire.Build(ApplicationSet, wire.Struct(new(Application), "*")))
}
```

## 🚀 Vantagens da Implementação

1. **Testabilidade**: Fácil criação de mocks e testes unitários
2. **Manutenibilidade**: Dependências centralizadas e organizadas
3. **Flexibilidade**: Fácil adição de novos módulos/dependências
4. **Performance**: Resolução de dependências em tempo de compilação
5. **Type Safety**: Detecção de erros em tempo de compilação

## 📝 Como Adicionar Novas Dependências

1. **Criar o Provider**:
```go
func ProvideNewService(dependency SomeDependency) *NewService {
    return NewService{dependency: dependency}
}
```

2. **Adicionar ao Wire Set apropriado**:
```go
var NewDomainSet = wire.NewSet(ProvideNewService)
```

3. **Atualizar o ApplicationSet se necessário**:
```go
var ApplicationSet = wire.NewSet(ConfigSet, DatabaseSet, NewDomainSet, UserSet, ServerSet)
```

4. **Regenerar o código**:
```bash
cd internal/wire && wire
```

## 🔄 Fluxo de Dependências

```
Config
├── Database
├── Logger ──── ZerologLogger
└── TelemetryCleanup

Database ──── UserRepository ──── UserService ──── UserController

Server ← Config + Database + ZerologLogger + UserController
```

## 📋 Comandos Úteis

- **Gerar código**: `cd internal/wire && wire`
- **Verificar dependências**: `wire check ./internal/wire`
- **Limpar código gerado**: `rm wire_gen.go && wire`

## 🧪 Testes

Para testes, você pode criar providers específicos que retornam mocks:

```go
func ProvideTestDatabase() *gorm.DB {
    // Retorna database em memória para testes
}

func InitializeTestApplication() (*Application, error) {
    panic(wire.Build(TestProviders))
}
``` 