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

type jiradb struct {
	db *sql.db
}

type user struct {
	username        string
	email           string
	availableblocks int
}

/* db connections */

func dbconnect() jiradb {
	db, err := sql.open("postgres", "user=jira dbname=jira password=jira")
	check(err)
	err = db.ping()
	check(err)

	var jiradb jiradb
	jiradb.db = db

	return jiradb
}

/* tasks */

func (db jiradb) addtask(username string, tasktype string, key string, summary string) {
	stmt, err := db.db.prepare("insert into tasks(username, type, key, summary, time) values($1, $2, $3, $4, $5)")
	check(err)
	_, err = stmt.exec(username, tasktype, key, summary, time.now())
	check(err)
}

func (db jiradb) addtoavailableblocks(username string, tasktype string) {

	blockcount := 0

	if tasktype == tasktypereview {
		blockcount = 6
	} else if tasktype == tasktypetest {
		blockcount = 4
	} else if tasktype == tasktypedev {
		blockcount = 2
	}

	stmt, err := db.db.prepare("update users set available_blocks = available_blocks + $1 where username = $2")
	check(err)
	_, err = stmt.exec(blockcount, username)
	check(err)
}

func (db jiradb) getalltaskcount() *sql.rows {
	stmt, err := db.db.prepare("select username, available_blocks from users")
	defer stmt.close()
	check(err)
	res, err := stmt.query()
	check(err)
	return res
}

/* user*/

func (db jiradb) createuser(username string, email string) {
	fmt.println("creating user: ", username)
	stmt, err := db.db.prepare("insert into users(username, email) values($1, $2)")
	defer stmt.close()
	check(err)
	_, err = stmt.exec(username, email)
	check(err)
}

func (db jiradb) getuser(username string) (error, user) {
	var user user
	err := db.db.queryrow("select username, email, available_blocks from users where username=$1", username).scan(&user.username, &user.email, &user.availableblocks)

	if err != nil && err != sql.errnorows {
		check(err)
	}

	return err, user
}

func (db jiradb) updateavailableblocks(username string, availableblocks int) {
	fmt.println("updating available blocks for user: ", username)
	stmt, err := db.db.prepare("update users set available_blocks = $1 where username = $2")
	defer stmt.close()
	check(err)
	_, err = stmt.exec(availableblocks, username)
	check(err)
}

func (db jiradb) getuserstats(username string) (error, int) {
	var taskcount int
	err := db.db.queryrow("select count(*) from tasks where username=$1", username).scan(&taskcount)

	if err != nil && err != sql.errnorows {
		check(err)
	}

	return err, taskcount
}

/* tables and migrations */

func (db jiradb) cleantables() {
	fmt.println("cleaning tables..")

	_, err := db.db.exec("drop table tasks")
	if err != nil {
		fmt.println(err)
	}
	_, err = db.db.exec("drop table users")
	if err != nil {
		fmt.println(err)
	}
}

func (db jiradb) inittables() {
	log.print("initializing tables..")

	_, err := db.db.exec(`
		create table users (
			username varchar(50) primary key,
			email varchar(100),
			available_blocks integer default 0
		);`)

	if err != nil {
		fmt.println(err)
	}

	_, err = db.db.exec(`
		create table tasks (
			  id serial primary key,
			  username varchar(50) references users(username),
			  type varchar(50),
			  key varchar(50),
			  time timestamp
		);`)

	if err != nil {
		fmt.println(err)
	}

	_, err = db.db.exec(`
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
		fmt.println(err)
	}

	_, err = db.db.exec(`
		create table sql_update (
			current_version int primary key,
            		update_time timestamp default now()
		);`)

	if err != nil {
		fmt.println(err)
	}
}

type userstatscollection struct {
	users []*userstats
}

type userstats struct {
	username string
	tasks    int
}

type migrations struct {
	migrationentries []migrationentry `json:"migrations"`
}

type migrationentry struct {
	source  string `json:"source"`
	version int `json:"version"`
}


/*
{
  "migrations": [
    {
      "source": "1.sql",
      "version": "1"
    }
  ]
}
*/
func (db jiradb) migrate() {
	log.print("migrating..")

	var migrations migrations

	file := readfile("./sql-update/migrations.json")
	err := json.unmarshal(file.bytes(), &migrations)

	check(err)

	for _, migrationentry := range migrations.migrationentries {
		fmt.println(migrationentry);
	}

}
