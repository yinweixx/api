package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/coreos/etcd/clientv3"
)

//Config --
type Config struct {
	EndPoints []string
	Transport CancelableTransport
}

type ETCDTIME struct {
	TIME string
}

var ETCDTime *ETCDTIME

//DefaultTransport --
var DefaultTransport CancelableTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
}

//CancelableTransport --
type CancelableTransport interface {
	http.RoundTripper
	CancelRequest(req *http.Request)
}

func (cfg *Config) getTransport() CancelableTransport {
	if cfg.Transport == nil {
		return DefaultTransport
	}
	return cfg.Transport
}

//EtcdKeepAlive --
func EtcdKeepAlive(cli *clientv3.Client) {
	resp, err := cli.Grant(context.TODO(), 3)
	if err != nil {
		log.Fatal(err)
	}
	_TIME := time.Now().UTC().Format(time.UnixDate)
	ETCDTime = &ETCDTIME{
		TIME: _TIME,
	}
	_, err = cli.Put(context.TODO(), "/container/api/"+_TIME, _TIME, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
	ch, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}
	ka := <-ch
	fmt.Println("ttl:", ka.TTL)
}
