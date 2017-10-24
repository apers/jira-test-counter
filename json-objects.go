package main

type MinecraftEvent struct {
	Username string `json:"username"`
	AvailableBlocks int `json:"availableBlocks"`
}

type JiraEvent struct {
	User      JsonUser  `json:"user" `
	Issue     Issue     `json:"issue"`
	ChangeLog ChangeLog `json:"changelog"`
}

type JsonUser struct {
	Name  string `json:"name"`
	Email string `json:"emailAddress"`
}

type ChangeLog struct {
	Entries []ChangeEntry `json:"items"`
}

type ChangeEntry struct {
	Field string `json:"field"`
	From  string `json:"fromString"`
	To    string `json:"toString"`
}

type Issue struct {
	Fields Fields `json:"fields"`
	Key    string `json:"key"`
}

type Fields struct {
	Flagged []CustomField `json:"customfield_10602"`
	Assignee Assignee `json:"assignee"`
	Summary string `json:"summary"`
}

type Assignee struct {
	Name string `json:"name"`
}

type CustomField struct {
	Value string `json:"value"`
}

func (issue Issue) getAssignee() string {
	return issue.Fields.Assignee.Name
}

func (issue Issue) isFlagged() bool {
	if len(issue.Fields.Flagged) > 0 {
		return issue.Fields.Flagged[0].Value == "Impediment"
	} else {
		return false
	}
}

func (changelog ChangeLog) hasStatusChange() bool {
	for _, v := range changelog.Entries {
		if v.Field == "status" {
			return true
		}
	}
	return false
}

func (changelog ChangeLog) getStatusChange() (string, string) {
	for _, v := range changelog.Entries {
		if v.Field == "status" {
			return v.From, v.To
		}
	}
	return "", ""
}
