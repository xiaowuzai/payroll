package data

import "time"

type User struct {
	Id string `xorm:"id"`
	Name string `xorm:"name"`
	Telephone string `xorm:"telephone"`
	Email string `xorm:"email"`
	Role string `xorm:"role"`
	Password string `xorm:"password"`
	Salt string `xorm:"salt"`
	Status int32 `xorm:"status"`  // 0 正常、1 禁用
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}