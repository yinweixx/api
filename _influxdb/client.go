package _influxdb

import (
	"time"

	log "github.com/sirupsen/logrus"

	client "github.com/influxdata/influxdb/client/v2"
)

const (
	MyDB     = "golang_yw_test"
	username = " "
	password = " "
)

func CreateClient() client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://192.168.254.239:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	return c
}

//InsertDB --
func (str *InfluxDBStructs) InsertDB() {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: str.DataBase,
		// Precision: "s",
		RetentionPolicy: str.RetentionPolicy,
	})
	if err != nil {
		log.Error(err)
	}

	pt, err := client.NewPoint(str.Table, str.Tags, str.Fields, time.Now())
	if err != nil {
		log.Error(err)
	}
	bp.AddPoint(pt)
	// Write the batch
	cli := *str.Client
	if err := cli.Write(bp); err != nil {
		log.Error(err)
	}
}
