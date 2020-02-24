package service

import (
	"math/rand"
	"stepstep/define"
	"stepstep/models"
)

type ReqAwardArgs struct {
	AwardType string `json:"awardType" binding:"required"`
}

type AwardRet struct {
	U           *models.User
	AwardStep   int
	AwardMoney  float64
	NextTimeout int64
}

func GetAwardByType(unionId string, args *ReqAwardArgs) (*AwardRet, error) {
	ret := &AwardRet{}
	u, err := GetUser(unionId)
	if err != nil {
		return nil, err
	}
	ret.U = u
	if args.AwardType == define.TIMEOUT_KEY_VIDEO {
		err := awardVideoStep(unionId, ret)
		return ret, err
	}
	if args.AwardType == define.TIMEOUT_KEY_LUCKSTEP {
		err := awardLuckStep(unionId, ret)
		return ret, err
	}
	if args.AwardType == define.TIMEOUT_KEY_REDPACK {
		err := awardRedPack(unionId, ret)
		return ret, err
	}
	return nil, define.ERROR_AWARD_NOT_EXISTS
}

func awardLuckStep(unionId string, ret *AwardRet) error {
	u := ret.U
	if !u.CheckTimeOutStatus(define.TIMEOUT_KEY_LUCKSTEP) {
		return define.ERROR_AWARD_FAIL
	}
	u.AddAwardStep(define.AWARD_LUCK_STEP)
	ret.AwardStep = define.AWARD_LUCK_STEP
	u.SetTimeOutFinish(define.TIMEOUT_KEY_LUCKSTEP)
	u.CheckTaskFinish()
	return nil
}

func awardVideoStep(unionId string, ret *AwardRet) error {
	u := ret.U
	if !u.CheckTimeOutStatus(define.TIMEOUT_KEY_VIDEO) {
		return define.ERROR_AWARD_FAIL
	}
	u.AddAwardStep(define.AWARD_VIDEO_STEP)
	ret.AwardStep = define.AWARD_VIDEO_STEP
	u.SetTimeOutFinish(define.TIMEOUT_KEY_VIDEO)
	u.CheckTaskFinish()
	return nil
}

func awardRedPack(unionId string, ret *AwardRet) error {
	u := ret.U
	if !u.CheckTimeOutStatus(define.TIMEOUT_KEY_REDPACK) {
		return define.ERROR_AWARD_FAIL
	}
	i := float64(rand.Intn(define.AWARD_RED_PACK)+1) / 10.0
	u.AddMoney(i, define.SYSTEM_TIMEOUT_REDPACK, define.DESC_TYPE_TIMEOUT_REDPACK)
	ret.AwardMoney = i
	u.SetTimeOutFinish(define.TIMEOUT_KEY_REDPACK)
	tdata := u.GetTimeoutData(define.TIMEOUT_KEY_REDPACK)
	if tdata != nil {
		ret.NextTimeout = tdata.TimeOut
	}
	return nil
}

type ReqStepAward struct {
	Idx int `json:"idx"`
}

type RetStepAward struct {
	U     *models.User
	Money float64
}

func GetStepTaskAward(unionId string, args *ReqStepAward) (*RetStepAward, error) {
	u, err := GetUser(unionId)
	if err != nil {
		return nil, err
	}
	idx := args.Idx
	money, err := u.GetStepAwardMoney(idx)
	if err != nil {
		return nil, err
	}
	u.AddMoney(money, define.SYSTEM_TASK, define.DESC_TYPE_TASK)
	u.SetAwardStepFinish(idx)
	ret := &RetStepAward{
		U:     u,
		Money: money,
	}
	return ret, nil
}
