package data

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/service"
	"time"
)

var _ service.MenuRepo = (*MenuRepo)(nil)

type Menu struct {
	Id string `xorm:"id"`
	Name string `xorm:"name"`
	Keys string `xorm:"keys"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

type RoleMenu struct {
	Id string `xorm:"id varchar(36) pk"`
	RoleId string `xorm:"role_id varchar(36) notnull"`
	MenuId string `xorm:"menu_id varchar(36) notnull"`
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

type MenuRepo struct {
	data *Data
}

// map[key]id
func menuToMap(menus []*Menu) map[string]string{
	menusMap := make(map[string]string, len(menus))
	for _, v := range menus {
		menusMap[v.Keys] = v.Id
	}
	return menusMap
}

func NewMenuRepo(data *Data) service.MenuRepo{
	return &MenuRepo{
		data: data,
	}
}


func (mr *MenuRepo) ListMenu(ctx context.Context) (*service.Menu, error) {
	menus, err := mr.data.listMenus(ctx)
	if err != nil {
		return nil, err
	}

	menusMap := menuToMap(menus)

	return &service.Menu{
		MenuKeys: menusMap,
	},nil
}
