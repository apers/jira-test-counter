package main

import (
	"net/http"
	"encoding/json"
	"bytes"
)

type UserStatsCollection struct {
	Users []*UserStats
}

type UserStats struct {
	Username string
	Tasks    int
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

func statsHandler(w http.ResponseWriter, r *http.Request) {
	// Post request coming from mineraft server
	// Todo - move to own handler and endpoint
	if r.Method == "POST" {
		buf := readReader(r.Body)
		if buf.Len() == 0 {
			return
		}
		mineCraftEvent := convertToUpdateBlockCountJson(buf)
		db.updateAvailableBlocks(mineCraftEvent.Username, mineCraftEvent.AvailableBlocks)
	}

	var statsColl UserStatsCollection
	var userStats *UserStats

	rows := db.getAllTaskCount()
	for rows.Next() {
		userStats = &UserStats{}
		rows.Scan(&userStats.Username, &userStats.Tasks)
		statsColl.Users = append(statsColl.Users, userStats)
	}

	js, err := json.Marshal(statsColl)
	check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}
