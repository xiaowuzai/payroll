package idcard

import (
	"regexp"
	"strconv"
	"time"
)

// 检查地区
func checkProv (val string) bool {
	var pattern = "/^[1-9][0-9]/"
	var provs = map[int]string{11:"北京",12:"天津",13:"河北",14:"山西",15:"内蒙古",21:"辽宁",22:"吉林",23:"黑龙江 ",31:"上海",32:"江苏",33:"浙江",34:"安徽",35:"福建",36:"江西",37:"山东",41:"河南",42:"湖北 ",43:"湖南",44:"广东",45:"广西",46:"海南",50:"重庆",51:"四川",52:"贵州",53:"云南",54:"西藏 ",61:"陕西",62:"甘肃",63:"青海",64:"宁夏",65:"新疆",71:"中国台湾",81:"中国香港",82:"中国澳门"}

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	has := reg.MatchString(val)
	if !has {
		return has
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return false
	}
	 _, has = provs[i]
	 return has
}

// 检查出生日期
func checkDate (val string) bool{
	var pattern = "/^(18|19|20)\\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)$/"
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	has := reg.MatchString(val)
	if !has {
		return false
	}

	year, err  := strconv.Atoi(val[:4])
	if err != nil {
		return false
	}
	month, err := strconv.Atoi(val[4:6])
	if err != nil {
		return false
	}
	day, err := strconv.Atoi(val[6:8])
	if err != nil {
		return false
	}

	t := time.Date(year, time.Month(month), day,0,0,0,0,time.Local)
	if int(t.Month()) != month {
		return false
	}

	return true
}

// 检查身份校验码
func checkCode (val string) {
var p = "/^[1-9]\\d{5}(18|19|20)\\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$/"
var factor = []int{ 7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2 }
var parity = []int{ 1, 0, 'X', 9, 8, 7, 6, 5, 4, 3, 2 }

var code = val[:17]
if(p.test(val)) {
var sum = 0;
for(var i=0;i<17;i++) {
sum += val[i]*factor[i];
}
if(parity[sum % 11] == code.toUpperCase()) {
return true;
}
}
return false;
}