package main

import (
	"github.com/spf13/viper"
"github.com/xiaowuzai/payroll/internal/data"
)

func main(){
	viper.ReadInConfig()
	data.NewData()
}


