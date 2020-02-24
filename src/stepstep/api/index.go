package api

import (
	"github.com/gin-gonic/gin"
	"stepstep/define"
	"stepstep/package/appdata"
	"stepstep/service"
)

/**
* showdoc
* @catalog
* @title 首页
* @description 获取首页数据接口
* @method post
* @url /step/index
* @param encrypteData 必选 string 步数的加密数据
* @param wxIv 必选 string 步数的IV向量
* @return {"code":0,"errMsg":"","data":{"step":199,"money":2.9,"luckStep":0,"redpack":0,"adVideo":0,"friendHelpCnt":0,"taskList":[{"AwardStep":1,"AwardMoney":0.5,"Status":0}],"friendList":["http://xxx"]}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息，当code不为0时才有
* @return_param data json  返回数据
* @return_param data->step int 步数
* @return_param data->money float 当前金币数
* @return_param data->luckStep int 幸运步数是否开启0:未开启1:开启
* @return_param data->redpack int 红包倒计时间戳
* @return_param data->adVideo int 看视频得奖励是否开启0:未开启1:开启
* @return_param data->friendHelpCnt int 好友帮助次数
* @return_param data->taskList array 步数任务列表
* @return_param data->friendList array 好友互助列表
* @return_param data->taskList->AwardStep int 达到奖励的步数
* @return_param data->taskList->AwardMoney float 奖励的金额
* @return_param data->taskList->Status int 状态0:未达到,1:达到未领取,2:已领取
* @remark 这里是备注信息
* @number 99
 */
func Index(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	args := &service.ReqWxStepInfo{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}
	u, err := service.AuthWxStepData(unionId, args)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.U = u
	luckStepData := u.GetTimeoutData(define.TIMEOUT_KEY_LUCKSTEP)
	adVideoData := u.GetTimeoutData(define.TIMEOUT_KEY_VIDEO)
	redPackData := u.GetTimeoutData(define.TIMEOUT_KEY_REDPACK)

	ret.Body = gin.H{
		RET_KEY_STEP:  u.GetStep(),
		RET_KEY_MONEY: u.Money,
		//是否显示奖励步数
		RET_KEY_LUCKSTEP: luckStepData.Status,
		//显示红包倒计时
		RET_KEY_REDPACK: redPackData.TimeOut,
		//视频红包是否显示
		RET_KEY_ADVIDEO:        adVideoData.Status,
		RET_KEY_FRIEND_HELPCNT: u.GetFriendHelpCnt(),
		RET_KEY_TASK_LIST:      u.ResetData.AwardList,
		RET_KEY_FRIEND_LIST:    u.GetFriendIcons(),
	}
	return ret
}

/**
* showdoc
* @catalog
* @title 获取侧边滚动数据
* @description 获取运动数据
* @method post
* @url /step/roll/list
* @return {"code":0,"errMsg":"","data":{"list":[{"name":"xx","image":"xx","money":50},{"name":"yy","image":"yy","money":150}]}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data->list->name string 姓名
* @return_param data->list->image string 头像
* @return_param data->list->money int 奖励的钱
* @remark 这里是备注信息
* @number 99
 */
func ApiRollAward(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	ret.Body = gin.H{
		"list": appdata.GetRandomUserInfo(50),
	}
	return ret
}

/**
* showdoc
* @catalog
* @title 刷新好友信息
* @description 刷新好友信息
* @method post
* @url /step/friend/refresh
* @return {"code":0,"errMsg":"","data":{"friendHelpCnt":0,"friendList":["http://xxx"]}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息，当code不为0时才有
* @return_param data json  返回数据
* @return_param data->friendHelpCnt int 好友帮助次数
* @return_param data->friendList array 好友互助列表
* @remark 这里是备注信息
* @number 99
 */
func ApiFriendRefresh(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	u, err := service.GetUser(unionId)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Body = gin.H{
		RET_KEY_FRIEND_HELPCNT: u.GetFriendHelpCnt(),
		RET_KEY_FRIEND_LIST:    u.GetFriendIcons(),
	}
	return ret
}
