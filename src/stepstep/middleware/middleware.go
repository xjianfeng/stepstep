package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"github.com/xjianfeng/gocomm/logger"
	"io/ioutil"
	"net/http"
	"stepstep/api"
	"stepstep/define"
	"stepstep/models"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var ignorePath = "/step/wechat/login"

var errResponse = `{"code":%d, "errMsg":"%s", "data":{}}`
var normalResponse = `{"code":0, "errMsg":"", "data":%s}`

func abortRequest(c *gin.Context) {
	c.String(200, errResponse, 500, "请登录")
	c.Abort()
}

func CustomMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		header := c.GetHeader("Content-Type")
		unionId := c.GetHeader("unionId")
		token := c.GetHeader("token")

		//校验token是否合法
		if path != ignorePath {
			if token == "" {
				abortRequest(c)
				return
			}
			serverToken := models.GetRedisUInfo(unionId, define.REDIS_HEKY_USER_TOKEN)
			if string(serverToken) != token {
				abortRequest(c)
				return
			}
		}
		//打印请求数据
		if strings.ToLower(header) == "application/json" {
			body, _ := c.GetRawData()
			logger.LogInfo("[Request] Path:%s, UnionId:%s, Content-Type:%s, body:%s", path, unionId, header, body)

			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		} else {
			logger.LogInfo("[Request] Path:%s, UnionId:%s, Content-Type:%s", path, unionId, header)
		}
		c.Next()
	}
}

func DoSave(data *api.ApiRetData) {
	if data.Err != nil {
		return
	}
	if data.U == nil {
		return
	}
	if !data.U.SaveDbFlg {
		return
	}
	data.U.Save()
}

func DecorateFunc(fun func(*gin.Context) *api.ApiRetData) gin.HandlerFunc {
	var err error
	var respBody string
	var respType = define.RESPONSE_TYPE_STRING
	var respData interface{}

	return func(c *gin.Context) {
		ret := fun(c)
		err = ret.Err
		respType = ret.RespType
		respData = ret.Body

		defer func() {
			DoSave(ret)
			path := c.Request.URL.Path
			logger.LogInfo("[Response] Path:%s, err:%v, respType:%s, respData:%s", path, err, respType, respBody)
		}()
		if err != nil {
			c.String(http.StatusOK, errResponse, 500, err.Error())
			return
		}
		if respType == define.RESPONSE_TYPE_STRING {
			respBody = respData.(string)
			c.String(http.StatusOK, respBody)
			return
		}
		if respType == define.RESPONSE_TYPE_JSON {
			byteData, err := json.Marshal(respData)
			respBody = string(byteData)
			if err != nil {
				c.String(http.StatusOK, errResponse, 500, err.Error())
				return
			}
			c.String(http.StatusOK, normalResponse, respBody)
			return
		}
	}
}
