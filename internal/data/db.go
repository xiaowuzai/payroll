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
	)
	if err != nil {
		logger.Error("NewDB Sync table error: ", err.Error())
		return nil, err
	}

	if err := initSQL(engine); err != nil {
		logger.Error("NewDB initSQL error: ", err.Error())
		return nil, err
	}

	logger.Infof("NewDB success")
	return engine, nil
}


func initSQL(engine *xorm.Engine)error{
	session := engine.NewSession()

	r := &Role{}
	has, err := session.ID("4fba1999-a7f9-4b34-b82a-b3f34cfe4d81").Get(r)
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
	role := &Role{
		Id: "4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",
		Name: "系统管理员"	,
		Description: "系统管理员",
	}
	_, err = session.Insert(role)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	menus := []*Menu{
		{Id:"3f946719-679d-4d0a-b265-a33e16cdd7be",Name: "salary-management", Keys:"工资管理"},
		{Id:"4b1c567e-c856-4053-afb0-f7c3094f174f",Name:"salary-preview",Keys: "当月工资预览"},
		{Id:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",Name:"salary-approve", Keys:"当月工资审批"},
		{Id:"559a1bf5-03dc-4dd3-8b79-ab816b8c30ac",Name:"salary-query", Keys:"当月工资统计"},
		{Id:"57344b90-e8ab-4ef6-8163-67c105abf6ef",Name:"salary-carry-forward-record", Keys:"结转记录"},
		{Id:"65ca014a-bf3a-42a5-8e8a-00fd52b22699",Name:"salary-note-history", Keys:"备忘记录"},
		{Id:"67b18bdb-dadf-4a0d-b6fb-95cfe360e21a",Name:"salary-statistics", Keys:"工资统计"},
		{Id:"6d2b6450-6cd3-4956-ab55-c18c400bb036",Name:"cost-statistics", Keys:"成本统计"},
		{Id:"7497609c-8fba-4c33-8bdd-bbc443865b73",Name:"project-cost-statistic", Keys:"项目费用统计"},
		{Id:"7fad0844-4563-408c-a269-e3a2d401ae88",Name:"system-config", Keys:"系统设置"},
		{Id:"8668919e-1112-49cc-8ae0-74808cb9e9cc",Name:"user-management", Keys:"用户管理"},
		{Id:"8d7fe0fb-4b98-450b-8e48-aac8a5809370",Name:"role-management", Keys:"角色管理"},
		{Id:"76b18bdb-dadf-4a0d-b6fb-95cfe360e21a",Name:"organization",Keys: "组织机构"},
		{Id:"6d7fe0fb-4b98-450b-8e48-aac8a5809370",Name:"staff-management", Keys:"员工管理"},
	}
	_, err = session.Insert(menus)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	roleMenu := []*RoleMenu{
		{Id:"de504569-42f5-4a90-9fc4-2481f05cdd89",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"3f946719-679d-4d0a-b265-a33e16cdd7be"},
		{Id:"73853828-d922-4c6a-b64c-6acc03df6831",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"4b1c567e-c856-4053-afb0-f7c3094f174f"},
		{Id:"2a9ccd2b-192a-4865-a3e0-b50656441305",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81"},
		{Id:"4e821e0b-297c-4578-b1ee-379a84091192",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"559a1bf5-03dc-4dd3-8b79-ab816b8c30ac"},
		{Id:"e7e255c3-e608-49f8-9db7-d66b72e1382c",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"57344b90-e8ab-4ef6-8163-67c105abf6ef"},
		{Id:"e4af79ee-2ae8-49c1-8a1f-e8e4c7b5e2e3",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"65ca014a-bf3a-42a5-8e8a-00fd52b22699"},
		{Id:"e3488d80-b05c-442d-8feb-94c2eadebc46",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"67b18bdb-dadf-4a0d-b6fb-95cfe360e21a"},
		{Id:"d86b6b91-60bf-42b6-a846-b1831baa0344",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"6d2b6450-6cd3-4956-ab55-c18c400bb036"},
		{Id:"7ccee999-37ca-49b9-881b-a6c95809a6a1",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"7497609c-8fba-4c33-8bdd-bbc443865b73"},
		{Id:"1b73b515-16e0-4841-8397-5442552ec708",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"7fad0844-4563-408c-a269-e3a2d401ae88"},
		{Id:"26566951-6b8d-4bed-a3d2-d12f29c90fa2",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"8668919e-1112-49cc-8ae0-74808cb9e9cc"},
		{Id:"16f5f7e7-e021-406c-9a8f-aeb2f136d518",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"8d7fe0fb-4b98-450b-8e48-aac8a5809370"},
		{Id:"40c7b944-8c0a-421a-b3f4-b52c996b5524",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"76b18bdb-dadf-4a0d-b6fb-95cfe360e21a"},
		{Id:"40c7b944-8c0a-421a-b3f4-b52c996b1234",RoleId:"4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",MenuId:"6d7fe0fb-4b98-450b-8e48-aac8a5809370"},
	}
	_, err = session.Insert(roleMenu)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	organization := &Organization{
		Id      : "16f5f7e7-e021-406c-9a8f-aeb2f136d518",
		ParentId    : "",
		Name         : "机关总部",
		Path        : ".16f5f7e7-e021-406c-9a8f-aeb2f136d518", // .ParentId.Id
		SalaryType  :"",   //   工资类型：手动输入
		FeeType     : 1,      // 1:工资 2:福利 3: 退休    费用类型
		Type        : 1,          // 1 单位、 2 工资表
		EmployeeType : 1, // 员工类型： 1: 公务员  2:事业 3: 企业
	}
	_, err = session.Insert(organization)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	password := "Admin@123"
	user := &User{
		Id: uuid.CreateUUID(),
		Username:    "管理员",
		AccountName: "admin",
		Email:       "",
		RoleId:      "4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",
	}
	user.Password, user.Salt = addSalt(password)
	_, err = session.Insert(user)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_ = session.Commit()
	return nil
}
