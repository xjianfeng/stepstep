package service

import (
	"github.com/xjianfeng/gocomm/logger"
	"gopkg.in/mgo.v2"
	"stepstep/define"
	"stepstep/models"
)

func GetUnionIdByUid(uid int) (string, error) {
	unionId, err := models.LoadUnionIdByUid(uid)
	return unionId, err
}

func GetUserData(unionId, openId string) (*models.User, error) {
	u, err := models.LoadUser(unionId)
	if err != nil && err != mgo.ErrNotFound {
		logger.LogError("%s", err.Error())
		return nil, err
	}

	imageCover := false
	if err == mgo.ErrNotFound {
		if openId == "" {
			return nil, define.ERROR_USER_NOT_EXISTS
		}
		u, err = models.Newuser(unionId, openId)
		if err != nil {
			logger.LogError("%s", err.Error())
			return nil, err
		}
		imageCover = true
		u.Save()
	}
	if openId != "" {
		imageCover = u.ResetEveryDay()
	}
	//创建用户，和每日重置数据时设置封面头像
	if imageCover {
		SetSportCover(u.UnionId)
	}
	return u, nil
}

func GetUser(unionId string) (*models.User, error) {
	return GetUserData(unionId, "")
}

func GetRefreshUser(unionId string) (*models.User, error) {
	u, err := GetUserData(unionId, "")
	if err != nil {
		return nil, err
	}
	u.CheckTimeout()
	return u, err
}
