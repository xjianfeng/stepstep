package service

import (
	"math/rand"
	"stepstep/define"
	"stepstep/models"
)

type ReqHelpFriend struct {
	UnionId string `json:"unionId" binding:"required"`
}

//帮助好友
func SetHelpFriend(unionId string, args *ReqHelpFriend) error {
	u, err := GetUser(unionId)
	if err != nil {
		return err
	}
	helpUnionId := args.UnionId
	if unionId == args.UnionId {
		return define.ERROR_HELP_FAIL
	}
	tar, err := GetUser(helpUnionId)
	if err != nil {
		return err
	}
	avatarUrl := u.AvatarUrl
	tar.AddHelpFriend(unionId, avatarUrl)
	return nil
}

type RetFriendRedPack struct {
	U          *models.User
	AwardMoney float64
}

//好友互助红包
func GetFriendHelpRedpack(unionId string) (*RetFriendRedPack, error) {
	u, err := GetUser(unionId)
	if err != nil {
		return nil, err
	}
	if u.GetFriendHelpCnt() <= 0 {
		return nil, define.ERROR_AWARD_NOT_EXISTS
	}
	m := float64(rand.Intn(define.AWARD_FRIEND_REDPACK)+1) / 10.0
	u.AddMoney(m, define.SYSTEM_FIREND, define.DESC_TYPE_TASK)
	u.SubHelpFriendCnt()
	ret := &RetFriendRedPack{
		U:          u,
		AwardMoney: m,
	}
	return ret, nil
}
