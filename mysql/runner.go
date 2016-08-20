package mysql

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Runner struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	Last_connection time.Time `json:"last_connection"`
	Ip              string    `json:"ip"`
}

func (db *Database) ListRunners(filter map[string]string, limit, offset int) []Runner {
	var rows *sql.Rows
	var err error

	limit_req := " ORDER BY last_conn DESC"

	if limit > 0 {
		limit_req += " LIMIT " + strconv.Itoa(limit)
	}

	if offset > 0 {
		limit_req += " OFFSET " + strconv.Itoa(offset)
	}

	// Execute the query
	if len(filter) == 0 {
		rows, err = db.conn.Query("SELECT * FROM runner" + limit_req)
	} else {
		req := "SELECT * FROM runner WHERE "
		f := make([]string, len(filter))
		i := 0
		l := make([]interface{}, len(filter))
		for key, value := range filter {
			//Here we trust that keys are legit
			f[i] = key + "=?"
			l[i] = value
			i = i + 1
		}
		rows, err = db.conn.Query(req+strings.Join(f, " AND ")+limit_req, l...)
	}

	if err != nil {
		log.Fatal("ERROR: Unable to select runner: ", err.Error())
	}

	var runners []Runner

	// Fetch rows
	for rows.Next() {
		var id int
		var name string
		var status string
		var last_connection time.Time
		var ip string
		// get RawBytes from data
		err = rows.Scan(&id, &name, &status, &last_connection, &ip)
		if err != nil {
			log.Fatal("ERROR: Unable to get next row: ", err.Error())
		}

		runners = append(runners, Runner{id, name, status, last_connection, ip})
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
func (db *Database) GetRunner(runner_id int) (*Runner, error) {

	var id int
	var name string
	var status string
	var last_connection time.Time
	var ip string
	// Execute the query
	req := "SELECT * FROM runner WHERE runner.id='" + strconv.Itoa(runner_id) + "'"
	err := db.conn.QueryRow(req).Scan(&id, &name, &status, &last_connection, &ip)
	if err != nil {
		log.Print("ERROR: Unable to select runner: ", err.Error())
		return nil, err
	}

	return &Runner{id, name, status, last_connection, ip}, nil
}

// Add runner
// Return:
// - Runner id if succes
// - -1 elsewhere
func (db *Database) AddRunner(ip string, name string) (int, error) {

	// Prepare statement for inserting data
	req := "INSERT INTO runner VALUES(NULL, ?, 'waiting', NULL, ?)"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add runner: ", err.Error())
		return -1, err
	}
	defer insert.Close()

	res, err := insert.Exec(name, ip)
	if err != nil {
		log.Print("ERROR: Unable to insert runner: ", err.Error())
		return -1, err
	}

	id, _ := res.LastInsertId()
	return int(id), nil
}

// Add runner
// Return:
// - Runner id if succes
// - -1 elsewhere
func (db *Database) UpdateRunner(id int, state string) error {

	// Prepare statement for inserting data
	req := "update runner set status='" + state + "' where id=" + strconv.Itoa(id)
	update, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add runner: ", err.Error())
		return err
	}
	defer update.Close()

	_, err = update.Exec()
	if err != nil {
		log.Print("ERROR: Unable to update runner: ", err.Error())
		return err
	}

	return nil
}

func (db *Database) StatRunners() (total, waiting, dead int64, err error) {
	// Execute the query

	var waitingNull sql.NullInt64
	var deadNull sql.NullInt64

	req := "SELECT COUNT(*) total, SUM(CASE WHEN status = 'waiting' then 1 else 0 end) waiting, SUM(CASE WHEN status = 'dead' then 1 else 0 end) dead FROM runner"

	err = db.conn.QueryRow(req).Scan(&total, &waitingNull, &deadNull)

	if err != nil {
		log.Print("ERROR: Unable to select run: ", err.Error())
		return
	}

	if waitingNull.Valid {
		waiting = waitingNull.Int64
	}

	if deadNull.Valid {
		dead = deadNull.Int64
	}

	return
}
