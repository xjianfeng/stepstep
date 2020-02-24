package define

import (
	"errors"
)

var (
	ERROR_WECHAT_AUTH          = errors.New("微信认证错误")
	ERROR_USER_NOT_EXISTS      = errors.New("用户不存在")
	ERROR_REQUEST_PARAMS       = errors.New("参数错误")
	ERROR_DENCRY_STEP_DATA     = errors.New("获取步数失败")
	ERROR_AWARD_FAIL           = errors.New("领取失败")
	ERROR_AWARD_NOT_EXISTS     = errors.New("奖励不存在")
	ERROR_AWARD_GOT            = errors.New("奖励已领取")
	ERROR_AWARD_TASK_NOTFINISH = errors.New("任务未完成")
	ERROR_SERVER_ERROR         = errors.New("服务器错误")
	ERROR_HELP_FAIL            = errors.New("自己点击无效")
	ERROR_PAY_PRICE            = errors.New("支付金钱错误")
)
