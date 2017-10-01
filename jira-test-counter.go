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

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type JiraEvent struct {
	User  User  `json:"user" `
	Issue Issue `json:"issue"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"emailAddress"`
}

type Issue struct {
	Fields Fields `json:"fields"`
}

type Fields struct {
	Flagged []CustomField `json:"customfield_10200"`
}

type CustomField struct {
	Value string `json:"value"`
}

func (issue Issue) isFlagged() bool {
	return issue.Fields.Flagged[0].Value == "Impediment"
}

func readFile(filename string) *bytes.Buffer {
	file, err := os.Open(filename)
	defer file.Close()
	check(err)
	return readReader(file)
}

func readReader(reader io.Reader) *bytes.Buffer {
	buf := new(bytes.Buffer)
	count, err := buf.ReadFrom(reader)
	if count == 0 {
	}
	check(err)
	return buf
}

func convertToJson(rawJson *bytes.Buffer) JiraEvent {
	var event JiraEvent
	fmt.Println(rawJson)
	err := json.Unmarshal(rawJson.Bytes(), &event)
	fmt.Println(event.User.Name)
	fmt.Println("Issue is flagged:", event.Issue.isFlagged())
	check(err)
	return event
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	buf := readReader(r.Body)
	convertToJson(buf)
}

func main() {
	fmt.Println(time.Now())
	fmt.Println("Staring server..")
	//buf := readFile("body.json")
	//convertToJson(buf)
	http.HandleFunc("/webhook", requestHandler)
	log.Fatal(http.ListenAndServe(":"+ServerPort, nil))
}
