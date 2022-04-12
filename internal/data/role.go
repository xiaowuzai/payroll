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
	Id string `xorm:"id varchar(36) pk"`
	Name string `xorm:"name varchar(255)"`
	Description string `xorm:"description varchar(255)"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	Deleted bool `xorm:"default false"`
}

type roleRepo struct {
	data *Data
}

func NewRoleRepo(data *Data) service.RoleRepo {
	return &roleRepo{
		data:data,
	}
}

func (rr *roleRepo)AddRole(ctx context.Context, userId string, sRole *service.Role) error {
	role := &Role{
		Id: uuid.CreateUUID(),
		Description: sRole.Description,
		Name: sRole.Name,
	}

	n, err := rr.data.db.Insert(role)
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("insert data fail")
	}

	roleMenus := make([]*RoleMenu,0, len(sRole.Menus))
	for _, v := range sRole.Menus {
		roleMenus = append(roleMenus, &RoleMenu{
			Id: uuid.CreateUUID(),
			RoleId: role.Id,
			MenuId: v,
		})
	}

	return rr.data.insertRoleMenus(ctx, roleMenus)
}

// 获取角色列表
func (rr *roleRepo)ListRole(ctx context.Context, userId string) ([]*service.Role, error) {
	roles, err := rr.listRole(ctx)
	if err != nil {
		return nil, err
	}

	sRoles := make([]*service.Role, 0, len(roles))
	for _, v :=  range roles {
		sRoles = append(sRoles, &service.Role{
			Id: v.Id,
			Name:v.Name,
			Description: v.Description,
			Created: v.Created,
		})
	}
	return sRoles, nil
}


func (rr *roleRepo) GetRole(ctx context.Context, userId, roleId string) (*service.Role, error) {
	role, err := rr.getRole(ctx, roleId)
	if err != nil {
		return nil, err
	}

	// 获取当前角色选择的 menuId
	roleMenus, err := rr.data.listRoleMenuIds(ctx, roleId)
	if err != nil {
		return nil, err
	}
	menuIds := make([]string,0, len(roleMenus))
	for _, v := range roleMenus {
		menuIds = append(menuIds,v.MenuId)
	}

	// 获取所有 menus
	menus, err := rr.data.listMenus(ctx)
	if err != nil {
		return nil, err
	}
	menusMap := make(map[string]string, len(menus))
	for _, v := range menus {
		menusMap[v.Keys] = v.Id
	}

	sRole := &service.Role{
		Id: role.Id,
		Name: role.Name,
		Description: role.Description,
		MenuKey: menusMap,
		Menus: menuIds,
		Created: role.Created,
	}
	return sRole, nil
}

func (rr *roleRepo) listRole(ctx context.Context) ([]*Role, error) {
	roles := make([]*Role, 0)
	err := rr.data.db.Find(&roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (rr *roleRepo)getRole(ctx context.Context, roleId string) (*Role, error) {
	role :=  &Role{}
	_, err := rr.data.db.ID(roleId).Get(role)
	return role, err
}

