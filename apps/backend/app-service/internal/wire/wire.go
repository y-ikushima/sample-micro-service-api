//go:build wireinject
// +build wireinject

package wire

import (
	"sample-micro-service-api/apps/backend/app-service/internal"
	systemsHandler "sample-micro-service-api/apps/backend/app-service/internal/handler/systems"
	systemsService "sample-micro-service-api/apps/backend/app-service/internal/service/systems"
	"sample-micro-service-api/package-go/database"

	"github.com/google/wire"
)

// ProvideDatabaseClient は database.NewClient をラップしてクリーンアップ機能を提供
func ProvideDatabaseClient() (*database.Client, func(), error) {
	client, err := database.NewClient()
	if err != nil {
		return nil, nil, err
	}
	
	cleanup := func() {
		client.Close()
	}
	
	return client, cleanup, nil
}

// Providers
var DatabaseSet = wire.NewSet(
	ProvideDatabaseClient,
)

var ServiceSet = wire.NewSet(
	systemsService.NewService,
)

var HandlerSet = wire.NewSet(
	systemsHandler.NewHandler,
)

var ServerSet = wire.NewSet(
	internal.NewServer,
)

// Wire everything together
var AppSet = wire.NewSet(
	DatabaseSet,
	ServiceSet,
	HandlerSet,
	ServerSet,
)

// InitializeApp は Wire によって自動生成される関数
func InitializeApp() (*internal.Server, func(), error) {
	wire.Build(AppSet)
	return nil, nil, nil
} 