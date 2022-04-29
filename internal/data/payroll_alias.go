package data

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"xorm.io/xorm"
)

type PayrollAlias struct {
	Id string `xorm:"varchar(36) pk"`
	OrganizationId string `xorm:"varchar(36) notnull"`
	AliasName string `xorm:"varchar(150) notnull unique"`
}

func (pa *PayrollAlias)insertList(ctx context.Context, session *xorm.Session,logger *logger.Logger, alias []*PayrollAlias) error {
	log := logger.WithRequestId(ctx)
	log.Infof("PayrollAlias insertList %+v", alias)

	_, err := session.Insert(alias)
	if err != nil {
		log.Error("PayrollAlias insertList error: ", err.Error())
		return errors.ErrDataInsert(err.Error())
	}

	return nil
}

func (pa *PayrollAlias) listAll(ctx context.Context, session *xorm.Session,logger *logger.Logger) ([]*PayrollAlias, error) {
	log := logger.WithRequestId(ctx)
	log.Info("PayrollAlias listAll")

	alias := make([]*PayrollAlias, 0)
	err := session.Find(alias)
	if err != nil {
		log.Error("PayrollAlias listAll error: ", err.Error())
		return nil, errors.ErrDataGet(err.Error())
	}

	return alias, nil
}

func (pa *PayrollAlias) getByName(ctx context.Context, session *xorm.Session,logger *logger.Logger) (bool, error) {
	log := logger.WithRequestId(ctx)
	log.Info("PayrollAlias listAll")

	has, err := session.Get(pa)
	if err != nil {
		log.Error("PayrollAlias listAll error: ", err.Error())
		return false, errors.ErrDataGet(err.Error())
	}

	return has, nil
}

func (pa *PayrollAlias) listByOrgId(ctx context.Context, session *xorm.Session,logger *logger.Logger) ([]*PayrollAlias, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("PayrollAlias listByOrgId input %+v", *pa)

	alias := make([]*PayrollAlias, 0)
	err := session.Where("organization_id = ?", pa.OrganizationId).Find(&alias)
	if err != nil {
		log.Error("PayrollAlias listByOrgId  error: ", err.Error())
		return nil, errors.ErrDataGet(err.Error())
	}

	return alias, nil
}

func (pa *PayrollAlias) deleteByOrgId(ctx context.Context, session *xorm.Session,logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("PayrollAlias deleteByOrgId input %+v", *pa)

	alias := &PayrollAlias{}
	_, err := session.Where("organization_id = ?", pa.OrganizationId).Delete(alias)
	if err != nil {
		log.Error("PayrollAlias listByOrgId  error: ", err.Error())
		return  errors.ErrDataDelete(err.Error())
	}

	return nil
}
