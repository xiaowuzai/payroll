package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/config"
	"github.com/xiaowuzai/payroll/internal/router"
)

type Server struct {
	engine *gin.Engine
	apiRouter *router.Router
}

func (s *Server)Start(conf config.Server) {
	addr := fmt.Sprintf("%s:%d", conf.Addr, conf.Port)

	err := s.engine.Run(addr)
	if err != nil {
		panic(err)
	}
}

func NewServer(engine *gin.Engine, apiRouter *router.Router) *Server {
	return &Server{
		engine: engine,
		apiRouter: apiRouter,
	}
}
