package wire

import "github.com/google/wire"

// ConfigSet provides configuration-related dependencies
var ConfigSet = wire.NewSet(
	ProvideConfig,
)

// LoggerSet provides logging-related dependencies
var LoggerSet = wire.NewSet(
	ProvideLogger,
	ProvideZerologLogger,
)

// DatabaseSet provides database-related dependencies
var DatabaseSet = wire.NewSet(
	ProvideDatabase,
)

// UserSet provides user domain dependencies
var UserSet = wire.NewSet(
	ProvideUserRepository,
	ProvideUserService,
	ProvideUserController,
)

// ServerSet provides server dependencies
var ServerSet = wire.NewSet(
	ProvideServer,
)

// ApplicationSet combines all application dependencies
var ApplicationSet = wire.NewSet(
	ConfigSet,
	LoggerSet,
	DatabaseSet,
	UserSet,
	ServerSet,
)
