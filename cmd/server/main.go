package main

import (
	"github.com/spf13/viper"
	"github.com/xiaowuzai/payroll/internal/config"
	"log"
	"os"
)

func main(){
	vip := viper.New()
	dir, err := os.Getwd()
	panicErr(err)
	log.Println(dir)

	vip.SetConfigFile("../../configs/server.yml")
	err = vip.ReadInConfig()
	panicErr(err)

	log.Println("configDatabase: ", vip.GetString("server.database.host"))
	conf := parseConf(vip)
	//db, err := data.NewData(conf.Database)
	//panicErr(err)
	//
	//ge := gin.New()
	//r := router.NewRouter()
	//server := app.NewServer(ge,r)
	server, err := InitServer(conf, conf.Database)
	panicErr(err)
	server.Start()

}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func parseConf(vip *viper.Viper) *config.Server {
	return &config.Server{
		Host:vip.GetString("server.host"),
		Port:vip.GetInt("server.port"),
		Name:vip.GetString("server.name"),
		Database: &config.Database{
			Host:     vip.Get("server.database.host").(string),
			Port:     vip.Get("server.database.port").(int),
			Passwd:   vip.Get("server.database.passwd").(string),
			Username: vip.Get("server.database.username").(string),
		},
	}
}
