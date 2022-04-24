package data

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"xorm.io/xorm"
)

var ErrBeginSession = errors.New("数据库异常")

func BeginSession(ctx context.Context, db *xorm.Engine, logger *logger.Logger) (*xorm.Session, error) {
	session := NewSession(ctx, db)
	return Begin(ctx, session, logger)

}

func NewSession(ctx context.Context, db *xorm.Engine) *xorm.Session {
	return db.NewSession()
}

func Begin(ctx context.Context, session *xorm.Session, logger *logger.Logger) (*xorm.Session, error) {
	log := logger.WithRequestId(ctx)
	if err := session.Begin(); err != nil {
		log.Errorf("开启 Session 异常")
		return nil, ErrBeginSession
	}
	return session, nil
}
