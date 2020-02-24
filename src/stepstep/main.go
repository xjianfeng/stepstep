package main

import (
	"github.com/xjianfeng/gocomm/db/dbredis"
	"github.com/xjianfeng/gocomm/db/mongo"
	"github.com/xjianfeng/gocomm/logger"
	"os"
	"os/signal"
	"stepstep/conf"
	"stepstep/package/cos"
	"stepstep/routers"
	"syscall"
)

func handelSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM)
	s := <-c
	logger.LogInfo("handelSignal signal %s", s.String())
	os.Exit(0)
}

func init() {
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

func main() {
	go handelSignal()
	r := routers.SetRouter()
	r.Run(conf.CfgServer.HTTPAddr)
}
