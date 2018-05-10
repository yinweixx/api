package server

import (
	"regexp"

	"e.coding.net/anyun-cloud-api-gateway/common"
	"e.coding.net/anyun-cloud-api-gateway/connection"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//BYTEMiddleware -- 中间件
func (_this *AnyunCloudGateway) BYTEMiddleware(gatewayclient *connection.GatewayClientContext) gin.HandlerFunc {
	BuiltInAPIRegexp := "^/byte/"
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		if matched, _ := regexp.MatchString(BuiltInAPIRegexp, url); matched {
			log.WithFields(log.Fields{
				"prefix":         "server.BYTEMiddleware",
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
			}, "/container_byte/")

			if err != nil {
				log.WithFields(log.Fields{
					"prefix": "server.ByteMiddleware.docheck",
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
