package discovery

import (
	"context"
	"fmt"
	"time"

	"testing"

	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func befor() {
	logLevel, _ := log.ParseLevel("debug")
	log.SetLevel(logLevel)
	formatter := new(prefixed.TextFormatter)
	log.SetFormatter(formatter)
}

func TestMessage(t *testing.T) {
	b := buildMessage()
	fmt.Println(string(b))
	client, _ := nats.Connect("nats://message-business.service.consul:4222")
	mess, err := client.Request("annyun-cloud-micro-service-0AAEB6F5B4112BCC79D6C270451E63E8", b, 10*time.Millisecond)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(mess.Data))
}

func TestEtcd(t *testing.T) {
	befor()
	if queryClient, err := NewQueryClient(context.Background(), []string{"server.etcd.service.dc-anyuncloud.consul:2379"}, "192.168.252.254", "web-manager.service.consul:6379"); err != nil {
		log.Fatalln(err.Error())
	} else {
		// // queryClient.redisClient.Set("name", "yw", 0)
		defer queryClient.RedisClient.Close()

		var s1 []string
		s1 = append(s1, "param1-value")
		var s2 []string
		s2 = append(s2, "param2-value1", "param2-value2")
		m := make(map[string][]string)
		m["key1"] = s1
		m["key2"] = s2

		request, err := queryClient.CheckQueryURL(&RequestDetail{
			url:         "/api-gateway/v1/index",
			method:      "GET",
			values:      m,
			contentType: "application/json",
			acceptType:  "application/json",
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		client, _ := nats.Connect("nats://message-business.service.consul:4222")
		mess, err := client.Request("C5CF1B523BBABB296C63963FF39F114D", request, 10*time.Millisecond)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(mess.Data))

		// v, _ := queryClient.redisClient.Get("test_key123").Result()
		// fmt.Println(v)

		// queryClient.etcdClient.Put(queryClient.ParentContext, "/api-gateway/index", "B2481F3B-4146-4251-AC01-58205CB589BC")
		// queryClient.etcdClient.Put(queryClient.ParentContext, "/meta/api_meta/B2481F3B-4146-4251-AC01-58205CB589BC/desc", "this is my api")
		// queryClient.etcdClient.Put(queryClient.ParentContext, "/meta/api_meta/B2481F3B-4146-4251-AC01-58205CB589BC/versions", "v1,v2,v3")
		// queryClient.etcdClient.Put(queryClient.ParentContext, "/meta/api_meta/B2481F3B-4146-4251-AC01-58205CB589BC/v1/method", "GET,PUT,POST,DELETE,OPTION,HEAD")

		// queryClient.etcdClient.Put(queryClient.ParentContext, "/meta/api_meta/B2481F3B-4146-4251-AC01-58205CB589BC/v1/request/GET/params/key1:0", "param1-value")
		// queryClient.etcdClient.Put(queryClient.ParentContext, "/meta/api_meta/B2481F3B-4146-4251-AC01-58205CB589BC/v1/request/GET/params/key2:0", "param2-value1")
		// queryClient.etcdClient.Put(queryClient.ParentContext, "/meta/api_meta/B2481F3B-4146-4251-AC01-58205CB589BC/v1/request/GET/params/key2:1", "param2-value2")

		// queryClient.etcdClient.Put(queryClient.ParentContext, "/meta/api_meta/B2481F3B-4146-4251-AC01-58205CB589BC/v1/request/POST/body", "application/json")
		// queryClient.etcdClient.Put(queryClient.ParentContext, "/meta/api_meta/B2481F3B-4146-4251-AC01-58205CB589BC/v1/request/POST/product", "application/json")
	}

}
