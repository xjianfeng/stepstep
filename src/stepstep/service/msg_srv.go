package service

import (
	"stepstep/conf"
	"stepstep/define"
	"stepstep/models"
	"stepstep/package/auth/wechat"
	"stepstep/package/utils"
	"time"
)

type ReqMessage struct {
	Uid     int    `json:"uId" `
	UnionId string `json:"unionId"`
	Msg     string `json:"msg" binding:"required"`
}

func AddMesage(unionId string, args *ReqMessage) (string, error) {
	u, err := GetUser(unionId)
	if err != nil {
		return "", err
	}
	if args.UnionId == "" {
		return "", define.ERROR_REQUEST_PARAMS
	}
	stime := time.Now().Format(define.TIME_FORMAT)
	msgId, err := models.AddMsg(args.UnionId, u, args.Msg)
	//发送微信模板消息
	formId, openId := models.GetRedisFormId(args.UnionId)
	if formId != "" && openId != "" {
		template := &utils.TemplateData{
			AccessToken: wechat.GetAccessToken(conf.CfgWechat.Appid, false),
			FormId:      formId,
			OpenId:      openId,
			TemplateId:  conf.CfgWechat.TemplateId,
			Page:        conf.CfgWechat.Page,
			KwData: map[string]interface{}{
				"keyword1": map[string]string{"value": "消息留言"},
				"keyword2": map[string]string{"value": args.Msg},
				"keyword3": map[string]string{"value": "看看是谁"},
				"keyword4": map[string]string{"value": stime},
			},
		}
		utils.SendWxTemplateMsg(template)
	}
	return msgId, err
}

type ReqMessageList struct {
	Uid       int    `json:"uId"`
	UnionId   string `json:"unionId"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize binding:"required"`
}

type RetMessageList struct {
	Timestamp int64
	PayFlg    bool
	NextPage  bool
	List      []models.MsgInfo
}

func GetMessageList(reqUnionId string, args *ReqMessageList) (*RetMessageList, error) {
	u, err := GetUser(args.UnionId)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	result, err := models.GetMsgList(args.UnionId, args.PageIndex*args.PageSize, args.PageSize)
	if err != nil {
		return nil, err
	}
	ret := &RetMessageList{}
	ret.List = result.MsgList
	ret.NextPage = result.MsgCnt > args.PageIndex*args.PageSize+args.PageSize
	ret.Timestamp = now
	if reqUnionId == u.UnionId {
		u.ReqMsgTime = now
		ret.PayFlg = u.PayFlg
	}
	return ret, nil
}

type ReqReply struct {
	MsgId string
	Msg   string
}

func AddReply(unionId string, args *ReqReply) error {
	u, err := GetUser(unionId)
	if err != nil {
		return err
	}
	models.AddReply(args.MsgId, u, args.Msg)
	//发送微信模板消息
	stime := time.Now().Format(define.TIME_FORMAT)
	formId, openId := models.GetRedisFormId(unionId)
	if formId != "" && openId != "" {
		template := &utils.TemplateData{
			AccessToken: wechat.GetAccessToken(conf.CfgWechat.Appid, false),
			FormId:      formId,
			OpenId:      openId,
			TemplateId:  conf.CfgWechat.TemplateId,
			Page:        conf.CfgWechat.Page,
			KwData: map[string]interface{}{
				"keyword1": map[string]string{"value": "消息回复"},
				"keyword2": map[string]string{"value": args.Msg},
				"keyword3": map[string]string{"value": "查看回复"},
				"keyword4": map[string]string{"value": stime},
			},
		}
		utils.SendWxTemplateMsg(template)
	}
	return err
}
