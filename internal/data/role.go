package data

import (
	"context"
	"errors"
	"github.com/xiaowuzai/payroll/internal/service"
	"github.com/xiaowuzai/payroll/pkg/uuid"
	"log"
	"time"
	"xorm.io/xorm"
)

var _ service.RoleRepo = (*roleRepo)(nil)

type roleRepo struct {
	data *Data
}

func NewRoleRepo(data *Data) service.RoleRepo {
	return &roleRepo{
		data:data,
	}
}

type Role struct {
	Id string `xorm:"id varchar(36) pk"`
	Name string `xorm:"name varchar(255)"`
	Description string `xorm:"description varchar(255)"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	Deleted bool `xorm:"default false"`
}


func (r *Role) toService()*service.Role {
	return &service.Role{
		Id :r.Id,
		Name: r.Name,
		Description:r.Description,
		Created: r.Created,
	}
}

func (r *Role) fromService(sRole *service.Role) {
	r.Id  = sRole.Id
	r.Name = sRole.Name
	r.Description = sRole.Description
}

func (rr *roleRepo)AddRole(ctx context.Context, userId string, sRole *service.Role) error {
	session, err := BeginSession(ctx, rr.data.db)
	if err != nil {
		return err
	}
	defer session.Close()

	role := &Role{}
	role.fromService(sRole)
	role.Id= uuid.CreateUUID()

	err = role.insert(ctx, session)
	if err != nil {
		return err
	}

	err = role.insertRoleMenus(ctx, session, sRole.Menus)
	if err != nil {
		return err
	}

	return nil
}

func (rr *roleRepo) UpdateRole(ctx context.Context, sRole *service.Role)  error {
	session := rr.data.db.NewSession()
	if err := session.Begin(); err !=nil {
		log.Println("UpdateRole error : ", err.Error())
		return errors.New("数据库异常")
	}
	defer session.Close()

	// 查询 role
	role := &Role{Id: sRole.Id}
	if err := role.get(ctx,session); err != nil {
		return err
	}

	role.fromService(sRole)
	err := role.update(ctx, session)
	if err != nil {
		return err
	}

	err = role.deleteRoleMenus(ctx, session)
	if err != nil {
		return err
	}

	err = role.insertRoleMenus(ctx, session, sRole.Menus)
	if err != nil {
		return err
	}

	session.Commit()
	return nil
}

func (rr *roleRepo) DeleteRole(ctx context.Context, id string)  error {
	session, err := BeginSession(ctx, rr.data.db)
	if err != nil {
		return err
	}
	defer session.Close()

	role := &Role{Id: id}
	if err := role.get(ctx, session); err != nil {
		return err
	}

	// 查看是否有被用户绑定
	user := &User{RoleId: role.Id}
	has, err := user.getByRoleId(ctx,session)
	if err != nil {
		return err
	}
	if has {
		return errors.New("该角色被用户绑定，不能删除")
	}

	err = role.delete(ctx, session)
	if err != nil {
		return err
	}

	err = role.deleteRoleMenus(ctx, session)
	if err != nil {
		return err
	}

	session.Commit()
	return nil
}

// 获取角色列表
func (rr *roleRepo)ListRole(ctx context.Context, userId string) ([]*service.Role, error) {
	session, err := BeginSession(ctx, rr.data.db)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// 获取 role 列表
	roles, err := listRole(ctx, session)
	if err != nil {
		return nil, err
	}

	sRoles := make([]*service.Role, 0, len(roles))
	for _, role :=  range roles {
		sRoles = append(sRoles, role.toService())
	}

	session.Commit()
	return sRoles, nil
}


func (rr *roleRepo) GetRole(ctx context.Context, userId, roleId string) (*service.Role, error) {
	session, err := BeginSession(ctx, rr.data.db)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// 获取角色
	role := &Role{Id: roleId}
	err = role.get(ctx, session)
	if err != nil {
		return nil, err
	}

	// 获取当前角色选择的 menuId
	menuIds, err := role.listMenuIds(ctx, session)
	if err != nil {
		return nil, err
	}

	sr := role.toService()
	sr.Menus = menuIds

	session.Commit()
	return sr, nil
}

func listRole(ctx context.Context, session *xorm.Session) ([]*Role, error) {
	roles := make([]*Role, 0)
	err := session.Find(&roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *Role) deleteRoleMenus(ctx context.Context, session *xorm.Session) error {
	_, err := session.Where("role_id = ?", r.Id).Delete(&RoleMenu{})
	if err != nil {
		log.Println("deleteRoleMenus error: ", err.Error())
		return errors.New("删除角色对应的权限失败")
	}
	return nil
}

func (r *Role) insertRoleMenus(ctx context.Context, session *xorm.Session, menus []string) error {
	rmenus := make([]*RoleMenu, 0, len(menus))

	for _, menuId :=  range menus {
		if len(menuId) != len(uuid.CreateUUID()) {
			return errors.New("插入数据出错")
		}

		rmenus = append(rmenus, &RoleMenu{
			Id : uuid.CreateUUID(),
			RoleId : r.Id,
			MenuId : menuId,
		})
	}

	_, err := session.Insert(&rmenus)
	if err != nil {
		log.Printf("insertRoleMenus error: %s\n", err.Error())
		return errors.New("插入角色权限失败")
	}
	return nil
}

func (r *Role) listMenuIds(ctx context.Context, session *xorm.Session) ([]string, error) {
	rms := make([]*RoleMenu, 0)
	err := session.Where("role_id = ?", r.Id).Find(&rms)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(rms))
	for _, rm :=  range  rms{
		ids = append(ids, rm.MenuId)
	}
	return ids, nil
}

func (r *Role) update(ctx context.Context, session *xorm.Session) error {
	_, err := session.ID(r.Id).Update(r)
	if err != nil{
		return errors.New("更新角色信息失败")
	}
	return nil
}

func (r *Role)get(ctx context.Context, session *xorm.Session)error {
	has, err := session.ID(r.Id).Get(r)
	if err != nil {
		log.Printf("getRole error: %s\n", err.Error())
		return errors.New("获取角色信息出错")
	}
	if !has {
		return errors.New("该角色不存在")
	}

	return nil
}

func (r *Role) delete(ctx context.Context, session *xorm.Session) error {
	_, err := session.ID(r.Id).Delete(r)
	if err != nil {
		log.Println("delete role  error: ", err.Error())
		return errors.New("删除角色失败")
	}

	return err
}

func (r *Role)insert(ctx context.Context, session *xorm.Session) error{
	_, err := session.Insert(r)
	if err != nil {
		log.Println("Role insert error: ", err.Error())
		return errors.New("插入数据失败")
	}
	return nil
}
