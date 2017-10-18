package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type JiraDb struct {
	db *sql.DB
}

type User struct {
	Username      string
	Email         string
	AvailableBlocks int
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

/* Tasks */

func (db JiraDb) addTask(username string, taskType string, key string) {
	stmt, err := db.db.Prepare("INSERT INTO tasks(username, type, key, time) VALUES($1, $2, $3, $4)")
	check(err)
	_, err = stmt.Exec(username, taskType, key, time.Now())
	check(err)
}

func (db JiraDb) addToAvailableBlocks(username string, taskType string) {

	blockCount := 0

	if taskType == TaskTypeReview {
		blockCount = 6
	} else if taskType == TaskTypeTest {
		blockCount = 4
	} else if taskType == TaskTypeDev {
		blockCount = 2
	}

	stmt, err := db.db.Prepare("UPDATE users SET available_blocks = available_blocks + $1 WHERE username = $2")
	check(err)
	_, err = stmt.Exec(blockCount, username)
	check(err)
}

func (db JiraDb) getAllTaskCount() *sql.Rows {
	stmt, err := db.db.Prepare("SELECT username, available_blocks FROM users")
	defer stmt.Close()
	check(err)
	res, err := stmt.Query()
	check(err)
	return res
}

/* User*/

func (db JiraDb) createUser(username string, email string) {
	fmt.Println("Creating user: ", username)
	stmt, err := db.db.Prepare("INSERT INTO users(username, email) VALUES($1, $2)")
	defer stmt.Close()
	check(err)
	_, err = stmt.Exec(username, email)
	check(err)
}

func (db JiraDb) getUser(username string) (error, User) {
	var user User
	err := db.db.QueryRow("SELECT username, email, available_blocks FROM users WHERE username=$1", username).Scan(&user.Username, &user.Email, &user.AvailableBlocks)

	if err != nil && err != sql.ErrNoRows {
		check(err)
	}

	return err, user
}

func (db JiraDb) updateAvailableBlocks(username string, availableBlocks int) {
	fmt.Println("Updating available blocks for user: ", username)
	stmt, err := db.db.Prepare("UPDATE users SET available_blocks = $1 WHERE username = $2")
	defer stmt.Close()
	check(err)
	_, err = stmt.Exec(availableBlocks, username)
	check(err)
}

func (db JiraDb) getUserStats(username string) (error, int) {
	var taskCount int
	err := db.db.QueryRow("SELECT count(*) FROM tasks WHERE username=$1", username).Scan(&taskCount)

	if err != nil && err != sql.ErrNoRows {
		check(err)
	}

	return err, taskCount
}

/* Tables */

func (db JiraDb) cleanTables() {
	fmt.Println("Cleaning tables..")

	_, err := db.db.Exec("DROP TABLE tasks")
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.db.Exec("DROP TABLE users")
	if err != nil {
		fmt.Println(err)
	}
}

func (db JiraDb) initTables() {
	fmt.Println("Initializing tables..")

	_, err := db.db.Exec(`
		create table users (
			username varchar(50) primary key,
			email varchar(100),
			available_blocks integer DEFAULT 0
		);`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = db.db.Exec(`
		create table tasks (
			  id serial primary key,
			  username varchar(50) references users(username),
			  type varchar(50),
			  key varchar(50),
			  time timestamp
		);`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = db.db.Exec(`
		create table blocks (
			  username varchar(50) references users(username),
			  material varchar(50),
			  x int,
			  y int,
			  z int,
			  PRIMARY KEY(x, y, z)
			  time timestamp
		);`)

	if err != nil {
		fmt.Println(err)
	}
}
