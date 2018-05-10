package main

import (
	"e.coding.net/anyun-cloud-api-gateway/app"
	"e.coding.net/anyun-cloud-api-gateway/config"
)

func main() {
	app.RunApplication(config.GatewayInitFlags, config.Bootstrap)
}
