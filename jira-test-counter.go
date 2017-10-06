package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const ServerPort = "80"

var db JiraDb

func readFile(filename string) *bytes.Buffer {
	file, err := os.Open(filename)
	defer file.Close()
	check(err)
	return readReader(file)
}

func convertToJiraJson(rawJson *bytes.Buffer) JiraEvent {
	var event JiraEvent
	err := json.Unmarshal(rawJson.Bytes(), &event)
	check(err)
	return event
}

func convertToUpdateBlockCountJson(rawJson *bytes.Buffer) MinecraftEvent {
	var event MinecraftEvent
	err := json.Unmarshal(rawJson.Bytes(), &event)
	check(err)
	return event
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	buf := readReader(r.Body)
	if buf.Len() == 0 {
		return
	}
	event := convertToJiraJson(buf)
	if event.ChangeLog.hasStatusChange() && event.Issue.isFlagged() {
		from, to := event.ChangeLog.getStatusChange()
		var taskType string
		if from == CodeReviewCol && to == TestCol {
			fmt.Println("CodeReview")
			fmt.Println("PF: ", event.Issue.Key)
			fmt.Println("Flagged: ", event.Issue.isFlagged())
			fmt.Println("User: ", event.User.Name)

			taskType = TaskTypeReview
		} else if from == TestCol && to == DoneCol {
			fmt.Println("Test")
			fmt.Println("PF: ", event.Issue.Key)
			fmt.Println("Flagged: ", event.Issue.isFlagged())
			fmt.Println("User: ", event.User.Name)

			taskType = TaskTypeTest
		} else {
			// Unspported task
			return
		}

		err, _ := db.getUser(event.User.Name)

		// No such user
		if err != nil {
			db.createUser(event.User.Name, event.User.Email)
		}

		db.addTask(event.User.Name, taskType, event.Issue.Key)
		db.addToAvailableBlocks(event.User.Name, taskType)
	} else {
		fmt.Println("Non-status event..", time.Now())
	}
}

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
