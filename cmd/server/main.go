package main

import (
	"github.com/spf13/viper"
	"github.com/xiaowuzai/payroll/internal/config"
	"log"
	"os"
)

func main() {
	vip := viper.New()
	dir, err := os.Getwd()
	panicErr(err)
	log.Println(dir)

	vip.SetConfigFile("./configs/server.yml")
	err = vip.ReadInConfig()
	panicErr(err)

	env := getEnv()
	log.Println("configDatabase: ", vip.GetString(env+"server.database.host"))
	conf := parseConf(vip, env)
	server, err := InitServer(conf, conf.Database)
	panicErr(err)

	server.Start()
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

var defaultEnv = "local."

func getEnv() string {
	env := os.Getenv("CONF_LEVEL")
	if env == "" {
		env = defaultEnv
	} else {
		env += "."
	}
	return env
}

func parseConf(vip *viper.Viper, env string) *config.Server {
	return &config.Server{
		Host: vip.GetString(env + "server.host"),
		Port: vip.GetInt(env + "server.port"),
		Name: vip.GetString(env + "server.name"),
		Database: &config.Database{
			Host:     vip.Get(env + "server.database.host").(string),
			Port:     vip.Get(env + "server.database.port").(int),
			Passwd:   vip.Get(env + "server.database.passwd").(string),
			Username: vip.Get(env + "server.database.username").(string),
			ShowSQL:  vip.Get(env + "server.database.showsql").(bool),
		},
	}
}
