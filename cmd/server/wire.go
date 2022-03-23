//+build wireinject

package main

import (
	"github.com/xiaowuzai/payroll/internal/app"
	"github.com/xiaowuzai/payroll/internal/data"
	"github.com/xiaowuzai/payroll/internal/router"
	"github.com/xiaowuzai/payroll/internal/config"
	"github.com/xiaowuzai/payroll/internal/router/handler"
	"github.com/xiaowuzai/payroll/internal/service"
	"github.com/google/wire"
)

func InitServer(conf *config.Server, database *config.Database) (*app.Server, error){
	panic(wire.Build(
		data.ProviderSet,
		//app.ProviderSet,
		service.ProviderSet,
		router.ProviderSet,
		handler.ProviderSet,
		app.NewServer,
		app.NewGinEngine,
		))
}