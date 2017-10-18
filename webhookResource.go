package main

import (
	"net/http"
	"fmt"
	"time"
)

const NotPassedCol = "Ikke passert"
const DevelopmentCol = "In Progress"
const CodeReviewCol = "Klar til code review"
const TestCol = "Testbar"
const DoneCol = "Done"

const TaskTypeDev = "development"
const TaskTypeReview = "review"
const TaskTypeTest = "test"

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	buf := readReader(r.Body)

	if buf.Len() == 0 {
		return
	}

	fmt.Println(buf)

	event := convertToJiraJson(buf)
	if event.ChangeLog.hasStatusChange() {
		from, to := event.ChangeLog.getStatusChange()
		var taskType string
		if from == CodeReviewCol && to == TestCol && event.Issue.getAssignee() != event.User.Name {
			fmt.Println("CodeReview")
			fmt.Println("PF: ", event.Issue.Key)
			fmt.Println("User: ", event.User.Name)

			taskType = TaskTypeReview
		} else if from == TestCol && to == DoneCol && event.Issue.getAssignee() != event.User.Name {
			fmt.Println("Test")
			fmt.Println("PF: ", event.Issue.Key)
			fmt.Println("User: ", event.User.Name)

			taskType = TaskTypeTest
		} else if from == DevelopmentCol && to == CodeReviewCol && event.Issue.getAssignee() == event.User.Name {
			fmt.Println("Development")
			fmt.Println("PF: ", event.Issue.Key)
			fmt.Println("User: ", event.User.Name)

			taskType = TaskTypeDev
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
