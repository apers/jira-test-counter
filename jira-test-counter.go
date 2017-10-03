package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const ServerPort = "80"

var db JiraDb

const DevelopmentCol = "In Progress"
const CodeReviewCol = "Klar til code review"
const TestCol = "Testbar"
const DoneCol = "Done"

const TaskTypeTest = "test"
const TaskTypeReview = "review"

func main() {
	fmt.Println(time.Now())
	fmt.Println("Staring server..")

	db = dbConnect()
	db.cleanTables()
	db.initTables()

	http.HandleFunc("/webhook", webHookHandler)
	http.HandleFunc("/stats", statsHandler)
	log.Fatal(http.ListenAndServe(":"+ServerPort, nil))
}
