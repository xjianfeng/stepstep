package utils

import (
	"github.com/xjianfeng/gocomm/lru"
	"stepstep/conf"
)

var Userlru *lru.LruQueue
var Uid2UnionIdLru *lru.LruQueue

func init() {
	var cap int32 = 8000
	if conf.CfgServer.LruCap > 0 {
		cap = conf.CfgServer.LruCap
	}
	Userlru, _ = lru.InitLruCap(cap)
	Uid2UnionIdLru, _ = lru.InitLruCap(cap)
}
