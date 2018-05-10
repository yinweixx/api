package config

import (
	"e.coding.net/anyun-cloud-api-gateway/app"
	"github.com/urfave/cli"
)

//GatewayArgs -- API网关启动配置映射结构体
//
//* 日志的级别的配置
//* 分布式DNS服务器的配置
type GatewayArgs struct {
	LoggerLevel      string
	DNS              string `yaml:"dns"`
	DataCenter       string `yaml:"datacenter"`
	EthernetName     string `yaml:"ethernet_name"`
	ConsulLeaderName string `yaml:"consul_leader_name"`
	HostCategory     string `yaml:"host_category"`
	BlockCategory    string `yaml:"block_category"`
	NatsManager      string `yaml:"natsmanager"`
	NatsBusiness     string `yaml:"natesbusiness"`
	Etcd             string `yaml:"etcd"`
	Redis            string `yaml:"redis"`
	Name             string `yaml:"name"`
	ImageVersion     string `yaml:"imageversion"`
	Round            string `yaml:"round"`
	NetWorkID        string `yaml:"networkid"`
	ConsulAddress    string `yaml:"consuladdress"`
	ESURL            string `yaml:"esurl"`
	ESUserName       string `yaml:"esusername"`
	ESPassWord       string `yaml:"espassword"`
}

var args GatewayArgs

//GatewayInitFlags -- 网关启动参数初始化映射
func GatewayInitFlags() *app.ApplicationArgs {
	return &app.ApplicationArgs{
		Name:    "anyun-cloud-api-gateway",
		Usage:   "Anyun Cloud Distributed API Gateway",
		Version: "1.0.0",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "logger_level",
				Usage:       "logger level",
				Value:       "debug",
				Destination: &args.LoggerLevel,
			},
			cli.StringFlag{
				Name:        "discovery_dns",
				Usage:       "discovery dns list (eg. 192.168.1.1,192.168.1.2)",
				Destination: &args.DNS,
			},
			cli.StringFlag{
				Name:        "datacenter",
				Usage:       "data center name",
				Value:       "dc-anyuncloud",
				Destination: &args.DataCenter,
			},
			cli.StringFlag{
				Name:        "consul_leader_name",
				Usage:       "consul leader dns name",
				Destination: &args.ConsulLeaderName,
			},
			cli.StringFlag{
				Name:        "bind_eth",
				Usage:       "bind ethernet name",
				Value:       "eth1",
				Destination: &args.EthernetName,
			},
			cli.StringFlag{
				Name:        "container_type",
				Usage:       "container type,宿主机类型，例如docker,lxd,聚合",
				Value:       "docker",
				Destination: &args.HostCategory,
			},
			cli.StringFlag{
				Name:        "block_category",
				Usage:       "block category,IP段用途，例如 管理，宿主机，业务，网关",
				Value:       "网关",
				Destination: &args.BlockCategory,
			},
			cli.StringFlag{
				Name:        "nats_manager_address",
				Usage:       "nats manager address",
				Destination: &args.NatsManager,
			},
			cli.StringFlag{
				Name:        "nats_business_address",
				Usage:       "nats business address",
				Destination: &args.NatsBusiness,
			},
			cli.StringFlag{
				Name:        "etcd_address",
				Usage:       "etcd address",
				Destination: &args.Etcd,
			},
			cli.StringFlag{
				Name:        "redis_address",
				Usage:       "redis address",
				Destination: &args.Redis,
			},
			cli.StringFlag{
				Name:        "name",
				Usage:       "container images name，镜像名称",
				Destination: &args.Name,
			},
			cli.StringFlag{
				Name:        "images_version",
				Usage:       "container images version",
				Destination: &args.ImageVersion,
			},
			cli.StringFlag{
				Name:        "round",
				Usage:       "work conditions,for example:测试，正式",
				Value:       "正式",
				Destination: &args.Round,
			},
			cli.StringFlag{
				Name:        "networkid",
				Usage:       "docker network", //this is container connect to enps1 which ovs network
				Value:       "eb79ff1f13be",
				Destination: &args.NetWorkID,
			},
			cli.StringFlag{
				Name:        "consul_address",
				Usage:       "consul address,for example:etcd.service.consul",
				Destination: &args.ConsulAddress,
			},
			cli.StringFlag{
				Name:        "es_url",
				Usage:       "elasticsearch url",
				Destination: &args.ESURL,
			},
			cli.StringFlag{
				Name:        "es_username",
				Usage:       "elasticsearch username",
				Destination: &args.ESUserName,
			},
			cli.StringFlag{
				Name:        "es_password",
				Usage:       "elasticsearch password",
				Destination: &args.ESPassWord,
			},
		},
	}
}
