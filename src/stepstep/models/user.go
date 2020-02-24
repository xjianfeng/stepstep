package models

import (
	"fmt"
	"stepstep/define"
	"stepstep/package/utils"
	"strconv"
	"sync"
	"time"
)

/* 用户操作 */

// 微信数据
type WxData struct {
	OpenId    string
	UnionId   string
	NickName  string
	AvatarUrl string
	Gender    int
	City      string
	Province  string
	Country   string
}

//奖励信息
type AwardStatus struct {
	AwardStep  int
	AwardMoney float64
	Status     int8
}

type FriendData struct {
	sync.RWMutex
	HelpCnt        int
	FriendHeadIcon []string
	FriendList     map[string]string
}

// 每日数据
type DayData struct {
	AwardStep  int         //奖励步数
	RealStep   int         //真实步数
	Friend     *FriendData //unionId: AvatarUrl
	IsFinished bool        //任务全部完成
	AwardList  [define.AWARD_LIST_MAX_IDX]AwardStatus
}

// 过期配置
type TimeData struct {
	TimeOut    int64
	Status     int8 // 0 1 2
	Refresh    int
	MaxRefresh int //每日限制最大刷新次数
	Info       map[string]interface{}
}

type SoundInfo struct {
	Sound    string
	SoundSec int
}

type SportData struct {
	Like        int
	KeepDay     int
	Step        int
	TotalStep   int
	QrCode      string
	KeepRedPack [define.SPORT_REDPACK_NUM]int8
	StepList    [define.STEP_HISTORY_DAY]int
	SoundList   map[string]SoundInfo
	ImageList   map[string]string
}

// 用户数据
type User struct {
	*WxData
	Uid        int
	Money      float64
	RealDayNo  int
	CreateTime int64
	ReqMsgTime int64 //请求获取留言时间
	Timeout    map[string]*TimeData
	PayFlg     bool //是否付费
	SaveDbFlg  bool //是否存盘标志
	AwardCnt   int  //第几次领取奖励
	SportDayNo int
	ResetData  *DayData
	SportInfo  SportData
}

func (u *User) AddMoney(money float64, sys string, desc string) {
	u.Money, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", u.Money+money), 64)
	u.SaveDbFlg = true
	addAssetLog(&assestLogInfo{
		unionId: u.UnionId,
		modType: define.MOD_TYPE_ADD,
		system:  sys,
		desc:    desc,
		value:   money,
	})
}

func (u *User) SubMoney(money float64) {

}

func (u *User) SetRealStep(step int) {
	u.ResetData.RealStep = step
	u.SaveDbFlg = true
}

func (u *User) AddAwardStep(step int) {
	u.ResetData.AwardStep += step
	u.SaveDbFlg = true
}

//检查定时数据
func (u *User) CheckTimeout() {
	now := time.Now().Unix()
	for _, v := range u.Timeout {
		if v.TimeOut > now {
			continue
		}
		if v.MaxRefresh != define.TIMEOUT_NOT_LIMIT_REFRESH &&
			v.Refresh >= v.MaxRefresh {
			continue
		}
		if v.Status == define.STATUS_INIT {
			v.Status = define.STATUS_ALREDY
			u.SaveDbFlg = true
		}
	}
}

//检查步数是否达到领奖提交
func (u *User) CheckTaskFinish() {
	if u.ResetData.IsFinished {
		return
	}
	nowStep := u.GetStep()
	for i, item := range u.ResetData.AwardList {
		if item.AwardStep > nowStep {
			continue
		}
		if item.Status == define.STATUS_FINISH || item.Status == define.STATUS_ALREDY {
			continue
		}
		if i == define.AWARD_LIST_MAX_IDX {
			u.ResetData.IsFinished = true
		}
		u.ResetData.AwardList[i].Status = define.STATUS_ALREDY
		u.SaveDbFlg = true
	}
}

func (u *User) SetUserInfo(userInfo *WxData) {
	u.NickName = userInfo.NickName
	u.AvatarUrl = userInfo.AvatarUrl
	u.Gender = userInfo.Gender
	u.City = userInfo.City
	u.Province = userInfo.Province
	u.Country = userInfo.Country
	u.SaveDbFlg = true
}

func (u *User) ResetEveryDay() bool {
	dayno := utils.GetRealDayNo()
	if u.RealDayNo == 0 {
		return false
	}
	if u.RealDayNo == dayno {
		return false
	}
	u.ResetData = &DayData{
		Friend: &FriendData{
			FriendHeadIcon: []string{},
			FriendList:     make(map[string]string),
		},
	}
	u.RealDayNo = utils.GetRealDayNo()
	initAwardList(u)
	u.SaveDbFlg = true
	return true
}

func (u *User) Save() {
	SaveUser(u)
}

func (u *User) SaveAlway() {
	u.SaveDbFlg = true
	SaveUser(u)
}

func (u *User) GetFriendHelpCnt() int {
	u.ResetData.Friend.RLock()
	defer u.ResetData.Friend.RUnlock()

	return u.ResetData.Friend.HelpCnt
}

func (u *User) GetFriendIcons() []string {
	u.ResetData.Friend.RLock()
	defer u.ResetData.Friend.RUnlock()

	return u.ResetData.Friend.FriendHeadIcon
}

func (u *User) CheckTimeOutStatus(key string) bool {
	tdata, ok := u.Timeout[key]
	if !ok {
		return false
	}
	now := time.Now().Unix()
	if tdata.Status != define.STATUS_ALREDY && tdata.TimeOut > now {
		return false
	}
	return true
}

var refKeyValue = map[string]int64{
	define.TIMEOUT_KEY_VIDEO:    define.TIMOUT_READ_VIDEO,
	define.TIMEOUT_KEY_LUCKSTEP: define.TIMOUT_AWARD_STEP,
	define.TIMEOUT_KEY_REDPACK:  define.TIMOUT_RED_REDPACK,
}

func (u *User) SetTimeOutFinish(key string) {
	tdata, ok := u.Timeout[key]
	if !ok {
		return
	}
	now := time.Now().Unix()
	tdata.Status = define.STATUS_INIT
	tdata.TimeOut = now + refKeyValue[key]
	tdata.Refresh += 1
	u.SaveDbFlg = true
}

func (u *User) GetStep() int {
	return u.ResetData.RealStep + u.ResetData.AwardStep
}

func (u *User) GetTimeoutData(key string) *TimeData {
	data, ok := u.Timeout[key]
	if !ok {
		return nil
	}
	return data
}

func (u *User) GetStepAwardMoney(idx int) (float64, error) {
	if idx >= len(u.ResetData.AwardList) {
		return 0, define.ERROR_REQUEST_PARAMS
	}
	if u.ResetData.AwardList[idx].Status != define.STATUS_ALREDY {
		return 0, define.ERROR_AWARD_TASK_NOTFINISH
	}
	money := u.ResetData.AwardList[idx].AwardMoney
	u.SaveDbFlg = true
	return money, nil
}

func (u *User) SetAwardStepFinish(idx int) {
	u.ResetData.AwardList[idx].Status = define.STATUS_FINISH
	u.SaveDbFlg = true
}

func (u *User) AddHelpFriend(unionId, avatarUrl string) bool {
	u.ResetData.Friend.Lock()
	defer u.ResetData.Friend.Unlock()

	if _, ok := u.ResetData.Friend.FriendList[unionId]; ok {
		return false
	}
	u.ResetData.Friend.FriendHeadIcon = append(u.ResetData.Friend.FriendHeadIcon, avatarUrl)
	u.ResetData.Friend.HelpCnt += 1
	u.ResetData.Friend.FriendList[unionId] = avatarUrl
	u.Save()
	return true
}

func (u *User) GetHelpFriendCnt() int {
	u.ResetData.Friend.RLock()
	defer u.ResetData.Friend.RUnlock()

	return u.ResetData.Friend.HelpCnt
}

func (u *User) SubHelpFriendCnt() {
	u.ResetData.Friend.Lock()
	defer u.ResetData.Friend.Unlock()

	u.ResetData.Friend.HelpCnt -= 1
	u.SaveDbFlg = true
}

func (u *User) AddKeepDay(step int) {
	dayno := utils.GetRealDayNo()
	if u.SportDayNo == dayno {
		if step-u.SportInfo.Step > 0 {
			u.SportInfo.TotalStep += step - u.SportInfo.Step
		}
		u.SportInfo.Step = step
		return
	}
	//断开两天归零
	if dayno-u.SportDayNo >= 2 {
		u.SportInfo.KeepDay = 0
	}
	if step < 3000 {
		return
	}
	u.SportDayNo = dayno
	u.SportInfo.Step = step
	u.SportInfo.TotalStep += step
	u.SportInfo.KeepDay += 1
	u.CheckKeepRedPack(u.SportInfo.KeepDay)
	u.SaveDbFlg = true
}

func (u *User) CheckKeepRedPack(day int) {
	if day < 5 {
		return
	}
	if u.SportInfo.KeepRedPack[1] != define.STATUS_INIT {
		return
	}
	if day >= 5 && u.SportInfo.KeepRedPack[0] == define.STATUS_INIT {
		u.SportInfo.KeepRedPack[0] = define.STATUS_ALREDY
		u.SetKeepRedpack(0, define.STATUS_ALREDY)
	}
	if day >= 14 && u.SportInfo.KeepRedPack[1] == define.STATUS_INIT {
		u.SetKeepRedpack(1, define.STATUS_ALREDY)
	}
}

func (u *User) SetKeepRedpack(idx int, status int8) {
	u.SportInfo.KeepRedPack[idx] = status
}

func (u *User) AddSportLike() {
	u.SportInfo.Like += 1
}

func (u *User) SubSportLike() {
	if u.SportInfo.Like <= 0 {
		return
	}
	u.SportInfo.Like -= 1
}
