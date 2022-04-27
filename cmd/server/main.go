package main

import (
	"github.com/xiaowuzai/payroll/internal/config"
)

func main() {
	conf, err := config.Parse()
	panicErr(err)

	server, err := InitServer(conf.Payroll, conf.Database)
	panicErr(err)

	server.Start()
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

