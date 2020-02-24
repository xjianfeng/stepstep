package models

import (
	"github.com/satori/go.uuid"
	"github.com/xjianfeng/gocomm/db/mongo"
	"github.com/xjianfeng/gocomm/decry"
	"gopkg.in/mgo.v2/bson"
	"stepstep/conf"
	"time"
)

/*评论 留言*/

func AddMsg(unionId string, u *User, msg string) (string, error) {
	db := mongo.GetMongo(conf.CfgMongo.MongoDb, "message")
	defer db.Close()
	c := db.Collection

	uuv4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	msgId := decry.Md5Sum(uuv4.Bytes())
	obj := bson.M{
		"_id":        msgId,
		"unionid":    unionId,
		"avatarurl":  u.AvatarUrl,
		"nickname":   u.NickName,
		"uid":        u.Uid,
		"msg":        msg,
		"createtime": time.Now().Unix(),
	}
	err = c.Insert(obj)
	return msgId, err
}

func AddReply(msgId string, u *User, msg string) {
	db := mongo.GetMongo(conf.CfgMongo.MongoDb, "message")
	defer db.Close()
	c := db.Collection

	c.UpdateId(msgId, bson.M{"$set": bson.M{"reply": bson.M{"avatarurl": u.AvatarUrl, "nickname": u.NickName, "uid": u.Uid, "msg": msg}}})
}

type ReplyInfo struct {
	AvatarUrl string `bson:"avatarurl" json:"avatarUrl"`
	NickName  string `bson:"nickname" json"nickName"`
	Msg       string `bson:"msg" json:"msg"`
	Uid       int    `bson:"uid"  json:"uid"`
}

type MsgInfo struct {
	MsgId     string     `bson:"_id" json:"msgId"`
	AvatarUrl string     `bson:"avatarurl" json:"avatarUrl"`
	NickName  string     `bson:"nickname" json"nickName"`
	Msg       string     `bson:"msg" json:"msg"`
	Uid       int        `bson:"uid"  json:"uid"`
	Reply     *ReplyInfo `bson:"reply" json:"reply"`
}

type RetMsg struct {
	MsgCnt  int
	MsgList []MsgInfo
}

func GetMsgNum(selter bson.M) (int, error) {
	db := mongo.GetMongo(conf.CfgMongo.MongoDb, "message")
	defer db.Close()
	c := db.Collection

	n, err := c.Find(selter).Count()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func GetMsgList(unionId string, start, psize int) (*RetMsg, error) {
	db := mongo.GetMongo(conf.CfgMongo.MongoDb, "message")
	defer db.Close()
	c := db.Collection

	result := &RetMsg{
		MsgList: []MsgInfo{},
	}
	timestamp := time.Now().Unix() - 30*86400
	selobj := bson.M{"_id": 1, "avatarurl": 1, "reply": 1, "nickname": 1, "msg": 1, "uid": 1}
	filter := bson.M{"unionid": unionId, "createtime": bson.M{"$gt": timestamp}}
	n, err := c.Find(filter).Count()
	if err != nil {
		return result, err
	}
	result.MsgCnt = n
	err = c.Find(filter).Select(selobj).Sort("-createtime").Skip(start).Limit(psize).All(&result.MsgList)
	if err != nil {
		return result, err
	}
	return result, nil
}
