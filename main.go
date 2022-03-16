package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func main1() {
	f, err := excelize.OpenFile("/Users/zly/Desktop/payroll/data.xlsx")
	if err != nil {
		panic(err)
	}
	defer func(){
		if err:= f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	//for _, v := range f.GetSheetList() {
		rows, err := f.GetRows("机关南京银行")
		if err != nil {
			panic(err)
		}
		for _, row := range rows {
	  		for _, colCell := range row {
	  			fmt.Print(colCell, "\t")
	  		}
			fmt.Println()
	   }
	//}

}


