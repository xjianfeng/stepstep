package api

import (
	"github.com/gin-gonic/gin"
	"stepstep/define"
	"stepstep/models"
)

var (
	RET_KEY_STEP           = "step"
	RET_KEY_MONEY          = "money"
	RET_KEY_LUCKSTEP       = "luckStep"
	RET_KEY_REDPACK        = "redpack"
	RET_KEY_ADVIDEO        = "adVideo"
	RET_KEY_FRIEND_HELPCNT = "friendHelpCnt"
	RET_KEY_TASK_LIST      = "taskList"
	RET_KEY_FRIEND_LIST    = "friendList"
	RET_KEY_AWARD_STEP     = "awardStep"
	RET_KEY_AWARD_MONEY    = "awardMoney"
	RET_KEY_IMAGE          = "image"
	RET_KEY_SOUND          = "sound"
	RET_KEY_NEXT_TIME      = "nextTime"
)

type ApiRetData struct {
	U        *models.User
	Err      error
	RespType string
	Body     interface{}
}

func NewJsonRet() *ApiRetData {
	r := new(ApiRetData)
	r.RespType = define.RESPONSE_TYPE_JSON
	return r
}

func NewStrRet() *ApiRetData {
	r := new(ApiRetData)
	r.RespType = define.RESPONSE_TYPE_STRING
	return r
}

func GetUnionId(c *gin.Context) (string, error) {
	unionId := c.GetHeader("unionId")
	if unionId == "" {
		return unionId, define.ERROR_REQUEST_PARAMS
	}
	return unionId, nil
}

func Test(c *gin.Context) *ApiRetData {
	ret := new(ApiRetData)
	ret.RespType = define.RESPONSE_TYPE_STRING
	ret.Body = "Hello World!"
	return ret
}
