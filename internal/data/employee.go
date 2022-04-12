package data

import "time"

type Employee struct {
	Id string `xorm:"id varchar(36) notnull"`
	Number int `xorm:"number int notnull"`  // 编号
	IdCard string `xorm:"id_card varchar(18) notnull "` // 身份证号
	Telephone string `xorm:"telephone varchar(11)"`
	OfferTime time.Time `xorm:"offer_time"`// 入职日期
	RetireTime time.Time `xorm:"retire_time"`// 退休日期
	Duty string `xorm:"duty varchar(32)"`// 职务
	Post string  `xorm:"post  varchar(32)"`// 岗位
	Level string  `xorm:"level  varchar(32)"`// 级别
	BaseSalary int32   `xorm:"base_salary int"` // 基本工资
	Identity int32 `xorm:"identity"`  // 身份类型： 0:公务员、 1: 事业、2: 企业
}
