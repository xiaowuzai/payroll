package data

import (
	"context"
	"errors"
	"github.com/xiaowuzai/payroll/internal/service"
	"log"
	"time"
	"xorm.io/xorm"
)


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

func (e *Employee)toService() *service.Employee{
	return &service.Employee{
		Id : e.Id,
		Number: e.Number,
		IdCard: e.IdCard,
		Telephone: e.Telephone,
		OfferTime: e.OfferTime,
		RetireTime: e.RetireTime,
		Duty: e.Duty,
		Post : e.Post,
		Level: e.Level,
		BaseSalary : e.BaseSalary,
		Identity : e.Identity,
		//PayrollInfos []*PayrollInfo
	}
}

func (e *Employee)fromService(s *service.Employee)  {
	e.Id = s.Id
	e.Number =s.Number
	e.IdCard =s.IdCard
	e.Telephone =s.Telephone
	e.OfferTime = s.OfferTime
	e.RetireTime = s.RetireTime
	e.Duty = s.Duty
	e.Post = s.Post
	e.Level= s.Level
	e.BaseSalary= s.BaseSalary
	e.Identity = s.Identity
}

func (e *Employee)get(ctx context.Context, session *xorm.Session) (bool,error){
	has, err := session.Get(e)
	if err != nil {
		log.Println("Employee get error: ", err.Error())
		return false, errors.New("获取数据出错")
	}

	return has, nil
}

func (e *Employee)insert(ctx context.Context, session *xorm.Session) error {
	_, err := session.Insert(e)
	if err != nil {
		log.Println("Employee insert error: ", err.Error())
		return errors.New("插入数据错误")
	}
	return nil
}

func (e *Employee) update(ctx context.Context, session *xorm.Session) error {
	_, err := session.Where("id = ?", e.Id).Update(e)
	if err != nil {
		log.Println("Employee update error: ", err.Error())
		return errors.New("更新数据出错")
	}
	return nil
}

func (e *Employee) delete(ctx context.Context, session *xorm.Session) error {
	em := &Employee{}
	_, err := session.Where("id = ?", e.Id).Delete(em)
	if err != nil {
		log.Println("Employee delete error: ", err.Error())
		return errors.New("删除数据出错")
	}
	return nil
}

func (e *Employee) list(ctx context.Context, session *xorm.Session) ([]*Employee, error) {
	es := make([]*Employee, 0)
	err := session.Find(&es)
	if err != nil {
		log.Println("Employee list error: ", err.Error())
		return nil, errors.New("获取数据出错")
	}
	return es, nil
}


