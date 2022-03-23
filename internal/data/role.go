package data

import (
	"context"
	"errors"
	"github.com/xiaowuzai/payroll/internal/service"
	"github.com/xiaowuzai/payroll/pkg/uuid"
	"time"
)

var _ service.RoleRepo = (*roleRepo)(nil)

type Role struct {
	Id string `xorm:"id"`
	Description string `xorm:"description"`
	Roles string `xorm:"roles"`  //eg: key1.key2.key3
	Created time.Time `xorm:"created"`
}

type roleRepo struct {
	data *Data
}

func NewRoleRepo(data *Data) service.RoleRepo {
	return &roleRepo{
		data:data,
	}
}

func (rr *roleRepo)AddRole(ctx context.Context, bRole *service.Role) error {
	role := &Role{
		Id: uuid.CreateUUID(),
		Description: bRole.Description,
		Roles: bRole.Roles,
	}

	n, err := rr.data.db.Insert(role)
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("insert data fail")
	}

	return nil
}
