package data

import "context"

func (ur *userRepo)getMenus(ctx context.Context, user *User)([]string, error) {

	menuIds := make([]string, 0)
	rmList, err := ur.data.listRoleMenuIds(ctx,user.RoleId)
	if err != nil {
		return nil, err
	}

	for _, v := range rmList {
		menuIds = append(menuIds, v.MenuId)
	}

	return menuIds, nil
}
