package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"e.coding.net/anyun-cloud-api-gateway/_elasticsearch"
	"e.coding.net/anyun-cloud-api-gateway/common"
	"e.coding.net/anyun-cloud-api-gateway/connection"
	"e.coding.net/anyun-cloud-api-gateway/server"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	hashids "github.com/speps/go-hashids"
	"github.com/urfave/cli"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

//Bootstrap -- 网关启动入口函数
func Bootstrap(ctx *cli.Context) error {
	getEnv()
	logLevel, _ := log.ParseLevel(args.LoggerLevel)
	_ctx := context.Background()
	gateWayClient, _ := connection.NewGatewayClientContext(
		_ctx,
		&connection.GatewayClientConfig{
			NatsManager:  args.NatsManager,
			NatsBusiness: args.NatsBusiness,
			Etcd:         []string{args.Etcd},
			Redis:        args.Redis,
			ESURL:        args.ESURL,
			ESUserName:   args.ESUserName,
			ESPassword:   args.ESPassWord,
		})
	logLevel, err := log.ParseLevel(args.LoggerLevel)
	if err != nil {
		log.WithField("prefix", "application.AgentBootstrap").Fatalf("unsuppored logger level: %s", args.LoggerLevel)
	}
	log.SetLevel(logLevel)
	formatter := new(prefixed.TextFormatter)
	log.SetFormatter(formatter)
	gateway := server.GetTestAnyunCloudGateway()
	params := &common.APICONTROLLERPARAMS{
		ID:            gateWayClient.RandomNumber,
		Name:          "anyuncloud/gateway-api:v1",
		Version:       args.ImageVersion,
		Dc:            args.DataCenter,
		Hostcategory:  args.HostCategory,
		Blockcategory: args.BlockCategory,
		Dnsname:       args.DataCenter,
		Round:         args.Round,
		ConsulAddress: args.ConsulAddress,
		DiscoveryDNS:  args.DNS,
		NetWorkID:     args.NetWorkID,
	}
	gateway.SetUpMiddlewares(_ctx, gateway, gateWayClient, params)
	gateway.Start()
	gateway.Initialization()
	server.EtcdKeepAlive(gateWayClient.Etcd)
	common.PrintFilesName("/etc/consul.d",
		common.GetIPFromNetWork(args.EthernetName))
	insertDB(_ctx, params, gateWayClient)
	runConsul(ctx)
	gateway.Join()
	return nil
}

func insertES(_ctx context.Context,
	cli *elastic.Client,
	params *common.APICONTROLLERPARAMS,
	_type,
	_result string) {
	_elasticsearch.InsertESDB2(
		_ctx,
		cli,
		&common.ElasticSearchParam2{
			ID:         params.ID,
			Name:       params.Name,
			DC:         params.Dc,
			Version:    params.Version,
			CreateTime: time.Now(),
			Type:       _type,
			Result:     _result,
		})
}

//generateID --
func generateID() string {
	hd := hashids.NewData()
	hd.Salt = fmt.Sprintf("%d", time.Now().UnixNano())
	hd.MinLength = 6
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{12, 34, 56, 78, 90})
	return e
}

func runConsul(ctx *cli.Context) {
	go func() {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		dir += "/consul"
		log.WithFields(log.Fields{
			"consul_dir":    dir,
			"datacenter":    args.DataCenter,
			"consul_leader": args.ConsulLeaderName,
		}).Info("consul app dir")
		bin := dir
		a1 := "agent"
		a2 := "-data-dir=/opt/consul"
		a3 := "-node=api-gateway-" + generateID()
		a4 := "-bind=" + common.GetIPFromNetWork("eth0")
		a5 := "-enable-script-checks=true"
		a6 := "-config-dir=/etc/consul.d"
		a7 := "-datacenter=" + args.DataCenter
		a8 := "-retry-join=" + args.ConsulLeaderName
		c := exec.Command(bin, a1, a2, a3, a4, a5, a6, a7, a8)
		// c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	}()
}

func getEnv() {
	if os.Getenv("ANYUN_NETWORKID") != "" {
		args.NetWorkID = os.Getenv("ANYUN_NETWORID")
	}
	if os.Getenv("ANYUN_DNS") != "" {
		args.DNS = os.Getenv("ANYUN_DNS")
	}
	if os.Getenv("ANYUN_DC") != "" {
		args.DataCenter = os.Getenv("ANYUN_DC")
	}
	if os.Getenv("ANYUN_CONSUL_LEADER") != "" {
		args.ConsulLeaderName = os.Getenv("ANYUN_CONSUL_LEADER")
	}
	if os.Getenv("ANYUN_IFNAME") != "" {
		args.EthernetName = os.Getenv("ANYUN_IFNAME")
	}
	if os.Getenv("ANYUN_NATS_MGR") != "" {
		args.NatsManager = os.Getenv("ANYUN_NATS_MGR")
	}
	if os.Getenv("ANYUN_NATS_BUSINESS") != "" {
		args.NatsBusiness = os.Getenv("ANYUN_NATS_BUSINESS")
	}
	if os.Getenv("ANYUN_ETCD") != "" {
		args.Etcd = os.Getenv("ANYUN_ETCD")
	}
	if os.Getenv("ANYUN_REDIS") != "" {
		args.Redis = os.Getenv("ANYUN_REDIS")
	}
	if os.Getenv("ANYUN_DNS_SRV_API") != "" {
		args.ConsulAddress = os.Getenv("ANYUN_DNS_SRV_API")
	}
	if os.Getenv("ANYUN_ES_URL") != "" {
		args.ESURL = os.Getenv("ANYUN_ES_URL")
	}
	if os.Getenv("ANYUN_ES_USERNAME") != "" {
		args.ESUserName = os.Getenv("ANYUN_ES_USERNAME")
	}
	if os.Getenv("ANYUN_ES_PASSWORD") != "" {
		args.ESPassWord = os.Getenv("ANYUN_ES_PASSWORD")
	}
}

func insertDB(
	_ctx context.Context,
	params *common.APICONTROLLERPARAMS,
	gateWayClient *connection.GatewayClientContext) {
	go func() {
		for {
			time.Sleep(300 * 1e9)
			log.Info("connect to container-manager to check out that is this api the last one.")
			if common.APICount == nil || *common.APICount == 0 {
				r, _ := json.Marshal(&common.APICONTROLLERPARAMS{
					TypeDO:        "1",
					ID:            params.ID,
					Name:          params.Name,
					Version:       params.Version,
					Dc:            params.Dc,
					ConsulAddress: params.ConsulAddress,
					DiscoveryDNS:  params.DiscoveryDNS,
				})
				mess, err := gateWayClient.NatsBusiness.Request(common.APIContainerManager, []byte(r), 10000*time.Millisecond)
				if err != nil {
					log.Error(err.Error())
				}
				if mess != nil && mess.Data != nil {
					log.Info("try to delete this container,then the server return message is ", string(mess.Data))
					if string(mess.Data) == "shutdown" {
						log.Info("try to delete this container")
						gateWayClient.Etcd.Delete(context.TODO(), "/container/api/"+server.ETCDTime.TIME)
						insertES(_ctx,
							gateWayClient.ESClient,
							params,
							"delete",
							"success",
						)
						log.Info("see you")
						os.Exit(0)
					} else {
						log.Info("connect to ES to insert recode that this request failed")
						go func() {
							insertES(_ctx,
								gateWayClient.ESClient,
								params,
								"delete",
								"failed,this is the last",
							)
						}()
					}
				} else {
					log.Error("nats return error")
					go func() {
						insertES(_ctx,
							gateWayClient.ESClient,
							params,
							"delete",
							"failed,error",
						)
					}()
				}
				log.Info("this request will try again 5 mins later")
			}
		}
	}()
}
