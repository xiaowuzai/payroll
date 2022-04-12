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

}

// 添加员工
func (es *EmployeeService)AddEmployee(ctx context.Context) {

}
