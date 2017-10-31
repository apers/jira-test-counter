package main

import (
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
)

type UserStatsCollection struct {
	users []*UserStats
	teams []*TeamStats
}

type UserStats struct {
	username         string
	available_blocks int
}

type TeamStats struct {
	teamname         string
	available_blocks int
}

type AchievementStat struct {
	most_reviews   string
	most_tests     string
	most_developed string
	ping_pong      string
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
	// Todo - move to own handler and endpoint
	// Post request coming from mineraft server
	if r.Method == "POST" {
		buf := readReader(r.Body)
		if buf.Len() == 0 {
			return
		}
		mineCraftEvent := convertToUpdateBlockCountJson(buf)
		db.updateAvailableBlocks(mineCraftEvent.Username, mineCraftEvent.AvailableBlocks)
	}
	// Todo - move to own handler and endpoint


	var statsColl UserStatsCollection
	var userStats *UserStats
	var teamStats *TeamStats;

	// Read user stats
	rows := db.getAllTaskCount()
	for rows.Next() {
		userStats = &UserStats{}
		rows.Scan(&userStats.username, &userStats.available_blocks)
		statsColl.users = append(statsColl.users, userStats)
	}

	fmt.Println("Tasks")
	fmt.Println(statsColl)

	// Read main team stats
	rows = db.getMainTaskCount()
	teamStats = &TeamStats{}
	rows.Scan(&teamStats.teamname, &teamStats.available_blocks)
	statsColl.teams = append(statsColl.teams, teamStats)

	fmt.Println("Main")
	fmt.Println(statsColl)


	// Read core stats
	rows = db.getCoreTaskCount()
	teamStats = &TeamStats{}
	rows.Scan(&teamStats.teamname, &teamStats.available_blocks)
	statsColl.teams = append(statsColl.teams, teamStats)

	fmt.Println("Core")
	fmt.Println(statsColl)

	js, err := json.Marshal(statsColl)
	check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func achievementStatsHandler(w http.ResponseWriter, r *http.Request) {
	var statsColl UserStatsCollection
	var userStats *UserStats

	rows := db.getAllTaskCount()
	for rows.Next() {
		userStats = &UserStats{}
		rows.Scan(&userStats.username, &userStats.available_blocks)
		statsColl.users = append(statsColl.users, userStats)
	}

	js, err := json.Marshal(statsColl)
	check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}
