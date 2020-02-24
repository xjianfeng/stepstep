package models

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xjianfeng/gocomm/db/dbredis"
	"github.com/xjianfeng/gocomm/logger"
	"stepstep/define"
	"time"
)

func SetRedisUInfo(unionId string, hkey string, value string) {
	r := dbredis.GetRedis()
	defer r.Close()

	key := define.REDIS_KEY_USER + unionId
	r.Send("HSET", key, hkey, value)
	r.Send("EXPIRE", key, 86400)
	err := r.Flush()
	if err != nil {
		logger.LogError("SetRedisUInfo err:%s", err.Error())
	}
}

func GetRedisUInfo(unionId string, hkey string) []byte {
	r := dbredis.GetRedis()
	defer r.Close()

	key := define.REDIS_KEY_USER + unionId
	ret, err := redis.Bytes(r.Do("HGET", key, hkey))
	if err != nil {
		logger.LogError("GetRedisUInfo err:%s", err.Error())
	}
	return ret
}

func GenIncrUid() (int, error) {
	r := dbredis.GetRedis()
	defer r.Close()

	key := define.REDIS_KEY_INCR_UID
	r.Send("SETNX", key, 14578934)
	r.Send("INCR", key)
	r.Flush()

	uid, err := redis.Int(r.Do("GET", key))
	return uid, err
}

type FormIdInfo struct {
	Timestamp int64
	OpenId    string
	FormId    string
}

func AddRedisFormId(unionId, openId, formId string) {
	r := dbredis.GetRedis()
	defer r.Close()
	key := define.REDIS_KEY_FORMID + unionId
	l, err := redis.Int(r.Do("LLEN", key))
	if err != nil {
		return
	}
	if l >= 30 {
		r.Do("LPOP", key)
	}
	v, _ := json.Marshal(&FormIdInfo{
		Timestamp: time.Now().Unix(),
		OpenId:    openId,
		FormId:    formId,
	})
	r.Do("RPUSH", key, string(v))
	r.Do("EXPIRE", key, 7*86400)
}

func GetRedisFormId(unionId string) (string, string) {
	r := dbredis.GetRedis()
	defer r.Close()
	key := define.REDIS_KEY_FORMID + unionId
	v := &FormIdInfo{}

	for {
		ret, err := redis.Bytes(r.Do("LPOP", key))
		if err != nil {
			return "", ""
		}
		err = json.Unmarshal(ret, v)
		if err != nil {
			return "", ""
		}
		if time.Now().Unix()-v.Timestamp > define.TIME_WEEK_SEC {
			continue
		}
		return v.FormId, v.OpenId
	}
	return "", ""
}
