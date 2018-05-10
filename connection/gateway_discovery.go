package connection

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"e.coding.net/anyun-cloud-api-gateway/common"
	"e.coding.net/anyun-cloud-api-gateway/pool"
	"github.com/coreos/etcd/clientv3"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

//CheckQueryURL 检查请求地址
func (_this *GatewayClientContext) CheckQueryURL(request *common.RequestDetail) (string, error) {
	r, _ := json.Marshal(&request)
	md5param := md5Param(r)
	val, err := query(_this, md5param)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("illegar params")
	}
	if val != "" {
		return val, nil
	}
	mess, err := _this.NatsBusiness.Request(common.CheckQueryURL, []byte(r), 1000*time.Millisecond)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("illegar params")
		return "", err
	}
	pool.SetTTL(*_this.Redis, md5param, string(mess.Data), common.REDISTIMETOLIVE)
	return string(mess.Data), nil
}

//CheckFromEtcd -- check etcd
func (_this *GatewayClientContext) CheckFromEtcd(mess string) (string, error) {
	val, err := query(_this, mess)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("illegar params")
	}
	if val != "" {
		return val, nil
	}
	resp, _err := keyExist(_this, mess)
	log.Info("etcd return " + resp)
	if _err != nil {
		return "", err
	}

	if resp == "" {
		return "", nil
	}

	var etcdResp common.ETCDResp
	err = json.Unmarshal([]byte(resp), &etcdResp)

	if err != nil {
		return "", err
	}
	log.Info("service check return " + strings.Replace(etcdResp.Info, "subject://", "", -1))
	pool.SetTTL(*_this.Redis, mess, strings.Replace(etcdResp.Info, "subject://", "", -1), common.REDISTIMETOLIVE)

	return strings.Replace(etcdResp.Info, "subject://", "", -1), nil
}

func md5Param(param []byte) string {
	hasher := md5.New()
	hasher.Write(param)
	return hex.EncodeToString(hasher.Sum(nil))
}

func query(_this *GatewayClientContext, str string) (string, error) {
	conn := *_this.Redis
	value, err := redis.String(conn.Do("GET", str))
	if err != nil {
		log.Error("redis client working error")
		return "", err
	}
	log.WithFields(log.Fields{
		"prefix": "discovery.doCheck",
		"value":  value,
		"err":    err,
	}).Info("redisClient working")
	if value != "" {
		return value, nil
	}
	return "", nil
}

func queryEtcd(_this *GatewayClientContext, mess string) (string, error) {
	val, err := _this.Etcd.Get(_this.Context, mess)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("etcdClient error")
		return "", err
	}
	if len(val.Kvs) == 0 {
		return "", nil
	}
	return string(val.Kvs[0].Value), nil
}

func keyExist(_this *GatewayClientContext, mess string) (string, error) {
	resp, err := _this.Etcd.Get(_this.Context, mess, clientv3.WithPrefix())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("etcdClient error")
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}
