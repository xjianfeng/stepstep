package service

import (
	"errors"
	"github.com/xjianfeng/gocomm/lhttp"
	"stepstep/conf"
	"stepstep/define"
)

type PayRequestJson struct {
	OpenId   string
	UnionId  string
	UserId   int64
	UserName string
	Amount   int
	Channel  string
	Parform  string
	Desc     string
	GoodsId  string
	PayType  int32 //1 充值 2 购买
	Ip       string
	Extend   string
	StrData  map[string]string
	IntData  map[string]int32
}

type RespWxPrePay struct {
	Code   int                    `json:"code"`
	ErrMsg string                 `json:"errMsg"`
	Data   map[string]interface{} `json:"data"`
}

func WxPrePay(unionId string, money int) (interface{}, error) {
	u, err := GetUser(unionId)
	if err != nil {
		return nil, err
	}
	url := conf.CfgServer.PayServed
	postData := &PayRequestJson{
		OpenId:   u.OpenId,
		UnionId:  u.UnionId,
		UserId:   int64(u.Uid),
		UserName: u.NickName,
		Amount:   money,
		Channel:  "wechat",
		Parform:  "mini",
		Desc:     "充值",
		GoodsId:  "1",
		PayType:  1,
		StrData:  map[string]string{"CallBackUrl": conf.CfgServer.PayCallBack},
	}
	body, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}
	respContent, err := lhttp.HttpPost(url, body, define.HTTP_HEDER_JSON)
	if err != nil {
		return nil, err
	}
	retJson := &RespWxPrePay{}
	err = json.Unmarshal(respContent, retJson)
	if err != nil {
		return nil, err
	}
	if retJson.Code != 0 {
		return nil, errors.New(retJson.ErrMsg)
	}
	return retJson.Data, nil
}

func WxPayCallBack(unionId string) error {
	u, err := GetUser(unionId)
	if err != nil {
		return err
	}
	u.PayFlg = true
	u.SaveAlway()
	return nil
}
