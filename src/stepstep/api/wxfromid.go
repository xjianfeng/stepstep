package api

import (
	"github.com/gin-gonic/gin"
	"stepstep/service"
)

/**
* showdoc
* @catalog 微信相关接口
* @title 收集微信ForomId
* @description 收集微信ForomId
* @method post
* @url /step/wechat/formid
* @param formId 必选 string 用户名
* @return {"code":0,"errMsg":"","data":{}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @remark 这里是备注信息
* @number 99
 */
func ApiWxFromId(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	args := &service.ReqWxFormId{}
	ret.Err = c.BindJSON(args)

	if ret.Err != nil {
		return ret
	}
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	err = service.AddWxFormId(unionId, args)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Body = gin.H{}
	return ret
}
