package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"e.coding.net/anyun-cloud-api-gateway/_elasticsearch"
	"e.coding.net/anyun-cloud-api-gateway/_influxdb"
	"e.coding.net/anyun-cloud-api-gateway/common"
	"e.coding.net/anyun-cloud-api-gateway/connection"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//Statistics -- 网关统计信息
//* 运行时间统计
//* 访问统计信息
//* API部署情况统计
type Statistics struct {
	Timer struct {
		Uptime         int64 //当前API网关节点的启动时间
		LastUpdateTime int64 //当前API网关最后一次更改配置的重启时间
	} //网关时间指标
	Counter struct {
		AccessTotalCount int            `json:"total"`   //系统启动后的连接总量
		ProcessCount     int            `json:"process"` //当前处理器个数
		APIAccessCount   map[string]int `json:"access"`  //根据URL分类的API访问总量
	}
	API struct {
		DeploySuccessCount        int            //已成功部署的API个数
		DeployErrorCount          int            //未成功部署的API个数
		LastAPIDeployTime         int64          //API最后一次部署的UNIX时间
		TopTimeAPI                map[string]int //最长调用时间的API信息
		AuthenticationFailedCount map[string]int //认证失败的API调用信息
	}
}

var totalCount int = 5

//SetUptime -- 设置API网关进程第一次启动的时间
func (_this *Statistics) SetUptime() {
	if _this.Timer.Uptime == 0 { //启动时间只能设置一次
		return
	}
	now := time.Now().Unix()
	_this.Timer.Uptime = now
	_this.Timer.LastUpdateTime = now
}

//UpdateConfigTime -- 更新最后一次API网关配置修改的时间
func (_this *Statistics) UpdateConfigTime() {
	now := time.Now().Unix()
	_this.Timer.LastUpdateTime = now
}

//APIStatisticsMiddleware -- 统计中间件
func (_this *Statistics) APIStatisticsMiddleware(ctx context.Context, gateway *AnyunCloudGateway, gatewayclient *connection.GatewayClientContext, params *common.APICONTROLLERPARAMS) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: 需要检测API是否在缓存里面存在
		//* 不存在的API不需要作为API调用统计
		//* 不在的API也不算作总的API调用个数里面
		//* 不存在的API记录在另外的统计
		t1 := time.Now()
		_this.Counter.AccessTotalCount++
		_this.Counter.ProcessCount++
		common.APICount = &_this.Counter.ProcessCount
		go func() {
			if _this.Counter.ProcessCount >= totalCount {
				r, _ := json.Marshal(&common.APICONTROLLERPARAMS{
					TypeDO:        "0",
					ID:            params.ID,
					Name:          params.Name,
					Version:       params.Version,
					Dc:            params.Dc,
					ConsulAddress: params.ConsulAddress,
					DiscoveryDNS:  params.DiscoveryDNS,
					Hostcategory:  params.Hostcategory,
					Blockcategory: params.Blockcategory,
					Round:         params.Round,
					NetWorkID:     params.NetWorkID,
				})
				mess, err := gatewayclient.NatsBusiness.Request(common.APIContainerManager, []byte(r), 1000*time.Millisecond)
				if err != nil {
					log.Error(err.Error())
				}
				if mess != nil && mess.Data != nil {
					log.Info(mess.Data)
				} else {
					log.Error("nats return error")
				}
			}
		}()
		if _this.Counter.APIAccessCount == nil {
			_this.Counter.APIAccessCount = map[string]int{}
		}
		if _, exist := _this.Counter.APIAccessCount[c.Request.URL.Path]; !exist {
			_this.Counter.APIAccessCount[c.Request.URL.Path] = 1
		} else {
			_this.Counter.APIAccessCount[c.Request.URL.Path]++
		}
		log.WithFields(log.Fields{
			"prefix":  "server.APIStatisticsMiddleware",
			"api-url": c.Request.URL,
		}).Debug("waiting for other resource")
		c.Next() //! 等待API处理器完成
		_this.Counter.ProcessCount--
		t2 := time.Now()
		log.WithFields(log.Fields{
			"prefix":            "server.APIStatisticsMiddleware",
			"api-total-count":   _this.Counter.AccessTotalCount,
			"api-process-count": _this.Counter.ProcessCount,
			"api-url":           c.Request.URL,
			"api-current-count": _this.Counter.APIAccessCount[c.Request.URL.Path],
			"time-comparison":   t2.Sub(t1),
		}).Debug("api counter")
		go func() {
			influx := &_influxdb.InfluxDBStructs{
				DataBase:        gatewayclient.InfluxDB,
				RetentionPolicy: gatewayclient.InfluxRetentionPolicy,
				Tags: map[string]string{
					"type": "api",
					"url":  c.Request.URL.Path,
				},
				Table: "api_gateway",
				Fields: map[string]interface{}{
					"exec_time": t2.Sub(t1),
				},
				Client: gatewayclient.InfluxDBClient,
			}
			influx.InsertDB()
		}()
	}
}

//Initialization -- 初始化内置管理API
//* API访问统计
func (_this *AnyunCloudGateway) Initialization() {
	gm1 := _this.Engine.Group("/gateway/v1/")
	{
		gm1.GET(_this.statistics())
	}
}

//statistics -- API访问统计服务
//* 返回API调用次数统计信息
func (_this *AnyunCloudGateway) statistics() (string, func(c *gin.Context)) {
	return "/statistics", func(c *gin.Context) {
		c.JSON(http.StatusOK, _this.Statistics.Counter)
	}
}

func doCheck(queryClient *connection.GatewayClientContext, requestDetail *common.RequestDetail, key string) (string, string, string, error) {
	// 查询API对应的ID
	mess, err := queryClient.CheckQueryURL(requestDetail)
	if err != nil {
		return "", "FAILED", "", err
	}
	var p common.ReponseDetail
	json.Unmarshal([]byte(mess), &p)
	if len(p.ServiceID) == 0 {
		return "无服务", "FAILED", "", nil
	}
	val, err := queryClient.CheckFromEtcd(key + p.ServiceID)
	if val == "" {
		return "服务未绑定", "FAILED", "", nil
	}
	result, err := queryClient.NatsBusiness.Request(val, buildMessage(requestDetail), 1000*time.Millisecond)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error")
		return "nats request error", "FAILED", "", err
	}
	log.Info("服务返回数据为 ", string(result.Data))
	var f interface{}
	err = json.Unmarshal(result.Data, &f)
	m := f.(map[string]interface{})
	if m["code"].(float64) == 200 {
		jsonString, _ := json.Marshal(m["result"].(map[string]interface{}))
		node := m["header"].(map[string]string)
		return string(jsonString), "SUCCESS", node["node"], nil
	}
	return "nats request error", "FAILED", "", nil
}

func buildMessage(request *common.RequestDetail) (data []byte) {
	header := &common.MessageHeader{
		Application: "api-gateway",
		Time:        time.Now().Unix(),
		Type:        "req",
		Version:     "1.0.0",
	}
	requestMessage := &common.RequestMessage{
		MessageHeader: *header,
		Business:      "service",
		Content:       request.Params,
	}
	data, _ = json.Marshal(requestMessage)
	return
}

func insertRecode(gatewayclient *connection.GatewayClientContext, bl, url, val string, c *gin.Context) {
	go func() {
		_elasticsearch.InsertESDB(
			gatewayclient.Context,
			gatewayclient.ESClient,
			&common.ElasticSearchParam{
				UserName:  "system",
				APIName:   c.Request.URL.Path,
				StartTime: time.Now(),
				Result:    bl,
			})
	}()
	go func() {
		influx := &_influxdb.InfluxDBStructs{
			DataBase:        gatewayclient.InfluxDB,
			RetentionPolicy: gatewayclient.InfluxRetentionPolicy,
			Tags: map[string]string{
				"type": "api",
				"url":  url,
			},
			Table: "api_gateway",
			Fields: map[string]interface{}{
				"exec_tx":     len(c.Request.URL.Query()),
				"exec_rx":     len(val),
				"exec_result": bl,
			},
			Client: gatewayclient.InfluxDBClient,
		}
		influx.InsertDB()
	}()
}
