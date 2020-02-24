package routers

import (
	"github.com/gin-gonic/gin"
	"stepstep/api"
	"stepstep/conf"
	"stepstep/middleware"
)

func SetRouter() *gin.Engine {
	mode := gin.DebugMode
	if conf.IsRelease() {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	r := gin.Default()
	g := r.Group("step/")
	g.Use(middleware.CustomMiddleWare())

	g.POST("index", middleware.DecorateFunc(api.Index))
	g.POST("upload", middleware.DecorateFunc(api.ApiUploadImg))
	g.POST("upload/sound", middleware.DecorateFunc(api.ApiUploadSound))
	g.POST("complain", middleware.DecorateFunc(api.ApiUploadImg))
	g.POST("roll/list", middleware.DecorateFunc(api.ApiRollAward))
	g.POST("wechat/login", middleware.DecorateFunc(api.ApiWxAuthCode))
	g.POST("wechat/userinfo", middleware.DecorateFunc(api.ApiAuthWxUserInfo))
	g.POST("wechat/formid", middleware.DecorateFunc(api.ApiWxFromId))
	g.POST("award/task", middleware.DecorateFunc(api.ApiTaskAward))
	g.POST("award/timeout", middleware.DecorateFunc(api.ApiTimeoutAward))
	g.POST("award/friend", middleware.DecorateFunc(api.ApiFriendRedpack))
	g.POST("award/sport", middleware.DecorateFunc(api.ApiSportRedpack))
	g.POST("friend/help", middleware.DecorateFunc(api.ApiHelpFriend))
	g.POST("friend/refresh", middleware.DecorateFunc(api.ApiFriendRefresh))
	g.POST("msg/add", middleware.DecorateFunc(api.ApiAddMessage))
	g.POST("msg/list", middleware.DecorateFunc(api.ApiMessageList))
	g.POST("msg/reply", middleware.DecorateFunc(api.ApiAddReply))
	g.POST("sport/data", middleware.DecorateFunc(api.ApiSportData))
	g.POST("sport/like", middleware.DecorateFunc(api.ApiSportLike))
	g.POST("sport/sound", middleware.DecorateFunc(api.ApiSportSound))
	g.POST("sport/cover", middleware.DecorateFunc(api.ApiSportCover))
	g.POST("wxpay/create", middleware.DecorateFunc(api.ApiWxPayCreate))
	g.POST("wxpay/callback", middleware.DecorateFunc(api.ApiWxPayCallBack))
	g.POST("asset/info", middleware.DecorateFunc(api.ApiAssetDetail))

	return r
}
