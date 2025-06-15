package fx

import (
	"context"

	"go.uber.org/fx"

	"github.com/rs/zerolog"
	"github.com/your-org/boilerplate-go/internal/config"
	"github.com/your-org/boilerplate-go/internal/database"
	"github.com/your-org/boilerplate-go/internal/logger"
	"github.com/your-org/boilerplate-go/internal/server"
	"github.com/your-org/boilerplate-go/internal/telemetry"
	"github.com/your-org/boilerplate-go/internal/user/application"
	"github.com/your-org/boilerplate-go/internal/user/infrastructure"
	"github.com/your-org/boilerplate-go/internal/user/presentation"
	"gorm.io/gorm"
)

// AppModule fornece todos os módulos da aplicação
var AppModule = fx.Module("app",
	ConfigModule,
	LoggerModule,
	TelemetryModule,
	DatabaseModule,
	UserModule,
	ServerModule,
)

// ConfigModule fornece a configuração
var ConfigModule = fx.Module("config",
	fx.Provide(config.Load),
)

// LoggerModule fornece o logger
var LoggerModule = fx.Module("logger",
	fx.Provide(NewLogger),
	fx.Provide(NewZerologLogger),
)

// TelemetryModule fornece telemetria
var TelemetryModule = fx.Module("telemetry",
	fx.Provide(NewTelemetryCleanup),
)

// DatabaseModule fornece conexão com banco de dados
var DatabaseModule = fx.Module("database",
	fx.Provide(NewDatabase),
	fx.Invoke(RunMigrations),
	fx.Invoke(SetupTracing),
)

// UserModule fornece componentes do domínio User
var UserModule = fx.Module("user",
	fx.Provide(infrastructure.NewGormUserRepository),
	fx.Provide(NewUserService),
	fx.Provide(NewUserController),
)

// ServerModule fornece o servidor HTTP
var ServerModule = fx.Module("server",
	fx.Provide(server.New),
)

// NewLogger adapter para o logger
func NewLogger(cfg *config.Config) *logger.Logger {
	appLogger := logger.InitLogger(cfg.Logger)
	return &appLogger
}

// NewZerologLogger extrai o zerolog.Logger do nosso wrapper
func NewZerologLogger(log *logger.Logger) zerolog.Logger {
	return log.Logger
}

// NewUserService adapter para o service de usuário
func NewUserService(userRepo *infrastructure.GormUserRepository, log *logger.Logger) *application.UserService {
	return application.NewUserService(userRepo, log)
}

// NewTelemetryCleanup adapter para telemetria
func NewTelemetryCleanup(lc fx.Lifecycle, cfg *config.Config) func() {
	if !cfg.Telemetry.Enabled {
		return func() {}
	}

	ctx := context.Background()
	cleanup := telemetry.InitTelemetry(ctx, cfg)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cleanup()
			return nil
		},
	})

	return cleanup
}

// NewDatabase adapter para conexão com banco
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	return database.Connect(cfg.Database)
}

// RunMigrations executa as migrações
func RunMigrations(db *gorm.DB) error {
	return database.Migrate(db)
}

// SetupTracing configura tracing do banco
func SetupTracing(db *gorm.DB, cfg *config.Config) error {
	return database.ConfigureTracing(db, cfg.Telemetry.Enabled)
}

// NewUserController adapter para o controller de usuário
func NewUserController(userService *application.UserService, log *logger.Logger) *presentation.UserController {
	return presentation.NewUserController(userService, log.Logger)
}
