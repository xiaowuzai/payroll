package data

import (
	"context"
	"fmt"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/service"
	"github.com/xiaowuzai/payroll/pkg/uuid"
	"time"
	"xorm.io/xorm"
)

var _ service.EmployeeRepo = (*EmployeeRepo)(nil)

type EmployeeRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewEmployeeRepo(data *Data, logger *logger.Logger) service.EmployeeRepo {
	return &EmployeeRepo{
		data:   data,
		logger: logger,
	}
}

func (er *EmployeeRepo) AddEmployee(ctx context.Context, se *service.Employee) error {
	session, err := BeginSession(ctx, er.data.db, er.logger)
	if err != nil {
		return err
	}
	defer session.Close()

	employee := &Employee{}
	employee.fromService(se)
	employee.Id = uuid.CreateUUID()
	err = employee.insert(ctx, session, er.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// 添加银行卡信息
	payInfos := make([]*PayrollInfo, 0, len(se.PayrollInfos))
	for _, spi := range se.PayrollInfos {

		payInfo := &PayrollInfo{}
		payInfo.fromService(spi)
		payInfo.Id = uuid.CreateUUID()
		payInfos = append(payInfos, payInfo)
	}

	payrollInfo := &PayrollInfo{}
	err = payrollInfo.insertList(ctx, session, er.logger, payInfos)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_ = session.Commit()
	return nil
}

func (er *EmployeeRepo) GetEmployeeList(ctx context.Context, name, organizationId string) ([]*service.Employee, error) {
	session := NewSession(ctx, er.data.db)

	employee := &Employee{}
	emps, err := employee.list(ctx, session, er.logger, name, organizationId)
	if err != nil {
		return nil, err
	}

	ses := make([]*service.Employee, 0, len(emps))
	for _, v := range emps {
		se := v.toService()
		ses = append(ses, se)
	}

	return ses, nil
}

func (er *EmployeeRepo) GetEmployee(ctx context.Context, id string) (*service.Employee, error) {
	session := NewSession(ctx, er.data.db)
	log := er.logger.WithRequestId(ctx)
	log.Infof("GetEmployee id = %s\n", id)

	employee := &Employee{}
	has, err := employee.get(ctx, session, er.logger)
	if err != nil {
		return nil, err
	}
	if !has {
		message := fmt.Sprintf("employee id = %s not exist", id)
		log.Error(message)
		return nil, errors.DataNotFound(message)
	}

	return employee.toService(), nil
}

func (er *EmployeeRepo) UpdateEmployee(ctx context.Context, se *service.Employee) error {
	session, err := BeginSession(ctx, er.data.db, er.logger)
	if err != nil {
		return err
	}
	defer session.Close()

	// 更新员工
	emp := &Employee{}
	emp.fromService(se)
	err = emp.update(ctx, session, er.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// 删除 payrollInfo
	pi := &PayrollInfo{}
	err = pi.deleteByEmployeeId(ctx, session, er.logger, se.Id)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// 插入 payrollInfo
	pis := make([]*PayrollInfo, 0, len(se.PayrollInfos))
	for _, v := range se.PayrollInfos {
		pi := &PayrollInfo{}
		pi.fromService(v)
		pi.Id = uuid.CreateUUID()

		pis = append(pis, pi)
	}
	err = pi.insertList(ctx, session, er.logger, pis)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_ = session.Commit()
	return nil
}

func (er *EmployeeRepo) DeleteEmployee(ctx context.Context, id string) error {
	session, err := BeginSession(ctx, er.data.db, er.logger)
	if err != nil {
		return err
	}
	defer session.Close()

	// 删员工
	em := &Employee{Id: id}
	err = em.delete(ctx, session, er.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// 删银行卡信息
	pi := &PayrollInfo{}
	err = pi.deleteByEmployeeId(ctx, session, er.logger, id)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_ = session.Commit()
	return nil
}

type Employee struct {
	Id         string    `xorm:"id varchar(36) notnull"`
	IdCard     string    `xorm:"id_card varchar(18) notnull "` // 身份证号
	Telephone  string    `xorm:"telephone varchar(11)"`
	OfferTime  time.Time `xorm:"offer_time"`         // 入职日期
	RetireTime time.Time `xorm:"retire_time"`        // 退休日期
	Duty       string    `xorm:"duty varchar(32)"`   // 职务
	Post       string    `xorm:"post  varchar(32)"`  // 岗位
	Level      string    `xorm:"level  varchar(32)"` // 级别
	Number     int       `xorm:"number int notnull"` // 编号
	BaseSalary int32     `xorm:"base_salary"`        // 基本工资
	Identity   int32     `xorm:"identity"`           // 身份类型： 0:公务员、 1: 事业、2: 企业
	Sex        int32     `xorm:"sex"`
	Status     int32     `xorm:"status"`
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
		Sex:        e.Sex,
		Status:     e.Status,
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
	e.Status = s.Status
	e.Sex = s.Sex
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

func (e *Employee) list(ctx context.Context, session *xorm.Session, logger *logger.Logger, name string, organizationId string) ([]*Employee, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("Employee list name = %s and organizationId = %s\n", name, organizationId)

	es := make([]*Employee, 0)
	if name != "" {
		session = session.Where("name = ?", name)
	}
	if organizationId != "" {
		session = session.Where("organizationId = ?", organizationId)
	}
	err := session.Find(&es)
	if err != nil {
		log.Error("Employee list error: ", err.Error())
		return nil, errors.ErrDataGet(err.Error())
	}

	return es, nil
}
