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

func (data *Data) insertRoleMenus(ctx context.Context, roleMenus []*RoleMenu) error {
	_, err := data.db.Insert(&roleMenus)
	return err
}

// 获取 role 对应的菜单列表
func ( data *Data)listRoleMenuIds(ctx context.Context, roleId string) ([]*RoleMenu, error){
	roleMenus := make([]*RoleMenu, 0)
	err := data.db.Where("role_id = ?", roleId).Find(&roleMenus)
	if err != nil {
		return nil, err
	}
	return roleMenus, nil
}

// 获取所有的菜单列表
func (data *Data)listMenus(ctx context.Context) ([]*Menu, error) {
	menus := make([]*Menu,0)
	err := data.db.Find(&menus)
	return menus, err
}