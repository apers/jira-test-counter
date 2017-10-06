package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const ServerPort = "80"

var db JiraDb

func main() {
	fmt.Println(time.Now())
	fmt.Println("Staring server..")

	db = dbConnect()
	//db.cleanTables()
	db.initTables()

	http.HandleFunc("/webhook", webhookHandler)
	http.HandleFunc("/stats", statsHandler)
	log.Fatal(http.ListenAndServe(":" + ServerPort, nil))
}
