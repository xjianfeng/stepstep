package api

import (
	"github.com/gin-gonic/gin"
	"stepstep/service"
)

/**
* showdoc
* @catalog
* @title 资产明细
* @description 获取资产奖励
* @method post
* @url /step/asset/info
* @param pageIndex 必选 int 第几页
* @param pageSize 必选 int 一页数量表示显示多少天数据，不是多少条记录
* @return {"code":0,"errMsg":"","data":{"nextPage":false,"list":[{"date":"2019-06-24","list":[{"htime":"16:31:52","desc":"运动奖励","value":1.2,"modtype":1}], "summary":"2.40"}]}}
* @return_param code int 返回状态码0为成功其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data->nextPage bool 是否有下一页
* @return_param data->list array 数组
* @return_param data->list->date string 天数
* @return_param data->list->list->htime string 时间 时分秒
* @return_param data->list->list->desc string 描述
* @return_param data->list->list->value string 变化值
* @return_param data->list->list->modtype string 变化类型1加2减
* @return_param data->list->summary string 汇总
* @remark 这里是备注信息
* @number 99
 */
func ApiAssetDetail(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	args := &service.ReqAssetDetail{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}
	ret.Body = service.GetAssetDetail(unionId, args)
	return ret
}
