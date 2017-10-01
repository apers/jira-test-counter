package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type JiraDb struct {
	db *sql.DB
}

func dbConnect() JiraDb {
	db, err := sql.Open("postgres", "user=jira dbname=jira password=jira")
	check(err)
	err = db.Ping()
	check(err)

	var jiraDb JiraDb
	jiraDb.db = db

	return jiraDb
}

func (db JiraDb) exec(sql string) sql.Result {
	res, err := db.db.Exec(sql)
	check(err)
	return res
}

func (db JiraDb) cleanTables() {
	fmt.Println("Cleaning tables..")

	_, err := db.db.Exec("DROP TABLE Users")
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.db.Exec("DROP TABLE Tasks")
	if err != nil {
		fmt.Println(err)
	}
}

func (db JiraDb) initTables() {
	fmt.Println("Initializing tables..")

	_, err := db.db.Exec(
		`create table users (
			id serial primary key,
			username varchar(50),
			email varchar(100),
			blocks_placed integer DEFAULT 0
		);`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = db.db.Exec("DROP TABLE Tasks")
	if err != nil {
		fmt.Println(err)
	}
}
