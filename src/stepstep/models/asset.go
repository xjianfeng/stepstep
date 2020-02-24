package models

/*资产信息*/
import (
	"github.com/xjianfeng/gocomm/db/mongo"
	"github.com/xjianfeng/gocomm/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"stepstep/conf"
	"strings"
	"time"
)

type assestLogInfo struct {
	unionId string
	modType int8
	system  string
	desc    string
	value   float64
}

type assetLogData struct {
	UnionId    string
	Date       string
	List       []map[string]interface{}
	CreateTime string
}

func addAssetLog(logData *assestLogInfo) {
	db := mongo.GetMongo(conf.CfgMongo.MongoDb, "assest")
	defer db.Close()
	c := db.Collection

	now := time.Now().Unix()
	formatTime := time.Now().Format("2006-01-02 15:04:05")
	retTime := strings.Split(formatTime, " ")
	date := retTime[0]
	htime := retTime[1]

	selobj := bson.M{
		"unionid": logData.unionId,
		"date":    date,
	}
	loginfo := bson.M{
		"modtype":    logData.modType,
		"system":     logData.system,
		"desc":       logData.desc,
		"value":      logData.value,
		"htime":      htime,
		"createtime": now,
	}
	dbData := assetLogData{}
	err := c.Find(selobj).One(&dbData)
	if err != nil && err != mgo.ErrNotFound {
		return
	}
	if err == nil {
		dbData.List = append(dbData.List, loginfo)
		c.Update(selobj, dbData)
		return
	}
	if err == mgo.ErrNotFound {
		obj := bson.M{
			"unionid":    logData.unionId,
			"date":       date,
			"createtime": now,
			"list":       []bson.M{loginfo},
		}
		c.Insert(obj)
	}
}

type AssetList struct {
	Htime   string  `json:"htime"`
	Desc    string  `json:"desc"`
	Value   float64 `json:"value"`
	ModType int8    `json:"modtype"`
}

type AssetInfo struct {
	Date string
	List []AssetList
}

type RetAssetData struct {
	RecordNum int
	AssetLog  []AssetInfo
}

func GetAssetLog(unionId string, start, psize int) *RetAssetData {
	db := mongo.GetMongo(conf.CfgMongo.MongoDb, "assest")
	defer db.Close()
	c := db.Collection

	result := &RetAssetData{
		AssetLog: []AssetInfo{},
	}

	selobj := bson.M{"date": 1, "list": 1}
	filter := bson.M{"unionid": unionId}
	n, err := c.Find(filter).Count()
	if err != nil {
		logger.LogError("GetAssetLog error:%s", err.Error())
		return result
	}
	err = c.Find(filter).Select(selobj).Sort("-createtime").Skip(start).Limit(psize).All(&result.AssetLog)
	if err != nil {
		logger.LogError("GetAssetLog error:%s", err.Error())
		return result
	}
	result.RecordNum = n
	return result
}
