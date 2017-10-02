package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func readFile(filename string) *bytes.Buffer {
	file, err := os.Open(filename)
	defer file.Close()
	check(err)
	return readReader(file)
}

func readReader(reader io.Reader) *bytes.Buffer {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	check(err)
	return buf
}

func convertToJson(rawJson *bytes.Buffer) JiraEvent {
	var event JiraEvent
	err := json.Unmarshal(rawJson.Bytes(), &event)
	check(err)
	return event
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	buf := readReader(r.Body)
	if buf.Len() == 0 {
		return
	}
	event := convertToJson(buf)
	if event.ChangeLog.hasStatusChange() && event.Issue.isFlagged() {
		from, to := event.ChangeLog.getStausChange()
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

	} else {
		fmt.Println("Non-status event..", time.Now())
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	var username string
	var count int
	rows := db.getAllTaskCount()
	for rows.Next() {
		rows.Scan(username, count)
		fmt.Fprint(w, "%s - %d", username, count)
	}
}

func main() {
	fmt.Println(time.Now())
	fmt.Println("Staring server..")

	db = dbConnect()
	db.cleanTables()
	db.initTables()

	http.HandleFunc("/webhook", webhookHandler)
	http.HandleFunc("/stats", statsHandler)
	log.Fatal(http.ListenAndServe(":" + ServerPort, nil))
}
