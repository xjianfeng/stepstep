package api

import (
	"github.com/gin-gonic/gin"
	"stepstep/service"
)

/**
* showdoc
* @catalog 领取奖励接口
* @title 获取超时奖励
* @description 获取超时奖励
* @method post
* @url /step/award/timeout
* @param awardType 必选 string 奖励类型,redpack:倒计时红包,advideo:视频步数,luckstep:幸运步数
* @return {"code":0,"errMsg":"","data":{"money":12,"step":1.3,"nextTime":11212,"awardStep":1,"awardMoney":1,"taskList":[{"awardStep":1,"status":0,"awardMoney":1.1}]}}
* @return_param code int 返回状态码0为成功其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data->money float 当前金币
* @return_param data->step int 当前步数
* @return_param data->nextTime int 下一个超时时间
* @return_param data->awardMoney float 奖励的金币
* @return_param data->awardStep int 奖励的步数
* @return_param data->taskList array 刷新步数任务列表
* @remark 这里是备注信息
* @number 99
 */
func ApiTimeoutAward(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	args := &service.ReqAwardArgs{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}

	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	result, err := service.GetAwardByType(unionId, args)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.U = result.U
	ret.Body = gin.H{
		RET_KEY_STEP:        result.U.GetStep(),
		RET_KEY_MONEY:       result.U.Money,
		RET_KEY_AWARD_STEP:  result.AwardStep,
		RET_KEY_AWARD_MONEY: result.AwardMoney,
		RET_KEY_TASK_LIST:   result.U.ResetData.AwardList,
		RET_KEY_NEXT_TIME:   result.NextTimeout,
	}
	return ret
}

/**
* showdoc
* @catalog 领取奖励接口
* @title 获取步数任务奖励
* @description 获取步数任务奖励
* @method post
* @url /step/award/task
* @param idx 必选 int 领取第几个红包,数组下标,从0开始
* @return {"code":0,"errMsg":"","data":{"money":12,"awardMoney":1.3}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data->money string 当前金币
* @return_param data->awardMoney string 获取的金币数
* @remark 这里是备注信息
* @number 99
 */
func ApiTaskAward(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	args := &service.ReqStepAward{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	result, err := service.GetStepTaskAward(unionId, args)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.U = result.U
	ret.Body = gin.H{
		RET_KEY_MONEY:       result.U.Money,
		RET_KEY_AWARD_MONEY: result.Money,
	}
	return ret
}
