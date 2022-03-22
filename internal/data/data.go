package data

import (
	"fmt"
	"github.com/xiaowuzai/payroll/internal/config"
	"xorm.io/xorm"
)

type Data struct {
	db *xorm.Engine
}

func NewData(conf config.Database) (*Data, error) {
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