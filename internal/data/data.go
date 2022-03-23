package data

import (
	"fmt"
	"github.com/xiaowuzai/payroll/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRoleRepo, NewOrganizationRepo)

type Data struct {
	db *xorm.Engine
}

func NewData(conf *config.Database) (*Data, error) {
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.Username, conf.Passwd, conf.Host, conf.Port, "payroll")
	engine, err := xorm.NewEngine("mysql", source)
	if err != nil {
		return nil, err
	}

	err = engine.Ping()
	if err != nil {
		return nil, err
	}

	return &Data{
		db:engine,
	},nil
}