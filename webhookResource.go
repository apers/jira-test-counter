package main

import (
	"net/http"
	"fmt"
	"time"
)

func webHookHandler(w http.ResponseWriter, r *http.Request) {
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
