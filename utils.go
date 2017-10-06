package main

import (
	"log"
	"io"
	"bytes"
	"encoding/json"
	"os"
)

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
