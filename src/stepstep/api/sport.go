package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"stepstep/package/utils"
	"stepstep/service"
	"strconv"
)

/**
* showdoc
* @catalog 运动接口
* @title 获取运动数据
* @description 获取运动数据
* @method post
* @url /step/sport/data
* @param uId 是 int 被查看人的用户ID
* @param unionId 是 string 被查看人的unionId
* @return {"code":0,"errMsg":"","data":{"msgNum":1,"qrcode":"xxx","distance":1,"keepday":1,"like":1,"stepList":[122,122,299],"soundList":[{"url":"xx","length":1}],"imagesList":["xxx","xxx","xxx"],"nickName":"","avatarUrl":"xxx","keepRedPack":[0,0]}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data->msgNum int 留言数量
* @return_param data->distance float 公里数
* @return_param data->keepday int 保持天数
* @return_param data->like int 点赞数
* @return_param data->qrcode string 小程序码附带参数uid=xxx
* @return_param data->newmsg bool 是否有新留言
* @return_param data->stepList array 最近3天步数
* @return_param data->soundList array 最近3天录音
* @return_param data->imagesList array 最近3张图片
* @return_param data->nickName string 昵称
* @return_param data->avatarUrl string 头像
* @return_param data->keepRedPack array 打卡红包状态0(不可领),1(可领取),2(已领取)
* @remark 这里是备注信息
* @number 99
 */
func ApiSportData(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	reqUnionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	args := &service.ReqSportData{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}
	//通过UID获取unionId
	if args.Uid != 0 && args.UnionId == "" {
		args.UnionId, ret.Err = service.GetUnionIdByUid(args.Uid)
		if ret.Err != nil {
			return ret
		}
	}
	data, err := service.GetSportData(reqUnionId, args)
	if err != nil {
		ret.Err = err
		return ret
	}
	v := data.SportData
	var k string
	nowDayNo := utils.GetRealDayNo()
	//三天前
	threeDayAgo := nowDayNo - 2
	imageList := [3]string{}
	idx := 0
	for i := nowDayNo; i >= threeDayAgo; i-- {
		k = strconv.Itoa(i)
		imageList[idx] = v.ImageList[k]
		idx++
	}
	soundList := [3]gin.H{}
	idx = 0
	for i := nowDayNo; i >= threeDayAgo; i-- {
		k = strconv.Itoa(i)
		data, ok := v.SoundList[k]
		if !ok {
			soundList[idx] = gin.H{"url": "", "length": 0}
		} else {
			soundList[idx] = gin.H{
				"url":    data.Sound,
				"length": data.SoundSec,
			}
		}
		idx++
	}
	ret.Body = gin.H{
		"like":        v.Like,
		"msgNum":      data.MsgNum,
		"distance":    fmt.Sprintf("%.2f", float64(v.TotalStep)/float64(1600)),
		"qrcode":      v.QrCode,
		"newmsg":      data.NewMsg,
		"keeyday":     v.KeepDay,
		"soundList":   soundList,
		"stepList":    v.StepList,
		"imagesList":  imageList,
		"nickName":    data.NickName,
		"avatarUrl":   data.AvatarUrl,
		"keepRedPack": v.KeepRedPack,
	}
	return ret
}

/**
* showdoc
* @catalog 运动接口
* @title 运动点赞
* @description 运动点赞
* @method post
* @url /step/sport/like
* @param uId 是 int 被点赞人的用户Id
* @param unionId 是 string 被点赞人的unionId
* @param opType 是 int 0添加1取消
* @return {"code":0,"errMsg":"","data":{}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @remark 这里是备注信息
* @number 99
 */
func ApiSportLike(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	_, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	args := &service.ReqSportLike{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}
	//通过UID获取unionId
	if args.Uid != 0 && args.UnionId == "" {
		args.UnionId, ret.Err = service.GetUnionIdByUid(args.Uid)
		if ret.Err != nil {
			return ret
		}
	}

	ret.Err = service.SetSportLike(args)
	if ret.Err != nil {
		return ret
	}
	ret.Body = gin.H{}
	return ret
}

/**
* showdoc
* @catalog 运动接口
* @title 换封面
* @description 换封面
* @method post
* @url /step/sport/cover
* @return {"code":0,"errMsg":"","data":{"image":"xxx"}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data->image string 图片路径
* @remark 这里是备注信息
* @number 99
 */
func ApiSportCover(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	image, err := service.SetSportCover(unionId)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Body = gin.H{"image": image}
	return ret
}

//暂时不用
func ApiSportSound(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	args := &service.ReqSportSound{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}

	ret.Err = service.SetSportSound(unionId, args)
	if ret.Err != nil {
		return ret
	}
	ret.Body = gin.H{}
	return ret
}

/**
* showdoc
* @catalog 运动接口
* @title 领红包
* @description 领红包
* @method post
* @url /step/award/sport
* @param idx 是 int 领取的下标
* @return {"code":0,"errMsg":"","data":{"money":100,"awardMoney":0.5,"idx":0}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data->money float 当前金币
* @return_param data->awardMoney float 获得的红包数
* @return_param data->idx int 领取的下标
* @remark 这里是备注信息
* @number 99
 */
func ApiSportRedpack(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	args := &service.ReqSportRedpack{}
	ret.Err = c.BindJSON(args)
	if ret.Err != nil {
		return ret
	}
	result, err := service.GetSportRedpack(unionId, args)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.U = result.U
	ret.Body = gin.H{
		"idx":        args.Idx,
		"money":      result.U.Money,
		"awardMoney": result.AwardMoney,
	}
	return ret
}
