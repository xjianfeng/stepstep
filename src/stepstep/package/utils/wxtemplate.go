package utils

import (
	"encoding/json"
	"github.com/xjianfeng/gocomm/lhttp"
	"stepstep/define"
)

type TemplateData struct {
	AccessToken string
	FormId      string
	OpenId      string
	TemplateId  string
	Page        string
	KwData      map[string]interface{}
}

func SendWxTemplateMsg(data *TemplateData) {
	link := "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token="
	url := link + data.AccessToken
	postData := map[string]interface{}{
		"touser":      data.OpenId,
		"template_id": data.TemplateId,
		"page":        data.Page,
		"form_id":     data.FormId,
		"data":        data.KwData,
	}
	body, err := json.Marshal(postData)
	if err != nil {
		return
	}
	go func() {
		lhttp.HttpPost(url, body, define.HTTP_HEDER_JSON)
	}()
}
