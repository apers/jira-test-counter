package main

import (
	"log"
	"io"
	"bytes"
	"encoding/json"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
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
