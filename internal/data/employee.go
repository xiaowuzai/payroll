package data

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/service"
	"time"
	"xorm.io/xorm"
)

type Employee struct {
	Id         string    `xorm:"id varchar(36) notnull"`
	Number     int       `xorm:"number int notnull"`           // 编号
	IdCard     string    `xorm:"id_card varchar(18) notnull "` // 身份证号
	Telephone  string    `xorm:"telephone varchar(11)"`
	OfferTime  time.Time `xorm:"offer_time"`         // 入职日期
	RetireTime time.Time `xorm:"retire_time"`        // 退休日期
	Duty       string    `xorm:"duty varchar(32)"`   // 职务
	Post       string    `xorm:"post  varchar(32)"`  // 岗位
	Level      string    `xorm:"level  varchar(32)"` // 级别
	BaseSalary int32     `xorm:"base_salary int"`    // 基本工资
	Identity   int32     `xorm:"identity"`           // 身份类型： 0:公务员、 1: 事业、2: 企业
}

func (e *Employee) toService() *service.Employee {
	return &service.Employee{
		Id:         e.Id,
		Number:     e.Number,
		IdCard:     e.IdCard,
		Telephone:  e.Telephone,
		OfferTime:  e.OfferTime,
		RetireTime: e.RetireTime,
		Duty:       e.Duty,
		Post:       e.Post,
		Level:      e.Level,
		BaseSalary: e.BaseSalary,
		Identity:   e.Identity,
		//PayrollInfos []*PayrollInfo
	}
}

func (e *Employee) fromService(s *service.Employee) {
	e.Id = s.Id
	e.Number = s.Number
	e.IdCard = s.IdCard
	e.Telephone = s.Telephone
	e.OfferTime = s.OfferTime
	e.RetireTime = s.RetireTime
	e.Duty = s.Duty
	e.Post = s.Post
	e.Level = s.Level
	e.BaseSalary = s.BaseSalary
	e.Identity = s.Identity
}

func (e *Employee) get(ctx context.Context, session *xorm.Session, logger *logger.Logger) (bool, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("Employee get input %+v\n", *e)

	has, err := session.Get(e)
	if err != nil {
		log.Errorf("Employee get error: %s\n", err.Error())
		return false, errors.ErrDataGet(err.Error())
	}

	return has, nil
}

func (e *Employee) insert(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("Employee insert input %+v\n", *e)

	_, err := session.Insert(e)
	if err != nil {
		log.Errorf("Employee insert error: %s\n", err.Error())
		return errors.ErrDataInsert(err.Error())
	}
	return nil
}

func (e *Employee) update(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("Employee update input %+v\n", *e)

	_, err := session.Where("id = ?", e.Id).Update(e)
	if err != nil {
		log.Error("Employee update error: ", err.Error())
		return errors.ErrDataUpdate(err.Error())
	}
	return nil
}

func (e *Employee) delete(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("Employee delete input %+v\n", *e)

	em := &Employee{}
	_, err := session.Where("id = ?", e.Id).Delete(em)
	if err != nil {
		log.Error("Employee delete error: ", err.Error())
		return errors.ErrDataDelete(err.Error())
	}
	return nil
}

func (e *Employee) list(ctx context.Context, session *xorm.Session, logger *logger.Logger) ([]*Employee, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("Employee list \n")

	es := make([]*Employee, 0)
	err := session.Find(&es)
	if err != nil {
		log.Error("Employee list error: ", err.Error())
		return nil, errors.ErrDataGet(err.Error())
	}

	return es, nil
}
