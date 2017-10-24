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
	availableblocks int
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

func (db JiraDb) addTask(username string, tasktype string, key string, summary string) {
	stmt, err := db.db.Prepare("insert into tasks(username, type, key, summary, time) values($1, $2, $3, $4, $5)")
	check(err)
	_, err = stmt.Exec(username, tasktype, key, summary, time.Now())
	check(err)
}

func (db JiraDb) addToAvailableBlocks(username string, tasktype string) {

	blockcount := 0

	if tasktype == TaskTypeReview{
		blockcount = 6
	} else if tasktype == TaskTypeTest {
		blockcount = 8
	} else if tasktype == TaskTypeDev {
		blockcount = 4
	}

	stmt, err := db.db.Prepare("update users set available_blocks = available_blocks + $1 where username = $2")
	check(err)
	_, err = stmt.Exec(blockcount, username)
	check(err)
}

func (db JiraDb) getAllTaskCount() *sql.Rows {
	stmt, err := db.db.Prepare("select username, available_blocks from users")
	defer stmt.Close()
	check(err)
	res, err := stmt.Query()
	check(err)
	return res
}

/* user*/

func (db JiraDb) createUser(username string, email string) {
	fmt.Println("creating user: ", username)
	stmt, err := db.db.Prepare("insert into users(username, email) values($1, $2)")
	defer stmt.Close()
	check(err)
	_, err = stmt.Exec(username, email)
	check(err)
}

func (db JiraDb) getUser(username string) (error, User) {
	var user User
	err := db.db.QueryRow("select username, email, available_blocks from users where username=$1", username).Scan(&user.username, &user.email, &user.availableblocks)

	if err != nil && err != sql.ErrNoRows {
		check(err)
	}

	return err, user
}

func (db JiraDb) updateAvailableBlocks(username string, availableblocks int) {
	fmt.Println("updating available blocks for user: ", username)
	stmt, err := db.db.Prepare("update users set available_blocks = $1 where username = $2")
	defer stmt.Close()
	check(err)
	_, err = stmt.Exec(availableblocks, username)
	check(err)
}

func (db JiraDb) getUserStats(username string) (error, int) {
	var taskcount int
	err := db.db.QueryRow("select count(*) from tasks where username=$1", username).Scan(&taskcount)

	if err != nil && err != sql.ErrNoRows {
		check(err)
	}

	return err, taskcount
}

/* Tables and Migrations */

func (db JiraDb) cleanTables() {
	fmt.Println("cleaning tables..")

	_, err := db.db.Exec("drop table tasks")
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.db.Exec("drop table users")
	if err != nil {
		fmt.Println(err)
	}
}

func (db JiraDb) initTables() {
	log.Print("initializing tables..")

	_, err := db.db.Exec(`
		create table users (
			username varchar(50) primary key,
			email varchar(100),
			available_blocks integer default 0
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
			  primary key(x, y, z),
			  time timestamp
		);`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = db.db.Exec(`
		create table sql_update (
			current_version int primary key,
            		update_time timestamp default now()
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
	version int `json:"version"`
}

func (db JiraDb) migrate() {
	log.Print("migrating..")

	var migrations Migrations

	file := readFile("./sql-update/migrations.json")
	err := json.Unmarshal(file.Bytes, &migrations)

	check(err)

	for _, migrationEntry := range migrations.MigrationEntries {
		fmt.Println(migrationEntry);
	}

}
