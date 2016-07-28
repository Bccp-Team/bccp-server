package mysql

import (
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Runner struct {
	Id              int
	Status          string
	Last_connection time.Time
	Ip              string
}

func (db *Database) ListRunners() []Runner {

	// Execute the query
	rows, err := db.conn.Query("SELECT * FROM runner")
	if err != nil {
		log.Fatal("ERROR: Unable to select runner: ", err.Error())
	}

	var runners []Runner

	// Fetch rows
	for rows.Next() {
		var id int
		var status string
		var last_connection time.Time
		var ip string
		// get RawBytes from data
		err = rows.Scan(&id, &status, &last_connection, &ip)
		if err != nil {
			log.Fatal("ERROR: Unable to get next row: ", err.Error())
		}

		runners = append(runners, Runner{id, status, last_connection, ip})
	}
	if err = rows.Err(); err != nil {
		log.Fatal("ERROR: Undefined row err: ", err.Error())
	}

	return runners
}

// Get runner info by id
// Return:
// - Runner if succes
// - Runner with id < 0 elsewhere
func (db *Database) GetRunner(runner_id int) Runner {

	var id int
	var status string
	var last_connection time.Time
	var ip string
	// Execute the query
	req := "SELECT * FROM runner WHERE runner.id='" + strconv.Itoa(runner_id) + "'"
	err := db.conn.QueryRow(req).Scan(&id, &status, &last_connection, &ip)
	if err != nil {
		log.Print("ERROR: Unable to select runner: ", err.Error())
		return Runner{-1, "", time.Time{}, ""}
	}

	return Runner{id, status, last_connection, ip}
}

// Add runner
// Return:
// - Runner id if succes
// - -1 elsewhere
func (db *Database) AddRunner(ip string) int64 {

	// Prepare statement for inserting data
	req := "INSERT INTO runner VALUES(NULL,'waiting',NULL,'" + ip + "')"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add runner: ", err.Error())
		return -1
	}
	defer insert.Close()

	res, err := insert.Exec()
	if err != nil {
		log.Print("ERROR: Unable to insert runner: ", err.Error())
		return -1
	}

	id, _ := res.LastInsertId()
	return id
}
