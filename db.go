package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"encoding/json"
	"log"
)

/* objects */

type JiraDb struct {
	db *sql.DB
}

type User struct {
	username        string
	email           string
	availableBlocks int
}

/* db connections */

func dbConnect() JiraDb {
	db, err := sql.Open("postgres", "user=jira dbname=jira password=jira")
	check(err)
	err = db.Ping()
	check(err)

	var JiraDb JiraDb
	JiraDb.db = db

	return JiraDb
}

/* tasks */

func (db JiraDb) addTask(username string, taskType string, key string, summary string) {
	stmt, err := db.db.Prepare("INSERT INTO tasks(username, type, key, summary, time) VALUES($1, $2, $3, $4, $5)")
	check(err)
	_, err = stmt.Exec(username, taskType, key, summary, time.Now())
	check(err)
}

func (db JiraDb) addToAvailableBlocks(username string, taskType string) {

	blockCount := 0

	if taskType == TaskTypeReview {
		blockCount = 6
	} else if taskType == TaskTypeTest {
		blockCount = 8
	} else if taskType == TaskTypeDev {
		blockCount = 4
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

func (db JiraDb) getMainTaskCount() *sql.Rows {
	stmt, err := db.db.Prepare(`
		SELECT 'MAIN', COUNT(*)
		FROM tasks
		WHERE username IN ('ajo', 'rbr', 'vd', 'kah', 'hvg');`)

	defer stmt.Close()
	check(err)
	res, err := stmt.Query()
	check(err)
	return res
}

func (db JiraDb) getCoreTaskCount() *sql.Rows {
	stmt, err := db.db.Prepare(`
		SELECT 'CORE', COUNT(*)
		FROM tasks
		WHERE username IN ('rmg', 'ap', 'wab', 'jbe', 'tkg');`)

	defer stmt.Close()
	check(err)
	res, err := stmt.Query()
	check(err)
	return res
}

/* user*/

func (db JiraDb) createUser(username string, email string) {
	fmt.Println("creating user: ", username)
	stmt, err := db.db.Prepare("INSERT INTO users(username, email) VALUES($1, $2)")
	defer stmt.Close()
	check(err)
	_, err = stmt.Exec(username, email)
	check(err)
}

func (db JiraDb) getUser(username string) (error, User) {
	var user User
	err := db.db.
		QueryRow("SELECT username, email, available_blocks FROM users WHERE username=$1", username).
		Scan(&user.username, &user.email, &user.availableBlocks)

	if err != nil && err != sql.ErrNoRows {
		check(err)
	}

	return err, user
}

func (db JiraDb) updateAvailableBlocks(username string, availableBlocks int) {
	fmt.Println("updating available blocks for user: ", username)
	stmt, err := db.db.Prepare("UPDATE users SET available_blocks = $1 WHERE username = $2")
	defer stmt.Close()
	check(err)
	_, err = stmt.Exec(availableBlocks, username)
	check(err)
}

func (db JiraDb) getUserStats(username string) (error, int) {
	var taskcount int
	err := db.db.QueryRow("SELECT count(*) FROM tasks WHERE username=$1", username).Scan(&taskcount)

	if err != nil && err != sql.ErrNoRows {
		check(err)
	}

	return err, taskcount
}

/* Tables and Migrations */

func (db JiraDb) cleanTables() {
	fmt.Println("cleaning tables..")

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
	log.Print("initializing tables..")

	_, err := db.db.Exec(`
		CREATE TABLE users (
			username VARCHAR(50) PRIMARY KEY,
			email VARCHAR(100),
			available_blocks INTEGER DEFAULT 0
		);`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = db.db.Exec(`
		CREATE TABLE tasks (
			  id SERIAL PRIMARY KEY,
			  username VARCHAR(50) REFERENCES users(username),
			  type VARCHAR(50),
			  key VARCHAR(50),
			  summary VARCHAR(200),
			  time TIMESTAMP
		);`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = db.db.Exec(`
		CREATE TABLE blocks (
			  username VARCHAR(50) REFERENCES users(username),
			  material VARCHAR(50),
			  x INT,
			  y INT,
			  z INT,
			  PRIMARY KEY(x, y, z),
			  time TIMESTAMP
		);`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = db.db.Exec(`
		CREATE TABLE sql_update (
			current_version INT PRIMARY KEY,
            		update_time TIMESTAMP DEFAULT now()
		);`)

	if err != nil {
		fmt.Println(err)
	}
}

type Migrations struct {
	MigrationEntries []MigrationEntry `json:"migrations"`
}

type MigrationEntry struct {
	source  string `json:"source"`
	version int    `json:"version"`
}

func (db JiraDb) migrate() {
	log.Print("migrating..")

	var migrations Migrations

	file := readFile("./sql-migrations/migrations.json")
	err := json.Unmarshal(file.Bytes(), &migrations)

	check(err)

	for _, migrationEntry := range migrations.MigrationEntries {
		fmt.Println(migrationEntry);
	}

}
