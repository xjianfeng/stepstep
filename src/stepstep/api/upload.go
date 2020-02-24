package api

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/xjianfeng/gocomm/decry"
	"stepstep/conf"
	"stepstep/define"
	"stepstep/package/cos"
	"stepstep/service"
	"strconv"
	"strings"
)

var (
	imgExts = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".git":  true,
		".bmg":  true,
	}
	soundExts = map[string]bool{
		".mp3": true,
	}
)

/**
* showdoc
* @catalog
* @title 图片上传
* @description 图片上传接口
* @method post
* @url /step/upload
* @return {"code":0,"errMsg":"","data":{"image":"http://xxx"}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data json 返回数据
* @return_param data->image string 上传图片返回的路径
* @remark 这里是备注信息
* @number 99
 */
func ApiUploadImg(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	file, _ := c.FormFile("file")
	fileName := file.Filename
	if len(fileName) < 4 {
		ret.Err = define.ERROR_REQUEST_PARAMS
		return ret
	}
	eIdx := strings.LastIndex(fileName, ".")
	fileExt := fileName[eIdx:]
	if !imgExts[fileExt] {
		ret.Err = define.ERROR_REQUEST_PARAMS
		return ret
	}
	uuv4, err := uuid.NewV4()
	if err != nil {
		ret.Err = define.ERROR_SERVER_ERROR
		return ret
	}

	fileName = decry.Md5Sum(uuv4.Bytes())
	f, err := file.Open()
	if err != nil {
		ret.Err = err
		return ret
	}
	defer f.Close()

	fileName = conf.CfgCos.ImagePath + fileName + fileExt
	ret.Err = cos.UploadObject(fileName, f)
	if ret.Err != nil {
		return ret
	}
	image := conf.CfgCos.Domain + fileName
	service.SetSportImage(unionId, image)
	ret.Body = gin.H{
		RET_KEY_IMAGE: image,
	}
	return ret
}

/**
* showdoc
* @catalog
* @title 声音上传
* @description 声音上传接口
* @method post
* @url /step/upload/sound
* @return {"code":0,"errMsg":"","data":{"sound":"http://xxx"}}
* @return_param code int 返回状态码0为成功,其他为错误
* @return_param errMsg string 错误信息,当code不为0时才有
* @return_param data json 返回数据
* @return_param data->sound string 上传mp3后的地址
* @remark 这里是备注信息
* @number 99
 */
func ApiUploadSound(c *gin.Context) *ApiRetData {
	ret := NewJsonRet()
	unionId, err := GetUnionId(c)
	if err != nil {
		ret.Err = err
		return ret
	}
	file, _ := c.FormFile("file")
	soundSec, err := strconv.Atoi(c.PostForm("length"))
	if err != nil {
		ret.Err = err
		return ret
	}
	fileName := file.Filename
	if len(fileName) < 4 {
		ret.Err = define.ERROR_REQUEST_PARAMS
		return ret
	}
	eIdx := strings.LastIndex(fileName, ".")
	fileExt := fileName[eIdx:]
	if !soundExts[fileExt] {
		ret.Err = define.ERROR_REQUEST_PARAMS
		return ret
	}
	uuv4, err := uuid.NewV4()
	if err != nil {
		ret.Err = define.ERROR_SERVER_ERROR
		return ret
	}

	fileName = decry.Md5Sum(uuv4.Bytes())
	f, err := file.Open()
	if err != nil {
		ret.Err = err
		return ret
	}
	defer f.Close()
	fileName = conf.CfgCos.SoundPath + fileName + fileExt
	ret.Err = cos.UploadObject(fileName, f)
	if ret.Err != nil {
		return ret
	}
	sound := conf.CfgCos.Domain + fileName
	service.SetSportSound(unionId, &service.ReqSportSound{sound, soundSec})
	ret.Body = gin.H{
		RET_KEY_SOUND: sound,
	}
	return ret
}
