package api

import (
	"github.com/gin-gonic/gin"
	"stepstep/define"
	"stepstep/service"
)

var (
	payPrice = 500
)

/**
* showdoc
* @catalog 微信相关接口
* @title 微信支付
* @description 微信支付
* @method post
* @url /step/wxpay/create
* @return {"code":0,"errMsg":"","data":{"appId":"xx","timeStamp":11111,"nonceStr":"","package":"","signType":"","orderId":"","paySign":""}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data json 返回数据
* @return_param data->appId string appId
* @return_param data->timeStamp int 时间戳
* @return_param data->nonceStr string 随机串
* @return_param data->package string package
* @return_param data->signType string MD5签名方式
* @return_param data->paySign string 签名
* @remark data 说明
* @number 99
 */
func ApiWxPayCreate(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Body, ret.Err = service.WxPrePay(unionId, payPrice)
	return ret
}

func ApiWxPayCallBack(c *gin.Context) *ApiRetData {
	ret := NewStrRet()
	ret.Body = `{"ret":"Fail"}`
	args := map[string]interface{}{}
	ret.Err = c.BindJSON(&args)
	if ret.Err != nil {
		return ret
	}
	unionId := args["unionId"].(string)
	price := args["payPrice"].(int)
	if price != payPrice {
		ret.Err = define.ERROR_PAY_PRICE
		return ret
	}
	err := service.WxPayCallBack(unionId)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Body = `{"ret":"OK"}`
	return ret
}
