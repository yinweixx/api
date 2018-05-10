package server

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Gateway -- API网关服务器接口
type Gateway interface {
	Start() error   //启动网关服务
	Stop() error    //停止网关服务
	Restart() error //重启网关服务
	Join() error    //维持API网关运行
}

//AnyunCloudGateway -- API网关服务器结构体
//
//* 使用etcd v3客户端获取配置API网关配置
//* 网关配置监听器将监听etcd的网关配置，动态的调整网关配置
//* API监听器监听网关的API变动，实时的部署或者卸载API
type AnyunCloudGateway struct {
	Config     *AnyunCloudGatewayConfig //API网关服务器配置
	Engine     *gin.Engine
	Server     *http.Server
	Statistics *Statistics
	// discovery  *discovery.QueryClient
	ctx context.Context
}

//AnyunCloudGatewayConfig -- API网关服务器配置结构体
type AnyunCloudGatewayConfig struct {
	ListenerAddr string //API网关监听地址
	HTTPS        struct {
		TLSConfig  *tls.Config //HTTPS TLS配置
		TLSCert    string      //TLS证书
		TLSCertKey string      //TLS私钥
	}
}

//AnyunCloudGatewayContext -- API网关服务器上下文
type AnyunCloudGatewayContext struct {
}

//API -- 平台API描述
type API struct {
	Group   string   // API分组信息为应用程序的英文短名称
	Version string   // 版本信息为API在etcd的配置索引版本信息(每次修改在etcd产生的索引版本信息)
	URI     string   //API的URI信息
	Method  []string //API支持的HTTP方法(GET,POST,PUT,PATCH,OPTION)
}
