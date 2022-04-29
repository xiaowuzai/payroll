package data

import (
	"fmt"
	"github.com/xiaowuzai/payroll/internal/config"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/pkg/uuid"
	"xorm.io/xorm"
)

func NewDB(conf *config.Database, logger *logger.Logger) (*xorm.Engine, error) {
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.Username, conf.Passwd, conf.Host, conf.Port, "payroll")
	engine, err := xorm.NewEngine("mysql", source)
	if err != nil {
		logger.Error("NewDB NewEngine error: ", err.Error())
		return nil, err
	}

	engine.ShowSQL(conf.ShowSQL)

	err = engine.Ping()
	if err != nil {
		logger.Error("NewDB Engine Ping error: ", err.Error())
		return nil, err
	}

	err = engine.Sync2(
		new(Bank),
		new(Employee),
		new(Menu),
		new(RoleMenu),
		new(Organization),
		new(Role),
		new(User),
		new(PayrollInfo),
		new(PayrollAlias),
	)
	if err != nil {
		logger.Error("NewDB Sync table error: ", err.Error())
		return nil, err
	}

	initData := &InitData{}
	err = initData.initSQL(engine)
	if err != nil {
		logger.Error("NewDB initSQL error: ", err.Error())
		return nil, err
	}

	logger.Infof("NewDB success")
	return engine,nil
}

var menus = map[string]string {
	"salary-management": "工资管理",
	"salary-preview": "当月工资预览",
	"salary-approve": "当月工资审批",
	"salary-query": "当月工资统计",
	"salary-carry-forward-record": "结转记录",
	"salary-note-history": "备忘记录",
	"salary-statistics": "工资统计",
	"cost-statistics": "成本统计",
	"project-cost-statistic": "项目费用统计",
	"system-config": "系统设置",
	"user-management": "用户管理",
	"role-management": "角色管理",
	"organization": "组织机构",
	"staff-management": "员工管理",
}

type InitData struct {}

func (initData *InitData)initSQL(db *xorm.Engine) error {
	session := db.NewSession()

	has, err := initData.checkInit(session)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	if err := session.Begin(); err != nil {
		return err
	}
	defer session.Close()

	// init role
	roleId, err := initData.initRole(session)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	menuIds, err := initData.initMenus(session)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	err = initData.initRoleMenu(session, roleId, menuIds)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	err = initData.initOrganization(session)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	err = initData.initAdminUser(session, roleId)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	err = initData.initBank(session)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	return session.Commit()
}

// 检查是否已经初始化过
func (initData *InitData)checkInit(session *xorm.Session)(bool, error){
	r := &Role{}
	has, err := session.ID("4fba1999-a7f9-4b34-b82a-b3f34cfe4d81").Get(r)
	if err != nil {
		return false, err
	}
	return has, nil
}

func (initData *InitData)initRole (session *xorm.Session) (string, error){
	var name = "系统管理员"

	role := &Role{
		Id:          uuid.CreateUUID(),
		Name:        name,
		Description: name,
	}
	_, err := session.InsertOne(role)
	if err != nil {
		return "", err
	}
	return role.Id, nil
}

func (initData *InitData)initMenus(session *xorm.Session) ([]string, error){
	menuIds := make([]string, 0, len(menus))
	ms := make([]*Menu, 0,  len(menus))

	for name, keys :=  range menus {
		id := uuid.CreateUUID()
		ms = append(ms, &Menu{
			Id: id,
			Name:name,
			Keys: keys,
		})
		menuIds = append(menuIds, id)
	}
	_, err := session.Insert(menus)
	if err != nil {
		return nil,err
	}
	return menuIds,nil
}

func (initData *InitData)initRoleMenu(session *xorm.Session, roleId string, menuIds []string) error {
	roleMenus := make([]*RoleMenu, 0, len(menuIds))
	for _, v := range menuIds {
		roleMenus= append(roleMenus, &RoleMenu{
			Id: uuid.CreateUUID(),
			RoleId: roleId,
			MenuId: v,
		})
	}
	_, err := session.Insert(roleMenus)
	return err
}

func (initData *InitData)initOrganization (session *xorm.Session)error {
	id :=  uuid.CreateUUID()
	organization := &Organization{
		Id:         id,
		Path:  "."+id,
		ParentId:     "root",
		Name:         "机关总部",
		FeeType:      1,                                       // 1:工资 2:福利 3: 退休    费用类型
		Type:         1,                                       // 1 单位、 2 工资表
		EmployeeType: 1,                                       // 员工类型： 1: 公务员  2:事业 3: 企业
	}
	_, err := session.Insert(organization)
	return err
}

func (initData *InitData)initAdminUser(session *xorm.Session, roleId string) error{
	password := "Admin@123"
	user := &User{
		Id:          uuid.CreateUUID(),
		Username:    "管理员",
		AccountName: "admin",
		Email:       "",
		RoleId:      roleId,
	}
	user.Password, user.Salt = addSalt(password)

	_, err := session.Insert(user)
	return err
}

var banks = []string {
	"江苏银行",
	"南京银行",
	"农业银行",
}
func (*InitData)initBank(session *xorm.Session) error{
	bs := make([]*Bank, 0, len(banks))
	for _, name := range banks {
		bs = append(bs, &Bank{
			Id: uuid.CreateUUID(),
			Name: name,
		})
	}

	_, err := session.Insert(bs)
	return err
}

