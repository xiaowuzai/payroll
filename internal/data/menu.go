package data

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/service"
	"time"
	"xorm.io/xorm"
)

var _ service.MenuRepo = (*MenuRepo)(nil)

type MenuRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewMenuRepo(data *Data, logger *logger.Logger) service.MenuRepo {
	return &MenuRepo{
		data:   data,
		logger: logger,
	}
}

type Menu struct {
	Id      string    `xorm:"id varchar(36) pk"`
	Name    string    `xorm:"name notnull"`
	Keys    string    `xorm:"keys"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

//func (data *Data) insertRoleMenus(ctx context.Context, roleMenus []*RoleMenu) error {
//	_, err := data.db.Insert(&roleMenus)
//	return err
//}

// 获取 role 对应的菜单列表
//func ( data *Data)listRoleMenuIds(ctx context.Context, roleId string) ([]*RoleMenu, error){
//	roleMenus := make([]*RoleMenu, 0)
//	err := data.db.Where("role_id = ?", roleId).Find(&roleMenus)
//	if err != nil {
//		return nil, err
//	}
//	return roleMenus, nil
//}

//// 获取所有的菜单列表
//func (data *Data)listMenus(ctx context.Context) ([]*Menu, error) {
//	menus := make([]*Menu,0)
//	err := data.db.Find(&menus)
//	return menus, err
//}

// map[key]id
func menuToMap(menus []*Menu) map[string]string {
	menusMap := make(map[string]string, len(menus))
	for _, v := range menus {
		menusMap[v.Keys] = v.Id
	}
	return menusMap
}

func (mr *MenuRepo) ListMenu(ctx context.Context) (*service.Menu, error) {
	session := NewSession(ctx, mr.data.db)

	menu := &Menu{}
	menus, err := menu.list(ctx, session, mr.logger)
	if err != nil {
		return nil, err
	}

	menusMap := menuToMap(menus)
	return &service.Menu{
		MenuKeys: menusMap,
	}, nil
}

func (m *Menu) insert(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("RoleMenu insert list input")

	return nil
}

func (m *Menu) delete(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("RoleMenu insert list input")

	return nil
}

func (m *Menu) update(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("RoleMenu insert list input")

	return nil
}

func (m *Menu) get(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("RoleMenu insert list input")

	return nil
}

func (m *Menu) list(ctx context.Context, session *xorm.Session, logger *logger.Logger) ([]*Menu, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("RoleMenu insert list input")

	menus := make([]*Menu, 0)
	err := session.Find(menus)
	if err != nil {
		log.Errorf("Menu list error: %s\n", err.Error())
		return menus, errors.ErrDataGet(err.Error())
	}

	return menus, nil
}

type RoleMenu struct {
	Id      string    `xorm:"id varchar(36) pk"`
	RoleId  string    `xorm:"role_id varchar(36) notnull"`
	MenuId  string    `xorm:"menu_id varchar(36) notnull"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func (rm *RoleMenu) insert(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("RoleMenu insert input: %+v\n", *rm)

	_, err := session.Insert(rm)
	if err != nil {
		log.Errorf("RoleMenu insert error: %s\n", err.Error())
		return errors.ErrDataInsert(err.Error())
	}

	return nil
}

func (rm *RoleMenu) delete(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("RoleMenu delete input: %+v\n", *rm)

	_, err := session.Delete(rm)
	if err != nil {
		log.Errorf("RoleMenu delete error: %s\n", err.Error())
		return errors.ErrDataDelete(err.Error())
	}

	return nil
}

func (rm *RoleMenu) update(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("RoleMenu update input: %+v\n", *rm)

	_, err := session.Update(rm)
	if err != nil {
		log.Errorf("RoleMenu update error: %s\n", err.Error())
		return errors.ErrDataUpdate(err.Error())
	}

	return nil
}

func (rm *RoleMenu) get(ctx context.Context, session *xorm.Session, logger *logger.Logger) (bool, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("RoleMenu get input: %+v\n", *rm)

	has, err := session.Get(rm)
	if err != nil {
		log.Error("RoleMenu get error: ", err.Error())
		return false, errors.ErrDataGet(err.Error())
	}

	return has, nil
}

func (rm *RoleMenu) list(ctx context.Context, session *xorm.Session, logger *logger.Logger, roleId string) ([]*RoleMenu, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("RoleMenu list input, roleId = %s\n", roleId)

	roleMenus := make([]*RoleMenu, 0)
	err := session.Where("role_id = ?", roleId).Find(&roleMenus)
	if err != nil {
		log.Error("RoleMenu list error: %s\n", err.Error())
		return nil, errors.ErrDataGet(err.Error())
	}

	return roleMenus, nil
}

func (rm *RoleMenu) insertList(ctx context.Context, session *xorm.Session, logger *logger.Logger, roleMenus []*RoleMenu) error {
	log := logger.WithRequestId(ctx)
	log.Infof("RoleMenu insert list input")

	_, err := session.Insert(&roleMenus)
	if err != nil {
		log.Errorf("RoleMenu insert list error %s\n", err.Error())
		return errors.ErrDataInsert(err.Error())
	}

	return nil
}
