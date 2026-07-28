//go:build wireinject
// +build wireinject

package container

import (
	"github.com/google/wire"
	"mk-pos-billing/internal/api/routes"
)

// InitializeRouteConfig constructs the RouteConfig needed by the API server.
func InitializeRouteConfig() (routes.RouteConfig, error) {
	wire.Build(
		ApplicationSet,
		wire.Struct(new(routes.RouteConfig), "*"),
	)
	return routes.RouteConfig{}, nil
}
