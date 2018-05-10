package _elasticsearch

import (
	"context"
	"fmt"

	"e.coding.net/anyun-cloud-api-gateway/common"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

var (
	//URL --
	URL = "http://192.168.254.221:9200"
	//USERNAME --
	USERNAME = "elastic"
	//PASSWORD --
	PASSWORD = "p278DKNSSDtlPZGaobD7"
)

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"esp":{
			"properties":{
				"username":{
					"type":"keyword"
				},
				"apiname":{
					"type":"text"
				},
				"starttime":{
					"type":"date"
				},
				"result":{
					"type":"text"
				}
			}
		}
	}
}`

const mapping2 = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"esp":{
			"properties":{
				"id":{
					"type":"text"
				},
				"name":{
					"type":"text"
				},
				"dc":{
					"type":"text"
				},
				"version":{
					"type":"text"
				},
				"time":{
					"type":"date"
				},
				"containeripaddress":{
					"type":"text"
				},
				"networkid":{
					"type":"text"
				},
				"type":{
					"type":"text"
				},
				"result":{
					"type":"text"
				}
			}
		}
	}
}`

//NewClient -- new elasticsearch client
func NewClient(url, username, pass string) *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetBasicAuth(username, pass))
	if err != nil {
		// Handle error
		log.Error(err.Error())
	}
	return client
}

//InsertESDB --
func InsertESDB(ctx context.Context, client *elastic.Client, param *common.ElasticSearchParam) {
	exists, err := client.IndexExists("esearch").Do(ctx)
	if err != nil {
		// Handle error
		log.Error(err.Error())
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("esearch").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			log.Error(err.Error())
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	put1, err := client.Index().
		Index("esearch").
		Type("esp").
		BodyJson(param).
		Do(ctx)
	if err != nil {
		// Handle error
		log.Error(err.Error())
	}
	fmt.Printf("Indexed %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}

//InsertESDB2 --
func InsertESDB2(ctx context.Context, client *elastic.Client, param *common.ElasticSearchParam2) {
	exists, err := client.IndexExists("createcontainer").Do(ctx)
	if err != nil {
		log.Error(err.Error())
	}
	if !exists {
		createIndex, err := client.CreateIndex("createcontainer").BodyString(mapping2).Do(ctx)
		if err != nil {
			log.Error(err.Error())
		}
		if !createIndex.Acknowledged {
		}
	}
	put1, err := client.Index().
		Index("createcontainer").
		Type("esp").
		BodyJson(param).
		Do(ctx)
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("Indexed %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
