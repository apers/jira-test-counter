package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type UserStats struct {
	Username string
	Tasks    int
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	var userStatsArr []*UserStats
	var userStats *UserStats

	rows := db.getAllTaskCount()
	for rows.Next() {
		userStats = &UserStats{}
		rows.Scan(userStats.Username, userStats.Tasks)
		userStatsArr = append(userStatsArr, userStats)
	}

	js, err := json.Marshal(userStatsArr)
	check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
