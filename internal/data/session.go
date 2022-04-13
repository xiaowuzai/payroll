package data

import (
	"context"
	"errors"
	"log"
	"xorm.io/xorm"
)

var ErrBeginSession error = errors.New("数据库异常")

func BeginSession(ctx context.Context, db *xorm.Engine)(*xorm.Session, error) {
	session := db.NewSession()
	if err := session.Begin(); err != nil {
		log.Println("开启 Session 异常")
		return nil, ErrBeginSession
	}
	return session, nil
}
