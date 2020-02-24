package service

import (
	"fmt"
	"github.com/xjianfeng/gocomm/logger"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"stepstep/conf"
	"stepstep/define"
	"stepstep/models"
	"stepstep/package/auth/wechat"
	"stepstep/package/cos"
	"stepstep/package/utils"
	"strconv"
)

type ReqSportData struct {
	Uid     int
	UnionId string
}

type RespSportData struct {
	SportData *models.SportData
	NickName  string
	AvatarUrl string
	NewMsg    int
	MsgNum    int
}

func GetSportData(reqUnionId string, args *ReqSportData) (*RespSportData, error) {
	u, err := GetUser(args.UnionId)
	if err != nil {
		return nil, err
	}
	newMsg := 0
	//TODO 缓存优化
	msgNum, err := models.GetMsgNum(bson.M{"unionid": args.UnionId})
	if err != nil {
		logger.LogError("GetMsgNum error:%s", err.Error())
	}
	// 自己请求才有是否有留言
	if args.UnionId == reqUnionId {
		selobj := bson.M{"unionid": args.UnionId, "createtime": bson.M{"$gt": u.ReqMsgTime}}
		newMsg, err = models.GetMsgNum(selobj)
		if err != nil {
			logger.LogError("GetMsgNum error:%s", err.Error())
		}
	}
	ret := &RespSportData{
		SportData: &u.SportInfo,
		NickName:  u.NickName,
		AvatarUrl: u.AvatarUrl,
		NewMsg:    newMsg,
		MsgNum:    msgNum,
	}
	return ret, nil
}

// 设置微信运动数据
func setUserSportData(u *models.User, wxStep WxStepInfo) {
	l := len(wxStep.StepInfoList)
	if l == 0 {
		return
	}
	u.AddKeepDay(wxStep.StepInfoList[l-1].Step)
	u.SportInfo.StepList[0] = wxStep.StepInfoList[l-1].Step
	if wxStep.StepInfoList[l-2].Timestamp > u.CreateTime {
		u.SportInfo.StepList[1] = wxStep.StepInfoList[l-2].Step
	}
	if wxStep.StepInfoList[l-3].Timestamp > u.CreateTime {
		u.SportInfo.StepList[2] = wxStep.StepInfoList[l-3].Step
	}
	setQrCode(u)
}

func setQrCode(u *models.User) {
	if u.SportInfo.QrCode != "" {
		return
	}
	scene := fmt.Sprintf("uid=%d", u.Uid)
	page := "pages/poster_share/poster_share"
	page = ""
	bRet, err := wechat.GetQrCode(conf.CfgWechat.Appid, scene, page)
	if err != nil {
		logger.LogError("setQrCode error:%s", err.Error())
		return
	}
	fileName := fmt.Sprintf("stepstep/qrcode/qrcode_%s.png", u.UnionId)
	err = cos.UploadBytes(fileName, bRet)
	if err != nil {
		logger.LogError("setQrCode error:%s", err.Error())
		return
	}
	u.SportInfo.QrCode = conf.CfgCos.Domain + fileName
}

/*
设置三张封面图片，重复上传替换当前天的, 超出三天的删除
*/
func SetSportImage(unionId, image string) {
	u, err := GetUser(unionId)
	if err != nil {
		logger.LogError("SetSportImage error:%s", err.Error())
		return
	}
	dayNo := utils.GetRealDayNo()
	nowDayIdx := strconv.Itoa(utils.GetRealDayNo())
	if u.SportInfo.ImageList == nil {
		u.SportInfo.ImageList = make(map[string]string)
	}
	u.SportInfo.ImageList[nowDayIdx] = image
	l := len(u.SportInfo.ImageList)
	// 删掉超过了3天旧数据
	if l > 3 {
		dayIdx := 0
		for k, v := range u.SportInfo.ImageList {
			dayIdx, _ = strconv.Atoi(k)
			if dayNo-dayIdx < 3 {
				continue
			}
			cos.DeleteObject(v)
			delete(u.SportInfo.ImageList, k)
		}
	}
	u.SaveAlway()
}

type ReqSportLike struct {
	Uid     int
	UnionId string
	OpType  int8 //操作 0 添加 1 取消
}

func SetSportLike(args *ReqSportLike) error {
	unionId := args.UnionId
	u, err := GetUser(unionId)
	if err != nil {
		return err
	}
	if args.OpType == 0 {
		u.AddSportLike()
		logger.LogInfo("SportLike %d", u.SportInfo.Like)
		return nil
	}
	u.SubSportLike()
	logger.LogInfo("SportLike %d", u.SportInfo.Like)
	return nil
}

type ReqSportSound struct {
	SoundUrl string
	SoundSec int
}

func SetSportSound(unionId string, args *ReqSportSound) error {
	u, err := GetUser(unionId)
	if err != nil {
		return err
	}
	dayNo := utils.GetRealDayNo()
	nowDayIdx := strconv.Itoa(utils.GetRealDayNo())
	if u.SportInfo.SoundList == nil {
		u.SportInfo.SoundList = make(map[string]models.SoundInfo)
	}
	u.SportInfo.SoundList[nowDayIdx] = models.SoundInfo{
		Sound:    args.SoundUrl,
		SoundSec: args.SoundSec,
	}
	l := len(u.SportInfo.SoundList)
	// 删掉超过了3天旧数据
	if l > 3 {
		dayIdx := 0
		for k, v := range u.SportInfo.SoundList {
			dayIdx, _ = strconv.Atoi(k)
			if dayNo-dayIdx < 3 {
				continue
			}
			cos.DeleteObject(v.Sound)
			delete(u.SportInfo.SoundList, k)
		}
	}
	u.SaveAlway()
	return nil
}

func SetSportCover(unionId string) (string, error) {
	i := rand.Intn(3) + 1
	image := conf.CfgCos.Domain + fmt.Sprintf("stepstep/cover/%d.jpg", i)
	SetSportImage(unionId, image)
	return image, nil
}

type ReqSportRedpack struct {
	Idx int
}

type RespSportRedpack struct {
	U          *models.User
	AwardMoney float64
}

func GetSportRedpack(unionId string, args *ReqSportRedpack) (*RespSportRedpack, error) {
	u, err := GetUser(unionId)
	if err != nil {
		return nil, err
	}
	if args.Idx >= define.SPORT_REDPACK_NUM {
		return nil, define.ERROR_REQUEST_PARAMS
	}
	if u.SportInfo.KeepRedPack[args.Idx] == define.STATUS_INIT {
		return nil, define.ERROR_AWARD_TASK_NOTFINISH
	}
	if u.SportInfo.KeepRedPack[args.Idx] == define.STATUS_FINISH {
		return nil, define.ERROR_AWARD_GOT
	}
	randMoney := define.AWARD_SPORT_REDPACK[args.Idx]
	awardMoney := float64(rand.Intn(randMoney)+1) / 10.0
	u.AddMoney(awardMoney, define.SYSTEM_KEEP_REPACK, define.DESC_TYPE_KEEP_REDPACK)
	u.SetKeepRedpack(args.Idx, define.STATUS_FINISH)
	ret := &RespSportRedpack{
		U:          u,
		AwardMoney: awardMoney,
	}
	return ret, nil
}
