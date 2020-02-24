package api

import (
	"github.com/gin-gonic/gin"
	"stepstep/service"
)

/**
* showdoc
* @catalog 微信相关接口
* @title 微信code登录
* @description 微信code登录
* @method post
* @url /step/wechat/login
* @param code 必选 string 微信获取的code
* @return {"code":0,"errMsg":"","data":{"unionId":"xx","uId":112121, "token":"xxx"}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data json 返回数据
* @return_param data->uId int 用户Id
* @return_param data->unionId string 用户唯一标识
* @return_param data->token string api请求token
* @remark data 说明
* @number 99
 */
func ApiWxAuthCode(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	args := &service.ReqWxCode{}
	ret.Err = c.BindJSON(args)

	if ret.Err != nil {
		return ret
	}
	data, err := service.AuthWxCode(args)
	if err != nil {
		ret.Err = err
		return ret
	}
	openId := data["openId"]
	unionId := data["unionId"]
	token := data["token"]
	u, err := service.GetUserData(unionId, openId)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Body = gin.H{"uId": u.Uid, "unionId": u.UnionId, "token": token}
	return ret
}

/**
* showdoc
* @catalog 微信相关接口
* @title 用户授权信息
* @description 用户授权信息
* @method post
* @url /step/wechat/userinfo
* @param encrypteData 必选 string 微信加密数据
* @param wxIv  必选 string 微信获取的IV向量
* @return {"code":0,"errMsg":"","data":{}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @remark 这里是备注信息
* @number 99
 */
func ApiAuthWxUserInfo(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	args := &service.ReqWxUserInfo{}

	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Err = service.AuthWxUserInfo(unionId, args)
	if ret.Err != nil {
		return ret
	}
	ret.Body = gin.H{}
	return ret
}
