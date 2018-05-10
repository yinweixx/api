package _influxdb

import client "github.com/influxdata/influxdb/client/v2"

type InfluxDBStructs struct {
	DataBase        string
	RetentionPolicy string
	Tags            map[string]string
	Fields          map[string]interface{}
	Table           string
	Client          *client.Client
}
