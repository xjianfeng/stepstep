package units

import (
	"github.com/xjianfeng/gocomm/db/dbredis"
	"github.com/xjianfeng/gocomm/db/mongo"
	"github.com/xjianfeng/gocomm/logger"
	"stepstep/conf"
	"stepstep/package/cos"
	"stepstep/units/service"
)

func SetUp() {
	conf.SetUp("config.ini")
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
	cos.SetCosOption()
}

func UnitTest() {
	service.TestAssetDetail()
}
