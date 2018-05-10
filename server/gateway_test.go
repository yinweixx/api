package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"time"

	"e.coding.net/anyun-cloud-api-gateway/connection"
	"e.coding.net/anyun-cloud-service-manager/newconn"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/gin-gonic/gin"
	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

//根据一些手动默认配置测试API网关的启动
func TestAPIServerStartup(t *testing.T) {
	engine := gin.New()
	cfg := new(AnyunCloudGatewayConfig)
	cfg.HTTPS.TLSConfig = DefaultTLSConfig()
	cfg.ListenerAddr = ":9000"
	gateway := &AnyunCloudGateway{
		Engine: engine,
		Config: cfg,
	}
	gateway.Start()
	gateway.Join()
}

func getTestAnyunCloudGateway() *AnyunCloudGateway {
	gin.SetMode(gin.ReleaseMode)
	log.SetLevel(log.DebugLevel)
	formatter := new(prefixed.TextFormatter)
	log.SetFormatter(formatter)
	engine := gin.New()
	cfg := new(AnyunCloudGatewayConfig)
	cfg.HTTPS.TLSConfig = DefaultTLSConfig()
	cfg.ListenerAddr = ":9000"
	gateway := &AnyunCloudGateway{
		Engine: engine,
		Config: cfg,
	}
	return gateway
}

// //API网关中间件测试
// func TestMiddleware(T *testing.T) {
// 	gateway := getTestAnyunCloudGateway()
// 	gateway.SetUpMiddlewares()
// 	gateway.Start()
// 	vt := gateway.Engine.Group("/test")
// 	{
// 		vt.GET("/api1", func(c *gin.Context) {
// 			c.String(http.StatusOK, "test1")
// 		})
// 	}
// 	gateway.Join()
// }

// //测试内置的API
// func TestBuildInAPI(t *testing.T) {
// 	gateway := getTestAnyunCloudGateway()
// 	gateway.SetUpMiddlewares()
// 	gateway.Start()
// 	gateway.Initialization()
// 	gateway.Join()
// }

func TestHTTPS(t *testing.T) {
	crt := "/Users/twitchgg/Develop/Projects/goproject/src/e.coding.net/anyun-cloud-api-gateway/ssl/server.crt"
	key := "/Users/twitchgg/Develop/Projects/goproject/src/e.coding.net/anyun-cloud-api-gateway/ssl/server.key"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi there!")
	})
	http.ListenAndServeTLS(":8081", crt, key, nil)
}

func TestEtcdKeepALIVE(t *testing.T) {
	gateWayClient, _ := connection.NewGatewayClientContext(context.Background(), &connection.GatewayClientConfig{
		NatsManager:  "",
		NatsBusiness: "",
		Etcd:         nil,
		Redis:        "",
	})
	cli := gateWayClient.Etcd
	resp, err := cli.Grant(context.TODO(), 3)
	if err != nil {
		log.Fatal(err)
	}
	_, err = cli.Put(context.TODO(), "/container/api/"+string(resp.ID), time.Now().UTC().Format(time.UnixDate), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
	ch, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}
	ka := <-ch
	fmt.Println("ttl:", ka.TTL)
	gresp, err := cli.Get(context.TODO(), "/container/api/", clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("number of keys:", gresp.Count)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func Test3(t *testing.T) {
	cli := newconn.NewEtcdClient()

	defer cli.Close()

	// create two separate sessions for lock competition
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	m1 := concurrency.NewMutex(s1, "/my-lock/")

	s2, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s2.Close()
	m2 := concurrency.NewMutex(s2, "/my-lock/")

	// acquire lock for s1
	if err := m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for s1")

	m2Locked := make(chan struct{})
	go func() {
		defer close(m2Locked)
		// wait until s1 is locks /my-lock/
		if err := m2.Lock(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(10 * 1e9)
	if err := m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("released lock for s1")

	<-m2Locked
	fmt.Println("acquired lock for s2")
}

func Test1(t *testing.T) {
	client, err := nats.Connect("nats://tcp.message-business.service.dc-anyuncloud.consul:4222")
	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("manager nats connected")
	}

	i := 0
	for i <= 3 {
		// time.Sleep(100 * time.Millisecond)
		go func() {
			// for {
			req, _ := client.Request("API_CONTROLLER_MANAGER", []byte("help me"), 10000*time.Millisecond)
			fmt.Println(string(req.Data))
			// }
		}()

		i++
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
