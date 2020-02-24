package service

import (
	"github.com/xjianfeng/gocomm/db/dbredis"
	"github.com/xjianfeng/gocomm/db/mongo"
	"github.com/xjianfeng/gocomm/logger"
	"stepstep/conf"
	"testing"
)

func SetUp() {
	//初始化项目对应配置
	mongo.SetUp(conf.CfgMongo.MongoUri)
	logger.SetUp(
		conf.CfgServer.LogPath,
		conf.CfgServer.AppName,
		conf.CfgServer.RunMode,
	)
	dbredis.SetUp(&dbredis.RedisConf{
		Host:        conf.CfgRedis.Host,
		Password:    conf.CfgRedis.Password,
		MaxIdle:     conf.CfgRedis.MaxIdle,
		MaxActive:   conf.CfgRedis.MaxActive,
		DefaultDb:   conf.CfgRedis.DefaultDb,
		IdleTimeout: conf.CfgRedis.IdleTimeout,
	})
}

func TestUserInfo(t *testing.T) {
	SetUp()
	u, err := GetUserData("opYS45Yyx-o-UUxiKG_M_dDLP6Lk", "opYS45Yyx-o-UUxiKG_M_dDLP6Lk")
	t.Logf("u=========== :%v, err:%v", u, err)
	u, err = GetUserData("opYS45Yyx-o-UUxiKG_M_dDLP6Lk", "opYS45Yyx-o-UUxiKG_M_dDLP6Lk")
	t.Logf("u********** :%v, err:%v", u, err)
}
