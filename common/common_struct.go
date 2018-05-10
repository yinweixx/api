package common

import "time"

var (
	// CheckQueryURL -- 检查请求地址
	CheckQueryURL = "API_CONTROLLER_CHANNEL_TEST"
	// APIContainerManager --
	APIContainerManager = "API_CONTROLLER_MANAGER"
	//REDISTIMETOLIVE --
	REDISTIMETOLIVE = 600
)

/*
APICONTROLLERPARAMS -- api 控制器请求实体
*/
type APICONTROLLERPARAMS struct {
	TypeDO        string `json:"typedo"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	Dc            string `json:"dc"`
	Hostcategory  string `json:"hostcategory"`
	Blockcategory string `json:"blockcategory"`
	Dnsname       string `json:"dnsname"`
	Round         string `json:"round"`
	DiscoveryDNS  string `json:"DNS"`
	ConsulAddress string `json:"consuladdress"`
	NetWorkID     string `json:"networkid"`
}

//RequestDetail -- 请求数据实体
type RequestDetail struct {
	BUSSINESS   string              `json:"BUSINESS"`
	URL         string              `json:"url"`
	Method      string              `json:"method"`
	Params      map[string][]string `json:"params"`
	ContentType string              `json:"contenttype"`
	AcceptType  string              `json:"accepttype"`
}

//ReponseDetail -- 返回数据实体
type ReponseDetail struct {
	URL       string `json:"url"`
	ServiceID string `json:"serviceId"`
}

//NatsRequestDetail -- Nats数据请求实体
type NatsRequestDetail struct {
	channel string
	data    []byte
}

//MetaDataVersion -- 元数据版本信息结构体
type MetaDataVersion struct {
	//TODO: 元数据版本信息定义
	serviceBranch string
	tagVersion    string
	info          string
}

//EndpointAPIMetaData -- API元数据信息结构体
type EndpointAPIMetaData struct {
	//TODO: API元数据信息定义
	Service EndpointServiceMetaData
	Version MetaDataVersion
}

//EndpointServiceMetaData -- 服务元数信息据结构体
//* 如果服务支持缓存，必须设置缓存的TTL
type EndpointServiceMetaData struct {
	//TODO: 服务元数据信息定义
	projectName  string
	serviceName  string
	Version      MetaDataVersion
	SupportCache bool //服务是否支持缓存
	CacheTTL     bool //服务缓存的TTL
}

//DNSTypeSRVRecord -- DNS SRV记录
type DNSTypeSRVRecord struct {
	//TODO: DNS SRV记录定义
	SRV []struct {
		FQDN string //提供服务节点的域名信息
		IP   string //提供服务节点的IP地址信息
		Port int    //提供服务节点的端口信息
	}
}

//ServerConnectionInfo -- 服务器连接信息
type ServerConnectionInfo struct {
	Address string //URI,IP,ADDR
	Port    int    //服务端口
}

//MessageHeader -- 请求消息头
type MessageHeader struct {
	Version     string `json:"version"`
	Type        string `json:"type"`
	Application string `json:"application"`
	Time        int64  `json:"time"`
}

//RequestMessage -- 请求消息实体
type RequestMessage struct {
	MessageHeader MessageHeader `json:"header"`
	Business      string        `json:"business"`
	Content       interface{}   `json:"content"`
}

type header struct {
	Encryption  string `json:"encryption"`
	Timestamp   int64  `json:"timestamp"`
	Key         string `json:"key"`
	Partnercode int    `json:"partnercode"`
}

//Event -- 平台事件结构体
type Event struct {
	From     string //事件来源URI
	Describe string //事件的详细信息,JSON
}

//ETCDResp --
type ETCDResp struct {
	Name   string `json:"name"`
	Uptime string `json:"uptime"`
	Lang   string `json:"lang"`
	Info   string `json:"info"`
}

//APICount --
var APICount *int

//ElasticSearchParam --
type ElasticSearchParam struct {
	UserName  string    `json:"username"`
	APIName   string    `json:"apiname"`
	StartTime time.Time `json:"starttime"`
	Result    string    `json:"result"`
}

//ElasticSearchParam2 --
type ElasticSearchParam2 struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	DC                 string    `json:"dc"`
	Version            string    `json:"version"`
	CreateTime         time.Time `json:"time"`
	ContainerIPAddress string    `json:"containeripaddress"`
	NetworkID          string    `json:"networkid"`
	Type               string    `json:"type"`
	Result             string    `json:"result"`
}
