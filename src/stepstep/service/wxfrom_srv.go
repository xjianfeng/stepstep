package service

import (
	"stepstep/define"
	"stepstep/models"
	"strings"
)

type ReqWxFormId struct {
	FormId string `json:"formId" binding:"required"`
}

func AddWxFormId(unionId string, args *ReqWxFormId) error {
	if args.FormId == "" {
		return define.ERROR_REQUEST_PARAMS
	}
	u, err := GetUser(unionId)
	if err != nil {
		return err
	}
	if strings.Contains(args.FormId, " ") {
		return nil
	}
	models.AddRedisFormId(unionId, u.OpenId, args.FormId)
	return nil
}
