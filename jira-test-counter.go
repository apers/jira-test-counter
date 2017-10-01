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

const DevelopmentCol = "In Progress"
const CodeReviewCol = "Code Review"
const TestCol = "Test"
const DoneCol = "Done"

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

func requestHandler(w http.ResponseWriter, r *http.Request) {
	buf := readReader(r.Body)

	event := convertToJson(buf)

	if event.ChangeLog.hasStatusChange() && event.Issue.isFlagged() {
		from, to := event.ChangeLog.getStausChange()
		fmt.Println(from, to)
		if from == CodeReviewCol && to == TestCol {
			fmt.Println("CodeReview")
			fmt.Println("PF: ", event.Issue.Key)
			fmt.Println("Flagged: ", event.Issue.isFlagged())
			fmt.Println("User: ", event.User.Name)
		}

		if from == TestCol && to == DoneCol {
			fmt.Println("Test")
			fmt.Println("PF: ", event.Issue.Key)
			fmt.Println("Flagged: ", event.Issue.isFlagged())
			fmt.Println("User: ", event.User.Name)
		}

	} else {
		fmt.Println("Non-status event..", time.Now())
	}
}

func main() {
	fmt.Println(time.Now())
	fmt.Println("Staring server..")
	db := dbConnect()
	db.cleanTables()
	//buf := readFile("body.json")
	//convertToJson(buf)
	//http.HandleFunc("/webhook", requestHandler)
	//log.Fatal(http.ListenAndServe(":"+ServerPort, nil))
}
