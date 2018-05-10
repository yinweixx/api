package connection

import (
	"context"
	"math/rand"
	"time"

	"e.coding.net/anyun-cloud-api-gateway/_elasticsearch"
	"e.coding.net/anyun-cloud-api-gateway/pool"
	"github.com/coreos/etcd/clientv3"
	"github.com/garyburd/redigo/redis"
	influxDBClient "github.com/influxdata/influxdb/client/v2"
	nats "github.com/nats-io/go-nats"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

//GatewayClientContext -- all client
type GatewayClientContext struct {
	Context               context.Context
	NatsManager           *nats.Conn
	NatsBusiness          *nats.Conn
	Etcd                  *clientv3.Client
	Redis                 *redis.Conn
	InfluxDBClient        *influxDBClient.Client
	InfluxDBURL           string
	InfluxDB              string
	InfluxRetentionPolicy string
	ESClient              *elastic.Client
	RandomNumber          string //用以证明负载均衡，随机生成
}

//GatewayClientConfig -- client address
type GatewayClientConfig struct {
	NatsManager             string
	NatsBusiness            string
	Etcd                    []string
	Redis                   string
	InfluxDB                string
	InfluxDBURL             string
	InfluxDBUser            string
	InfluxDBPass            string
	InfluxDBRetentionPolicy string
	ESURL                   string
	ESUserName              string
	ESPassword              string
}

//NewGatewayClientContext -- create new client
func NewGatewayClientContext(ctx context.Context, config *GatewayClientConfig) (*GatewayClientContext, error) {
	_config := GetConnInfo()
	if config.Etcd == nil || len(config.Etcd) == 0 || config.Etcd[0] == "" {
		config.Etcd = _config.Etcd
	}
	if config.NatsBusiness == "" {
		config.NatsBusiness = _config.NatsBusiness
	}
	if config.NatsManager == "" {
		config.NatsManager = _config.NatsManager
	}
	if config.Redis == "" {
		config.Redis = _config.Redis
	}
	if config.InfluxDBURL == "" {
		config.InfluxDBURL = _config.InfluxDBURL
	}
	if config.InfluxDB == "" {
		config.InfluxDB = _config.InfluxDB
	}
	if config.InfluxDBUser == "" {
		config.InfluxDBUser = _config.InfluxDBUser
	}
	if config.InfluxDBPass == "" {
		config.InfluxDBPass = _config.InfluxDBPass
	}
	if config.InfluxDBRetentionPolicy == "" {
		config.InfluxDBRetentionPolicy = _config.InfluxDBRetentionPolicy
	}
	if config.ESURL == "" {
		config.ESURL = _config.ESURL
	}
	if config.ESUserName == "" {
		config.ESUserName = _config.ESUserName
	}
	if config.ESPassword == "" {
		config.ESPassword = _config.ESPassword
	}
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Etcd,
		DialTimeout: 6 * time.Second,
	})
	if err != nil {
		log.Error(err.Error())
	}
	managerNatsClient := make(chan *nats.Conn)
	serviceNatsClient := make(chan *nats.Conn)
	go func() {
		log.Info("manager nats url is ", config.NatsManager)
		client, err := nats.Connect(config.NatsManager)
		if err != nil {
			log.Error(err.Error())
		} else {
			log.Info("manager nats connected")
		}
		managerNatsClient <- client
	}()
	go func() {
		log.Info("business nats url is ", config.NatsBusiness)
		client, err := nats.Connect(config.NatsBusiness)
		if err != nil {
			log.Error(err.Error())
		} else {
			log.Info("business nats connected")
		}
		serviceNatsClient <- client
	}()
	ESClient := _elasticsearch.NewClient(config.ESURL, config.ESUserName, config.ESPassword)
	c, err := influxDBClient.NewHTTPClient(influxDBClient.HTTPConfig{
		Addr:     config.InfluxDBURL,
		Username: config.InfluxDBUser,
		Password: config.InfluxDBPass,
	})
	if err != nil {
		log.Error(err)
	}
	redisClient := pool.NewRedisPool(config.Redis).Get()
	return &GatewayClientContext{
		Context:               ctx,
		Etcd:                  etcd,
		NatsBusiness:          <-serviceNatsClient,
		NatsManager:           <-managerNatsClient,
		Redis:                 &redisClient,
		InfluxDBClient:        &c,
		InfluxDB:              config.InfluxDB,
		InfluxRetentionPolicy: config.InfluxDBRetentionPolicy,
		ESClient:              ESClient,
		RandomNumber:          GetRandomString(21),
	}, nil
}

//GetRandomString -- 随机生成该api的随机字符串，用来在前端证明负载均衡
func GetRandomString(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
