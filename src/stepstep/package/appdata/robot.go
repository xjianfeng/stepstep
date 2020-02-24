package appdata

import (
	"fmt"
	"math/rand"
	"stepstep/conf"
	"time"
)

var (
	xlsRobotFile  = "robot.json"
	robotNameList = []string{}

	cacheRobotList = CacheRobotInfo{}
)

type CacheRobotInfo struct {
	Timestamp int64
	RobotList []XlsUserData
}

type XlsUserData struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Money int    `json:"money"`
}

func GetRandomUserInfo(randCnt int) []XlsUserData {
	now := time.Now().Unix()
	if cacheRobotList.Timestamp > 0 && now-cacheRobotList.Timestamp < 180 {
		return cacheRobotList.RobotList
	}
	idx := 0
	name := ""
	image := ""
	money := 0
	l := len(robotNameList)
	cacheRobotList.Timestamp = now
	cacheRobotList.RobotList = []XlsUserData{}

	for i := 0; i < randCnt; i++ {
		idx = rand.Intn(l - 1)
		name = robotNameList[idx]
		if idx > 1000 {
			idx = idx % 1000
		}
		image = conf.CfgServer.ImageDomain + fmt.Sprintf("%d.jpg", idx)
		money = rand.Intn(100) + 50
		cacheRobotList.RobotList = append(cacheRobotList.RobotList, XlsUserData{
			Name:  name,
			Image: image,
			Money: money,
		})
	}
	return cacheRobotList.RobotList
}

func InitRobotData() {
	nameData, err := GetXlsData(xlsRobotFile, "name")
	onlyName, err := GetXlsData(xlsRobotFile, "onlyName")
	if err != nil {
		panic(err)
	}
	for _, sheetData := range nameData {
		firstName, _ := sheetData["FirstName"].(string)
		for _, vsheetData := range nameData {
			lastName, _ := vsheetData["LastName"].(string)
			name := firstName + lastName
			robotNameList = append(robotNameList, name)
		}
	}
	for _, sheetData := range onlyName {
		name, _ := sheetData["name"].(string)
		robotNameList = append(robotNameList, name)
	}
}

func init() {
	InitRobotData()
}
