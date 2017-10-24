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
	db.initTables()

	http.HandleFunc("/webhook", webHookHandler)
	http.HandleFunc("/stats", statsHandler)

	log.Fatal(http.ListenAndServe(":" + ServerPort, nil))
}
