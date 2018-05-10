package _influxdb

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

func TestInfluxdbCreateDB(t *testing.T) {
	_, err := queryDB(CreateClient(), fmt.Sprintf("CREATE DATABASE %s", MyDB))
	if err != nil {
		log.Fatal(err)
	}
}

func TestInfluxdbSelectDB(t *testing.T) {
	res, err := queryDB(CreateClient(), fmt.Sprintf("select * from %s WHERE time > now() - 12h", "MontiorInformation.aRetentionPolicy.api_gateway"))
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range res {
		for _, n := range r.Series {
			fmt.Println(n.Values)
		}
	}
}

func TestCreateBatch(t *testing.T) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "MontiorInformation",
		// Precision: "s",
		RetentionPolicy: "aRetentionPolicy",
	})
	if err != nil {
		log.Fatal(err)
	}
	tags := map[string]string{
		"type": "api",
		"url":  "/camera-app-rtp/v1/dc",
	}
	fields := map[string]interface{}{
		"exec_result": "SUCCESS",
		"exec_time":   time.Now(),
		"exec_tx":     "123",
		"exec_rx":     "345",
	}

	pt, err := client.NewPoint("api_gateway", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := CreateClient().Write(bp); err != nil {
		log.Fatal(err)
	}
}

func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func TestTime(t *testing.T) {
	// m1 := time.Now()
	// time.Sleep(5 * 1e9)
	// m2 := time.Now()
	// fmt.Println(m2.Sub(m1))
	conn, err := net.Dial("tcp", "sina.com.cn:80")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	// var buf bytes.Buffer
	// io.Copy(&buf, conn)
	// fmt.Println("total size:", buf.Len())
	c, _ := ioutil.ReadAll(conn)
	fmt.Println(c)
}
