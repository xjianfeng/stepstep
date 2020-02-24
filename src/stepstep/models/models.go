package models

import (
	"github.com/json-iterator/go"
	"github.com/xjianfeng/gocomm/db/mongo"
	"github.com/xjianfeng/gocomm/logger"
	"gopkg.in/mgo.v2/bson"
	"stepstep/conf"
	"stepstep/define"
	"stepstep/package/appdata"
	"stepstep/package/utils"
	"strconv"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func getCacheUid2UnionId(uid int) string {
	key := strconv.Itoa(uid)
	data := utils.Uid2UnionIdLru.GetNode(key)
	if data != nil {
		return data.(string)
	}
	return ""
}

func setCacheUid2UnionId(uid int, unionId string) {
	key := strconv.Itoa(uid)
	utils.Uid2UnionIdLru.AddNode(key, unionId)
}

func getCacheUser(unionId string) *User {
	data := utils.Userlru.GetNode(unionId)
	if data != nil {
		return data.(*User)
	}
	return nil
	// 暂时不用redis
	ret := GetRedisUInfo(unionId, define.REDIS_HEKY_USER_CACHE)

	if ret != nil {
		u := &User{}
		json.Unmarshal(ret, u)
		return u
	}
	return nil
}

func setCacheUser(u *User) error {
	utils.Userlru.AddNode(u.UnionId, u)
	return nil
	// 暂时不用redis
	value, err := json.Marshal(u)
	if err != nil {
		return err
	}
	SetRedisUInfo(u.UnionId, define.REDIS_HEKY_USER_CACHE, string(value))
	return nil
}

func SaveUser(u *User) {
	if !u.SaveDbFlg {
		return
	}
	u.SaveDbFlg = false
	m := mongo.GetMongo(conf.CfgMongo.MongoDb, "user")
	defer m.Close()
	c := m.Collection

	unionId := u.UnionId
	_, err := c.UpsertId(unionId, u)
	if err != nil {
		logger.LogError("SaveUser Err:%s", err.Error())
		return
	}
	logger.LogInfo("SaveUser unionId:%s", u.UnionId)
	setCacheUser(u)
}

type RetDbUnionId struct {
	UnionId string `bson:"_id"`
}

func LoadUnionIdByUid(uid int) (string, error) {
	unionId := getCacheUid2UnionId(uid)
	if unionId != "" {
		logger.LogInfo("LoadUnionId For Lru Cache %d", uid)
		return unionId, nil
	}
	m := mongo.GetMongo(conf.CfgMongo.MongoDb, "user")
	defer m.Close()
	c := m.Collection

	data := &RetDbUnionId{}
	err := c.Find(bson.M{"uid": uid}).Select(bson.M{"_id": 1}).One(data)
	if err != nil {
		logger.LogError("LoadUnionIdByUid err:%s", err.Error())
		return "", err
	}
	setCacheUid2UnionId(uid, data.UnionId)
	return data.UnionId, nil
}

func LoadUser(unionId string) (*User, error) {
	if unionId == "" {
		return nil, define.ERROR_REQUEST_PARAMS
	}
	u := getCacheUser(unionId)
	if u != nil {
		logger.LogInfo("LoadUser For Lru Cache %s", u.UnionId)
		return u, nil
	}
	m := mongo.GetMongo(conf.CfgMongo.MongoDb, "user")
	defer m.Close()
	c := m.Collection

	user := &User{}
	err := c.FindId(unionId).One(user)

	if err != nil {
		logger.LogError("LoadUser:%s", err.Error())
		return nil, err
	}
	setCacheUser(user)
	return user, nil
}

func initAwardList(u *User) {
	u.AwardCnt += 1
	for i, step := range define.AWARD_STEP_INTERVAL {
		//根据领取奖励的天数来获取对应第几天的奖励来初始化数据
		awardMoney, ok := appdata.GetAwardValue(u.AwardCnt, i+1)
		if !ok {
			continue
		}
		u.ResetData.AwardList[i].AwardStep = step
		u.ResetData.AwardList[i].AwardMoney = awardMoney
	}
}

func Newuser(unionId, openId string) (*User, error) {
	var err error
	u := new(User)
	u.Uid, err = GenIncrUid()
	if err != nil {
		return nil, err
	}
	u.WxData = new(WxData)
	u.OpenId = openId
	u.UnionId = unionId
	u.RealDayNo = utils.GetRealDayNo()
	u.CreateTime = time.Now().Unix()
	u.SaveDbFlg = true
	u.ResetData = &DayData{
		Friend: &FriendData{
			FriendHeadIcon: []string{},
			FriendList:     make(map[string]string),
		},
	}
	initTimeOutData(u)
	initAwardList(u)
	return u, nil
}

func initTimeOutData(u *User) {
	u.Timeout = make(map[string]*TimeData)

	now := time.Now().Unix()
	//每两小时奖励1000步
	u.Timeout[define.TIMEOUT_KEY_LUCKSTEP] = &TimeData{
		TimeOut:    now + define.TIMOUT_AWARD_STEP,
		MaxRefresh: define.TIMEOUT_NOT_LIMIT_REFRESH,
		Info: map[string]interface{}{
			define.AWARD_TYPE_STEP: 1000,
		},
	}

	//每30秒看一次视频奖励1000步
	u.Timeout[define.TIMEOUT_KEY_VIDEO] = &TimeData{
		TimeOut:    now + define.TIMOUT_READ_VIDEO,
		MaxRefresh: define.AD_VIDEO_MAX_CNT,
		Info: map[string]interface{}{
			define.AWARD_TYPE_STEP: 1000,
		},
	}

	//每两小时随机一个红包
	u.Timeout[define.TIMEOUT_KEY_REDPACK] = &TimeData{
		TimeOut:    now + define.TIMOUT_RED_REDPACK,
		MaxRefresh: define.TIMEOUT_NOT_LIMIT_REFRESH,
		Info: map[string]interface{}{
			define.AWARD_TYPE_MONEY: []int{1, 4},
		},
	}
}
