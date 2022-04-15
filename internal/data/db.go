package data

import (
	"fmt"
	"github.com/xiaowuzai/payroll/internal/config"
	"xorm.io/xorm"
)

func NewDB(conf *config.Database) (*xorm.Engine, error) {
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.Username, conf.Passwd, conf.Host, conf.Port, "payroll")
	engine, err := xorm.NewEngine("mysql", source)
	engine.ShowSQL(conf.ShowSQL)
	if err != nil {
		return nil, err
	}

	err = engine.Ping()
	if err != nil {
		return nil, err
	}

	err = engine.Sync2(new(Bank), new(Employee), new(Menu), new(Organization), new(Role), new(User))
	if err != nil {
		return nil, err
	}

	return engine, nil
}
