package server

import (
	"regexp"

	"e.coding.net/anyun-cloud-api-gateway/common"
	"e.coding.net/anyun-cloud-api-gateway/connection"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//FileMiddleware -- 网关只匹配URI以"/FILE/"开头的调用为应用的API调用
func (_this *AnyunCloudGateway) FileMiddleware(gatewayclient *connection.GatewayClientContext) gin.HandlerFunc {
	BuiltInAPIRegexp := "^/file/"
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		if matched, _ := regexp.MatchString(BuiltInAPIRegexp, url); matched {
			log.WithFields(log.Fields{
				"prefix":         "server.FILEMiddleware",
				"request-path":   url,
				"request-method": c.Request.Method,
				"client-ip":      c.ClientIP(),
				"params":         c.Request.URL.Query(),
				"content-type":   c.ContentType(),
				"content-length": c.Request.ContentLength,
				"accept-type":    c.Request.Header.Get("Accept"),
			}).Info("request info")
			val, bl, header, err := doCheck(gatewayclient, &common.RequestDetail{
				BUSSINESS:   "getByteVerification",
				URL:         url,
				Method:      c.Request.Method,
				Params:      c.Request.URL.Query(),
				ContentType: c.ContentType(),
				AcceptType:  c.Request.Header.Get("Accept"),
			}, "/container_file/")

			if err != nil {
				log.WithFields(log.Fields{
					"prefix": "server.FileMiddleware.docheck",
				}).Error("error")
				c.JSON(SYSTEMERROR, gin.H{
					"message": "error",
				})
			} else {
				c.Header("APIID", gatewayclient.RandomNumber)
				c.Header("container-node", header)
				c.String(OK, val)
			}
			insertRecode(gatewayclient, bl, url, val, c)
		}
	}
}
