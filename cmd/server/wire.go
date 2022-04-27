//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/xiaowuzai/payroll/internal/app"
	"github.com/xiaowuzai/payroll/internal/config"
	"github.com/xiaowuzai/payroll/internal/data"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/router"
	"github.com/xiaowuzai/payroll/internal/router/handler"
	"github.com/xiaowuzai/payroll/internal/service"
)

func InitServer(conf *config.Payroll, database *config.Database) (*app.Server, error) {
	panic(wire.Build(
		data.ProviderSet,
		service.ProviderSet,
		router.ProviderSet,
		handler.ProviderSet,
		app.NewServer,
		app.NewGinEngine,
		logger.NewLogger,
	))
}
