package service

import (
	"context"
	"time"
)

type EmployeeRepo interface {
	AddEmployee(context.Context, *Employee) error
	GetEmployeeList(context.Context) ([]*Employee, error)
	GetEmployee(context.Context, string) (*Employee, error)
	UpdateEmployee(context.Context,*Employee) (*Employee, error)
}

type EmployeeService struct {
	repo EmployeeRepo
}

type Employee struct {
	Id string
	Number int
	IdCard string
	Telephone string
	OfferTime time.Time
	RetireTime time.Time
	Duty string // 职务
	Post string  // 岗位
	Level string  // 级别
	BaseSalary int32   // 基本工资
	Identity int32  // 身份类型： 0:公务员、 1: 事业、2: 企业
	PayrollInfos []*PayrollInfo
}

type PayrollInfo struct {
	Id string `json:"id"`
	Employee string `json:"employeeId"`
	BankId string `json:"bankId"`
	CardNumber string `json:"cardNumber"`
	OrganizationId string `json:"organizationId"`
}

// 添加员工
func (es *EmployeeService)AddEmployee(ctx context.Context, em *Employee) error{
	return es.repo.AddEmployee(ctx, em)
}

// 获取员工列表
func (es *EmployeeService)GetEmployeeList(context.Context) ([]*Employee, error){

}

// 获取指定员工
func (es *EmployeeService)GetEmployee(context.Context, string) (*Employee, error){

}

// 更新员工
func (es *EmployeeService)UpdateEmployee(context.Context,*Employee) (*Employee, error){

}

