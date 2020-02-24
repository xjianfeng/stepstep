package api

import (
	"github.com/gin-gonic/gin"
	"stepstep/service"
)

/**
* showdoc
* @catalog 好友接口
* @title 帮助好友
* @description 点击链接帮助好友
* @method post
* @url /step/friend/help
* @param unionId 必选 string 好友的unionId
* @return {"code":0,"errMsg":"","data":{}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @remark 这里是备注信息
* @number 99
 */
func ApiHelpFriend(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	args := &service.ReqHelpFriend{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	err = service.SetHelpFriend(unionId, args)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Body = gin.H{}
	return ret
}

/**
* showdoc
* @catalog 领取奖励接口
* @title 好友互助红包
* @description 获取好友互助红包
* @method post
* @url /step/award/friend
* @return {"code":0,"errMsg":"","data":{"money":12,"awardMoney":1.3,"friendHelpCnt":1}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data->money string 当前金币
* @return_param data->awardMoney string 获取的金币数
* @return_param data->friendHelpCnt string 剩余的好友互助次数
* @remark 这里是备注信息
* @number 99
 */
func ApiFriendRedpack(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	result, err := service.GetFriendHelpRedpack(unionId)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.U = result.U
	ret.Body = gin.H{
		RET_KEY_MONEY:          result.U.Money,
		RET_KEY_AWARD_MONEY:    result.AwardMoney,
		RET_KEY_FRIEND_HELPCNT: result.U.GetFriendHelpCnt(),
	}
	return ret
}
