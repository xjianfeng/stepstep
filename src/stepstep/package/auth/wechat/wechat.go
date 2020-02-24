package wechat

import (
	//"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xjianfeng/gocomm/decry"
	"github.com/xjianfeng/gocomm/lhttp"
	"github.com/xjianfeng/gocomm/logger"
	"io/ioutil"
	"net/http"
	"time"
)

var log = logger.New("wechat_auth.log")

func Code2Session(appId, appSecret, wxCode string) []byte {
	var urlParam lhttp.UrlParam
	urlParam.Data = make(map[string]string)
	urlParam.Data["appid"] = appId
	urlParam.Data["secret"] = appSecret
	urlParam.Data["js_code"] = wxCode
	urlParam.Data["grant_type"] = "authorization_code"
	urlString := urlParam.UrlEncode()
	authUrl := "https://api.weixin.qq.com/sns/jscode2session?" + urlString

	response, ok := http.Get(authUrl)
	var empty []byte
	if ok != nil {
		return empty
	}
	defer response.Body.Close()
	result, ok := ioutil.ReadAll(response.Body)
	log.LogInfo("WxAuthLogin result %s", string(result))
	if ok != nil {
		return empty
	}
	return result
}

func DecryptData(encryData string, key string, iv string) []byte {
	//log.LogInfo("WxDecryptData encryData %s, key %s, iv %s", encryData, key, iv)
	base64Data, _ := base64.StdEncoding.DecodeString(encryData)
	base64Key, _ := base64.StdEncoding.DecodeString(key)
	base64Iv, _ := base64.StdEncoding.DecodeString(iv)
	block, _ := aes.NewCipher(base64Key)
	cipherMod := cipher.NewCBCDecrypter(block, base64Iv)
	decryStr := make([]byte, len(base64Data))
	cipherMod.CryptBlocks(decryStr, base64Data)
	return PKCS5UnPadding(decryStr)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	//log.LogInfo("PKCS5UnPadding %s, length %d, unppadding %d", string(origData), length, unpadding)
	if length < unpadding {
		return []byte{}
	}
	return origData[:(length - unpadding)]
}

type RetAccessToken struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	Expires      int    `json:"expires_in"`
	RefreshToekn string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
}

func WxOAuth2AccessToken(appId string, secret string, code string) *RetAccessToken {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", appId, secret, code)
	respData, err := lhttp.HttpGet(url)
	if err != nil {
		log.LogError("WxOAuth2AccessToken respData %v, err %v", respData, err)
		return nil
	}
	log.LogDebug("WxOAuth2AccessToken url:%s, respData:%s", url, respData)
	retData := &RetAccessToken{}
	err = json.Unmarshal(respData, retData)
	if err != nil {
		log.LogError("WxOAuth2AccessToken retData %v json Unmarshal %v", retData, err.Error())
		return nil
	}
	return retData
}

type RetWxUserInfo struct {
	OpenId     string `json:"openid"`
	NickName   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgUrl string `json:"headimgurl"`
	UnionId    string `json:"unionid"`
}

func GetUserInfoByOAuth(access_token string, openid string) *RetWxUserInfo {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", access_token, openid)
	respData, err := lhttp.HttpGet(url)
	if err != nil {
		log.LogError("GetUserInfoByOAuth respData %v, err %v", respData, err)
		return nil
	}
	log.LogDebug("GetUGetUserInfoByOAuth url:%s, respData:%s", url, respData)
	retData := &RetWxUserInfo{}
	err = json.Unmarshal(respData, retData)
	if err != nil {
		log.LogError("GetUserInfoByOAuth respData %v, err %v", respData, err)
		return nil
	}
	return retData
}

type accResult struct {
	Ret         string
	ErrorMsg    string
	AccessToken string
}

func GetAccessToken(appId string, refresh bool) string {
	url := fmt.Sprintf("http://server.kaizhan8.com/serverapi/accesstoken?appId=%s", appId)
	if refresh {
		url = fmt.Sprintf("%s&refresh=%d", url, 1)
	}
	response, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.LogError("err: %s", err.Error())
		return ""
	}

	ret := &accResult{}
	err = json.Unmarshal(result, ret)
	if err != nil || ret.Ret != "OK" {
		log.LogError("GetAccessToken error:%v, ret:%v", err, ret)
		return ""
	}
	return ret.AccessToken
}

type ResqTick struct {
	Ret      string
	ErrorMsg string
	Ticket   string
}

func GetH5JsApiTicket(appId string) string {
	url := fmt.Sprintf("http://server.kaizhan8.com/serverapi/jsapitick?appId=%s", appId)
	ret, err := lhttp.HttpGet(url)
	if err != nil {
		return ""
	}
	response := &ResqTick{}
	err = json.Unmarshal(ret, response)
	if err != nil {
		return ""
	}
	if response.Ret != "OK" {
		return ""
	}
	return response.Ticket
}

type JsApiConfig struct {
	Noncestr  string
	Timestamp int64
	Sign      string
}

func JsApiSignInfo(appId, url string) *JsApiConfig {
	ticket := GetH5JsApiTicket(appId)
	ts := time.Now().Unix()
	strTs := fmt.Sprintf("%d", ts)
	randomStr := decry.Md5Sum([]byte(strTs))[:10]

	signStr := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, randomStr, ts, url)
	sign := decry.Sha1(signStr)
	log.LogInfo("JsApiSignInfo %s", signStr)
	log.LogInfo("Sign %s", sign)

	return &JsApiConfig{
		Noncestr:  randomStr,
		Timestamp: ts,
		Sign:      sign,
	}
}

func GetQrCode(appId, scene, page string) ([]byte, error) {
	accessToken := GetAccessToken(appId, false)
	if accessToken == "" {
		return nil, errors.New("accessToken error")
	}
	url := "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=" + accessToken
	args := map[string]interface{}{
		"scene":      scene,
		"page":       page,
		"auto_color": true,
	}
	body, err := json.Marshal(&args)
	if err != nil {
		return nil, err
	}
	ret, err := lhttp.HttpPost(url, body, map[string]string{"Content-Type": "application/json"})
	return ret, err
}
