package main

import (
	"log"
	"miaoshaSystem/global"
	"miaoshaSystem/sql"
	"miaoshaSystem/web"
)

func main() {
	sql.ConnectMysql()
	sql.Init()

	if sql.DB == nil {
		log.Fatalf("MySQL database connection is not initialized")
	}

	global.CreateTable()

	go global.StartKafkaConsumer()

	web.Gin()
}
