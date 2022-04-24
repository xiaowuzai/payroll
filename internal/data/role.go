package data

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/service"
	"github.com/xiaowuzai/payroll/pkg/uuid"
	"time"
	"xorm.io/xorm"
)

var _ service.RoleRepo = (*roleRepo)(nil)

type roleRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewRoleRepo(data *Data, logger *logger.Logger) service.RoleRepo {
	return &roleRepo{
		data:   data,
		logger: logger,
	}
}

func (rr *roleRepo) AddRole(ctx context.Context, userId string, sRole *service.Role) error {
	session, err := BeginSession(ctx, rr.data.db, rr.logger)
	if err != nil {
		return err
	}
	defer session.Close()

	role := &Role{}
	role.fromService(sRole)
	role.Id = uuid.CreateUUID()

	err = role.insert(ctx, session, rr.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	err = role.insertRoleMenus(ctx, session, rr.logger, sRole.Menus)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_ = session.Commit()
	return nil
}

func (rr *roleRepo) UpdateRole(ctx context.Context, sRole *service.Role) error {
	session, err := BeginSession(ctx, rr.data.db, rr.logger)
	if err != nil {
		return err
	}

	defer session.Close()

	// 查询 role
	role := &Role{Id: sRole.Id}
	has, err := role.get(ctx, session, rr.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	if !has {
		_ = session.Rollback()
		return nil
	}

	role.fromService(sRole)
	err = role.update(ctx, session, rr.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	err = role.deleteRoleMenus(ctx, session, rr.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	err = role.insertRoleMenus(ctx, session, rr.logger, sRole.Menus)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_ = session.Commit()
	return nil
}

func (rr *roleRepo) DeleteRole(ctx context.Context, id string) error {
	session, err := BeginSession(ctx, rr.data.db, rr.logger)
	if err != nil {
		return err
	}
	defer session.Close()

	log := rr.logger.WithRequestId(ctx)

	role := &Role{Id: id}
	has, err := role.get(ctx, session, rr.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	if !has {
		message := fmt.Sprintf("Role id = %s, not found", id)
		log.Info(message)
		_ = session.Commit()
		return nil
	}

	// 查看是否有被用户绑定
	user := &User{RoleId: role.Id}
	has, err = user.hasByRoleId(ctx, session, rr.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	if has {
		message := "该角色被用户绑定，不能删除"
		log.Info(message)
		_ = session.Rollback()
		return errors.New(message)
	}

	err = role.delete(ctx, session, rr.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	err = role.deleteRoleMenus(ctx, session, rr.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_ = session.Commit()
	return nil
}

// 获取角色列表
func (rr *roleRepo) ListRole(ctx context.Context, userId string) ([]*service.Role, error) {
	session := NewSession(ctx, rr.data.db)

	// 获取 role 列表
	role := &Role{}
	roles, err := role.list(ctx, session, rr.logger)
	if err != nil {
		return nil, err
	}

	sRoles := make([]*service.Role, 0, len(roles))
	for _, role := range roles {
		sRoles = append(sRoles, role.toService())
	}

	return sRoles, nil
}

func (rr *roleRepo) GetRole(ctx context.Context, userId, roleId string) (*service.Role, error) {
	session, err := BeginSession(ctx, rr.data.db, rr.logger)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// 获取角色
	role := &Role{Id: roleId}
	has, err := role.get(ctx, session, rr.logger)
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}
	if !has {
		message := fmt.Sprintf("Role id = %s not found", roleId)
		log.Errorf(message)
		_ = session.Rollback()
		return nil, errors.DataNotFound(message)
	}

	// 获取当前角色选择的 menuId
	menuIds, err := role.listMenuIds(ctx, session, rr.logger)
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}

	sr := role.toService()
	sr.Menus = menuIds

	_ = session.Commit()
	return sr, nil
}

type Role struct {
	Id          string    `xorm:"id varchar(36) pk"`
	Name        string    `xorm:"name varchar(255)"`
	Description string    `xorm:"description varchar(255)"`
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
	Deleted     bool      `xorm:"default false"`
}

func (r *Role) toService() *service.Role {
	return &service.Role{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Created:     r.Created,
	}
}

func (r *Role) fromService(sRole *service.Role) {
	r.Id = sRole.Id
	r.Name = sRole.Name
	r.Description = sRole.Description
}

func (r *Role) deleteRoleMenus(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	_, err := session.Where("role_id = ?", r.Id).Delete(&RoleMenu{})
	if err != nil {
		log.Println("deleteRoleMenus error: ", err.Error())
		return errors.New("删除角色对应的权限失败")
	}
	return nil
}

func (r *Role) insertRoleMenus(ctx context.Context, session *xorm.Session, logger *logger.Logger, menus []string) error {
	rmenus := make([]*RoleMenu, 0, len(menus))

	for _, menuId := range menus {
		if len(menuId) != len(uuid.CreateUUID()) {
			return errors.New("插入数据出错")
		}

		rmenus = append(rmenus, &RoleMenu{
			Id:     uuid.CreateUUID(),
			RoleId: r.Id,
			MenuId: menuId,
		})
	}

	roleMenu := &RoleMenu{}

	return roleMenu.insertList(ctx, session, logger, rmenus)
}

func (r *Role) listMenuIds(ctx context.Context, session *xorm.Session, logger *logger.Logger) ([]string, error) {
	roleMenu := &RoleMenu{}
	rms, err := roleMenu.list(ctx, session, logger, r.Id)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(rms))
	for _, rm := range rms {
		ids = append(ids, rm.MenuId)
	}
	return ids, nil
}

func (r *Role) update(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("Role update input %+v\n", *r)

	_, err := session.ID(r.Id).Update(r)
	if err != nil {
		log.Errorf("Role update error: %s\n", err.Error())
		return errors.ErrDataUpdate(err.Error())
	}
	return nil
}

func (r *Role) get(ctx context.Context, session *xorm.Session, logger *logger.Logger) (bool, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("Role get input %+v\n", *r)

	has, err := session.ID(r.Id).Get(r)
	if err != nil {
		log.Errorf("Role get error: %s\n", err.Error())
		return false, errors.ErrDataGet(err.Error())
	}

	return has, nil
}

func (r *Role) delete(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("Role delete input %+v\n", *r)

	_, err := session.ID(r.Id).Delete(r)
	if err != nil {
		log.Error("delete role  error: ", err.Error())
		return errors.ErrDataDelete(err.Error())
	}

	return nil
}

func (r *Role) insert(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("Role insert input %+v\n", *r)

	_, err := session.Insert(r)
	if err != nil {
		log.Errorf("Role insert error: %s\n", err.Error())
		return errors.ErrDataInsert(err.Error())
	}
	return nil
}

func (r *Role) list(ctx context.Context, session *xorm.Session, logger *logger.Logger) ([]*Role, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("Role list input \n")

	roles := make([]*Role, 0)
	err := session.Find(&roles)
	if err != nil {
		log.Errorf("Role list error: %s\n", err.Error())
		return nil, err
	}
	return roles, nil
}
