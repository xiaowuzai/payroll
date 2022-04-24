package data

import (
	"fmt"
	"github.com/xiaowuzai/payroll/internal/config"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"xorm.io/xorm"
)

func NewDB(conf *config.Database, logger *logger.Logger) (*xorm.Engine, error) {
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.Username, conf.Passwd, conf.Host, conf.Port, "payroll")
	engine, err := xorm.NewEngine("mysql", source)
	if err != nil {
		logger.Error("NewDB NewEngine error: ", err.Error())
		return nil, err
	}

	engine.ShowSQL(conf.ShowSQL)

	err = engine.Ping()
	if err != nil {
		logger.Error("NewDB Engine Ping error: ", err.Error())
		return nil, err
	}

	err = engine.Sync2(
		new(Bank),
		new(Employee),
		new(Menu),
		new(RoleMenu),
		new(Organization),
		new(Role),
		new(User),
	)
	if err != nil {
		logger.Error("NewDB Sync table error: ", err.Error())
		return nil, err
	}

	logger.Infof("NewDB success")
	return engine, nil
}
