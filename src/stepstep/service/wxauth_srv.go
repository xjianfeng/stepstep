package service

import (
	"errors"
	"github.com/xjianfeng/gocomm/decry"
	"github.com/xjianfeng/gocomm/logger"
	"stepstep/conf"
	"stepstep/define"
	"stepstep/models"
	"stepstep/package/auth/wechat"
	"strconv"
	"time"
)

type ReqWxCode struct {
	Code string `json:"code" binding:"required"`
}

type DataWxCode struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func AuthWxCode(args *ReqWxCode) (map[string]string, error) {
	bRet := wechat.Code2Session(conf.CfgWechat.Appid, conf.CfgWechat.Appsecret, args.Code)

	v := &DataWxCode{}
	err := json.Unmarshal(bRet, v)
	if err != nil {
		return nil, err
	}
	if v.ErrMsg != "" {
		logger.LogError("AuthWxCode ErrMsg:%s", v.ErrMsg)
		return nil, define.ERROR_WECHAT_AUTH
	}
	if v.UnionId == "" {
		v.UnionId = v.OpenId
	}
	if v.SessionKey != "" {
		models.SetRedisUInfo(v.UnionId, define.REDIS_HEKY_USER_SESSIONKEY, v.SessionKey)
	}
	now := strconv.FormatInt(time.Now().Unix(), 10)
	token := decry.Md5Sum([]byte(now))
	models.SetRedisUInfo(v.UnionId, define.REDIS_HEKY_USER_TOKEN, token)

	ret := map[string]string{
		"openId":  v.OpenId,
		"unionId": v.UnionId,
		"token":   token,
	}
	return ret, nil
}

type ReqWxUserInfo struct {
	EncrypteData string `json:"encrypteData" binding:"required"`
	WxIv         string `json:"wxIv" binding:"required"`
}

func AuthWxUserInfo(unionId string, args *ReqWxUserInfo) error {
	sessionKey := models.GetRedisUInfo(unionId, define.REDIS_HEKY_USER_SESSIONKEY)
	if sessionKey == nil {
		return errors.New("can't find sessionKey")
	}
	decryData := wechat.DecryptData(args.EncrypteData, string(sessionKey), args.WxIv)
	userInfo := &models.WxData{}
	err := json.Unmarshal(decryData, userInfo)
	if err != nil {
		return err
	}
	user, err := GetUser(unionId)
	if err != nil {
		return err
	}
	user.SetUserInfo(userInfo)
	user.Save()
	return nil
}

type ReqWxStepInfo struct {
	EncrypteData string `json:"encrypteData" binding:"required"`
	WxIv         string `json:"wxIv" binding:"required"`
}

type WxStepData struct {
	Timestamp int64
	Step      int
}

type WxStepInfo struct {
	StepInfoList []WxStepData
}

func AuthWxStepData(unionId string, args *ReqWxStepInfo) (*models.User, error) {
	sessionKey := models.GetRedisUInfo(unionId, define.REDIS_HEKY_USER_SESSIONKEY)
	if sessionKey == nil {
		return nil, errors.New("can't find sessionKey")
	}
	user, err := GetRefreshUser(unionId)
	decryData := wechat.DecryptData(args.EncrypteData, string(sessionKey), args.WxIv)
	stepInfo := WxStepInfo{}
	err = json.Unmarshal(decryData, &stepInfo)
	if err != nil {
		logger.LogError("AuthWxStepData Error:%s", err.Error())
		return user, nil
	}
	if err != nil {
		return nil, err
	}
	stepList := stepInfo.StepInfoList
	l := len(stepList)
	if l == 0 {
		return nil, define.ERROR_DENCRY_STEP_DATA
	}
	newStep := stepList[l-1].Step
	if user.ResetData.RealStep == newStep {
		return user, nil
	}
	user.SetRealStep(newStep)
	user.CheckTaskFinish()
	setUserSportData(user, stepInfo)
	return user, nil
}
