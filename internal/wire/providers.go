package wire

import (
	"github.com/rs/zerolog"
	"github.com/your-org/boilerplate-go/internal/config"
	"github.com/your-org/boilerplate-go/internal/database"
	"github.com/your-org/boilerplate-go/internal/server"
	"github.com/your-org/boilerplate-go/internal/user/application"
	"github.com/your-org/boilerplate-go/internal/user/domain"
	"github.com/your-org/boilerplate-go/internal/user/infrastructure"
	"github.com/your-org/boilerplate-go/internal/user/presentation"
	"github.com/your-org/boilerplate-go/pkg/logger"
	"gorm.io/gorm"
)

// ProvideConfig provides application configuration
func ProvideConfig() (*config.Config, error) {
	return config.Load()
}

// ProvideLogger provides application logger
func ProvideLogger(cfg *config.Config) logger.Logger {
	return logger.InitLogger(cfg.Logger)
}

// ProvideZerologLogger provides zerolog.Logger from logger.Logger
func ProvideZerologLogger(appLogger logger.Logger) zerolog.Logger {
	return appLogger.Logger
}

// ProvideDatabase provides database connection
func ProvideDatabase(cfg *config.Config) (*gorm.DB, error) {
	return database.Connect(cfg.Database)
}

// ProvideUserRepository provides user repository
func ProvideUserRepository(db *gorm.DB) domain.UserRepository {
	return infrastructure.NewGormUserRepository(db)
}

// ProvideUserService provides user service
func ProvideUserService(userRepo domain.UserRepository) *application.UserService {
	return application.NewUserService(userRepo)
}

// ProvideUserController provides user controller
func ProvideUserController(userService *application.UserService, logger zerolog.Logger) *presentation.UserController {
	return presentation.NewUserController(userService, logger)
}

// ProvideServer provides HTTP server
func ProvideServer(cfg *config.Config, db *gorm.DB, logger zerolog.Logger, userController *presentation.UserController) *server.Server {
	return server.New(cfg, db, logger, userController)
}
