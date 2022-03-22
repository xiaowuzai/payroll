package data

import "time"

type Employee struct {
	Id string `xorm:"id"`
	Number int `xorm:"number"`  // 编号
	IdCard string `xorm:""` // 身份证号
	Telephone string `xorm:""`
	OfferTime time.Time // 入职日期
	RetireTime time.Time // 退休日期
	Duty string // 职务
	Post string // 岗位
	Level string // 级别
	BaseSalary int32   // 基本工资
	Identity int32 `xorm:""`  // 身份类型： 0:公务员、 1: 事业、2: 企业
}