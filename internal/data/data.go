package data

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"xorm.io/xorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRoleRepo, NewOrganizationRepo, NewUserRepo, NewMenuRepo, NewDB)

type Data struct {
	db *xorm.Engine
}

func NewData(db *xorm.Engine) (*Data, error) {
	return &Data{
		db: db,
	}, nil
}
