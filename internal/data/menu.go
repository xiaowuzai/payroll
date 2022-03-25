package data

import (
	"context"
	"time"
)

type Menu struct {
	Id string `xorm:"id"`
	Name string `xorm:"name"`
	Keys string `xorm:"keys"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

type RoleMenu struct {
	Id string `xorm:"id"`
	RoleId string `xorm:"role_id"`
	MenuId string `xorm:"menu_id"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func (rr *roleRepo) insertRoleMenus(ctx context.Context, roleMenus []*RoleMenu) error {
	_, err := rr.data.db.Insert(&roleMenus)
	return err
}

func (rr *roleRepo) listRoleMenuIds(ctx context.Context, roleId string) ([]*RoleMenu, error){
	roleMenus := make([]*RoleMenu, 0)
	err := rr.data.db.Where("role_id = ?", roleId).Find(roleMenus)
	if err != nil {
		return nil, err
	}
	return roleMenus, nil
}

func (rr *roleRepo)listMenus(ctx context.Context) ([]*Menu, error) {
	menus := make([]*Menu,0)
	err := rr.data.db.Find(&menus)
	return menus, err
}