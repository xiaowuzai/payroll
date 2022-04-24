package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/config"
	"github.com/xiaowuzai/payroll/internal/router"
)

type Server struct {
	addr      string
	engine    *gin.Engine
	apiRouter *router.Router
}

func (s *Server) Start() {
	err := s.engine.Run(s.addr)
	if err != nil {
		panic(err)
	}
}

func NewServer(engine *gin.Engine, apiRouter *router.Router, conf *config.Server) *Server {
	apiRouter.WithEngine(engine)
	return &Server{
		addr:      fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		engine:    engine,
		apiRouter: apiRouter,
	}
}
