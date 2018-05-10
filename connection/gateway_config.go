package connection

//GetConnInfo -- connection config
func GetConnInfo() *GatewayClientConfig {
	return &GatewayClientConfig{
		Etcd:                    []string{"192.168.252.15:2379"},
		NatsBusiness:            "nats://tcp.message-business.service.dc-anyuncloud.consul:4222",
		NatsManager:             "nats://tcp.message-business.service.dc-anyuncloud.consul:4222",
		Redis:                   "redis.service.dc-anyuncloud.consul:6379",
		InfluxDB:                "MontiorInformation",
		InfluxDBURL:             "http://192.168.254.239:8086",
		InfluxDBUser:            " ",
		InfluxDBPass:            " ",
		InfluxDBRetentionPolicy: "aRetentionPolicy",
		ESURL:      "http://192.168.252.15:9200",
		ESUserName: "elastic",
		ESPassword: "p278DKNSSDtlPZGaobD7",
	}
}

// MysqlClient:  "root:1234qwer@tcp(client.mysql.service.consul:3306)/test?charset=utf8", //http://pear.php.net/manual/en/package.database.db.intro-dsn.php
