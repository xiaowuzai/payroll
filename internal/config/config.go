package config

import (
	"github.com/spf13/viper"
	"os"
)

type Server struct {
	Payroll *Payroll  `mapstructure:"payroll"`
	Database *Database `mapstructure:"database"`
}

type Payroll struct {
	Host     string `mapstructure:"host"`
	Port     int `mapstructure:"port"`
	Name     string `mapstructure:"name"`
}

type Database struct {
	Host     string  `mapstructure:"host"`
	Username string  `mapstructure:"username"`
	Passwd   string  `mapstructure:"passwd"`
	Port     int    `mapstructure:"port"`
	ShowSQL  bool    `mapstructure:"showSQL"`
}

func Parse() (*Server, error){
	//vipAll := viper.New()

	vip := viper.New()
	vip.SetConfigName("default")

	vip.SetConfigType("yaml")
	vip.AddConfigPath("./configs/")
	vip.AddConfigPath("./configs")
	vip.AddConfigPath(".")
	err := vip.ReadInConfig()
	if err != nil {
		panic(err)
	}
	defaults := vip.AllSettings()
	for k, v :=  range defaults{
		vip.SetDefault(k, v)
	}

	env := os.Getenv("CONF_LEVEL")
	if env != "" && env != "local"{
		vip.SetConfigName(env)
		err = vip.ReadInConfig()
		if err != nil {
			panic(err)
		}

		envConf := vip.AllSettings()
		for k, v :=  range envConf{
			vip.SetDefault(k, v)
		}
	}

	server := &Server{}
	err = vip.Unmarshal(server)
	if err != nil {
		panic(err)
	}
	return server, nil
}