package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

type testInfo struct {
	A string
	S int
}

type testData struct {
	Alist [1]testInfo
}

func test() {
	for i := 0; i < 10; i++ {
		fmt.Printf("%v,", float64(rand.Intn(4)+1)/10.0)
	}
}

/**
* showdoc
* @catalog 测试接口
* @title 测试文档
* @description 步数宝测试接口
* @method post
* @url https://server.kaizhan8.com/step/test
* @param username 必选 string 用户名
* @param password 必选 string 密码
* @param name 可选 string 用户昵称
* @return {"error_code":0,"data":{"uid":"1","username":"12154545","name":"吴系挂","groupid":2,"reg_time":"1436864169","last_login_time":"0"}}
* @return_param groupid int 用户组id
* @return_param name string 用户昵称
* @remark 这里是备注信息
* @return_param groupid int 用户组id
* @return_param name string 用户昵称
* @number 99
 */
func ApiShowDoc() {
	println("test ApiShowDoc")
}

var (
	sinTimeSec, _ = time.ParseInLocation("2006-01-02 15:04:05", "2019-01-01 00:00:00", time.Local)
)

func GetRealDayNo() int {
	DayNo := int(math.Ceil(time.Now().Sub(sinTimeSec).Hours() / 24))
	return DayNo
}

func main() {
	fmt.Printf("GetRealDayNo %d", GetRealDayNo())
	fmt.Printf("Year Day %d", time.Now().YearDay())
	formatTime := time.Now().Format("2006-01-02 15:04:05")
	retTime := strings.Split(formatTime, " ")
	println("")
	println(retTime[0] + retTime[1])
}
