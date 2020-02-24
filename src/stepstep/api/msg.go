package api

import (
	"github.com/gin-gonic/gin"
	"stepstep/service"
)

/**
* showdoc
* @catalog 留言接口
* @title 添加留言
* @description 添加留言
* @method post
* @url /step/msg/add
* @param uId 必选 int 被留言者的uId
* @param unionId 必选 string 被留言者的unionId
* @param msg 必选 string 留言
* @return {"code":0,"errMsg":"","data":{"msgId":"xxxx"}}
* @return_param code int 返回状态码0:为成功,其他为错误
* @return_param errMsg string 错误信息，当code不为0时才有
* @return_param data->msgId string 留言ID
* @remark data 说明
* @number 99
 */
func ApiAddMessage(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	args := &service.ReqMessage{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	//通过UID获取unionId
	if args.Uid != 0 && args.UnionId == "" {
		args.UnionId, ret.Err = service.GetUnionIdByUid(args.Uid)
		if ret.Err != nil {
			return ret
		}
	}
	msgId, err := service.AddMesage(unionId, args)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Body = gin.H{"msgId": msgId}
	return ret
}

/**
* showdoc
* @catalog 留言接口
* @title 获取留言
* @description 获取留言
* @method post
* @url /step/msg/list
* @param uId 必选 int 被留言者的uId
* @param unionId 必选 string 被查看人的unionId
* @param pageIndex 必选 int 第几页从0开始
* @param pageSize 必选 int 每页数量
* @return {"code":0,"errMsg":"","data":{"timestamp":12332,"nextPage":false,"payflg":true,"list":[{"msgId":"ssss","uid":1212,"avatarUrl":"","nickName":"","msg":"","reply":{"avatarUrl":"","nickName":"","msg":"","uid":121}}]}}
* @return_param code int 返回状态码0为成功，其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data json 返回数据
* @return_param data->timestamp int 时间戳,下次请求带上
* @return_param data->payflg bool 是否支付,用来显示头像和名称
* @return_param data->nextPage bool 是否有下一页
* @return_param data->list->msgId string 留言ID
* @return_param data->list->avatarUrl string 微信头像(头像可能没有)
* @return_param data->list->nickName string 微信昵称(可能没有)
* @return_param data->list->msg string 留言信息
* @return_param data->list->uid int 用户uid
* @return_param data->list->reply json 回复留言
* @remark
* @number 99
 */
func ApiMessageList(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	args := &service.ReqMessageList{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	//通过UID获取unionId
	if args.Uid != 0 && args.UnionId == "" {
		args.UnionId, ret.Err = service.GetUnionIdByUid(args.Uid)
		if ret.Err != nil {
			return ret
		}
	}
	result, err := service.GetMessageList(unionId, args)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Body = gin.H{
		"timestamp": result.Timestamp,
		"list":      result.List,
		"payflg":    result.PayFlg,
		"nextPage":  result.NextPage,
	}
	return ret
}

/**
* showdoc
* @catalog 留言接口
* @title 回复留言
* @description 回复留言
* @method post
* @url /step/msg/reply
* @param msgId 必选 string 回复那条留言的Id
* @param msg 必选 string 回复内容
* @return {"code":0,"errMsg":"","data":{}}
* @return_param code int 返回状态码0:为成功,其他为错误
* @return_param errMsg string 错误信息，当code不为0时才有
* @remark data 说明
* @number 99
 */
func ApiAddReply(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	args := &service.ReqReply{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Err = service.AddReply(unionId, args)
	ret.Body = gin.H{}
	return ret
}
