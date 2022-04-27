package service

import (
	"context"
	"time"
)

type EmployeeRepo interface {
	AddEmployee(context.Context, *Employee) error
	GetEmployeeList(context.Context, string, string) ([]*Employee, error)
	GetEmployee(context.Context, string) (*Employee, error)
	UpdateEmployee(context.Context, *Employee) error
	DeleteEmployee(context.Context, string) error
}

type EmployeeService struct {
	repo EmployeeRepo
}

type Employee struct {
	Id           string
	Name 		string  // 姓名
	IdCard       string  // 身份证号码
	Telephone    string // 手机号码
	Duty         string // 职务
	Post         string // 岗位
	Level        string // 级别
	OfferTime    time.Time
	RetireTime   time.Time
	Number       int   // 员工编号
	Sex int32 // 性别： 1: 男、2: 女
	Status int32 // 状态 1: 在职、2: 离职 3: 退休
	BaseSalary   int32  // 基本工资
	Identity     int32  // 身份类别： 1:公务员、 2: 事业、3: 企业
	PayrollInfos []*PayrollInfo
}

type PayrollInfo struct {
	Id             string `json:"id"`
	EmployeeId       string `json:"employeeId"`
	BankId         string `json:"bankId"`
	CardNumber     string `json:"cardNumber"`
	OrganizationId string `json:"organizationId"`
}

// 添加员工
func (es *EmployeeService) AddEmployee(ctx context.Context, em *Employee) error {
	return es.repo.AddEmployee(ctx, em)
}

// 获取员工列表
func (es *EmployeeService) ListEmployee(ctx context.Context, name string, organizationId string) ([]*Employee, error) {

	return es.repo.GetEmployeeList(ctx, name, organizationId)
}

// 获取指定员工
func (es *EmployeeService) GetEmployee(ctx context.Context, id string) (*Employee, error) {
	return es.repo.GetEmployee(ctx, id)
}

// 更新员工
func (es *EmployeeService) UpdateEmployee(ctx context.Context, e *Employee) error {
	return es.repo.UpdateEmployee(ctx, e)
}

// 删除
func (es *EmployeeService) DeleteEmployee(ctx context.Context, id string) error {
	return es.repo.DeleteEmployee(ctx, id)
}
