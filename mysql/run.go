package mysql

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Run struct {
	Id        int    `json:"id"`
	Status    string `json:"status"`
	Runner_id int    `json:"runner_id"`
	Repo      int    `json:"repo"`
	Batch     int    `json:"batch`
	Logs      string `json:"logs"`
}

// List Runs
func (db *Database) ListRuns(filter map[string]string) ([]Run, error) {

	var id int
	var status string
	var runner_id int
	var repo int
	var batch int
	var logs string

	var rows *sql.Rows
	var err error
	// Execute the query
	if len(filter) == 0 {
		rows, err = db.conn.Query("SELECT * FROM run")
	} else {
		req := "SELECT * FROM run WHERE "
		f := make([]string, len(filter))
		i := 0
		l := make([]interface{}, len(filter))
		for key, value := range filter {
			//Here we trust that keys are legit
			f[i] = key + "=?"
			l[i] = value
			i = i + 1
		}
		rows, err = db.conn.Query(req+strings.Join(f, " AND "), l...)
	}
	if err != nil {
		log.Print("ERROR: Unable to select run: ", err.Error())
		return nil, err
	}

	var runs []Run

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(&id, &status, &runner_id, &repo, &batch, &logs)
		if err != nil {
			log.Print("ERROR: Unable to get next row: ", err.Error())
			return nil, err
		}

		runs = append(runs, Run{id, status, runner_id, repo, batch, logs})
	}
	if err = rows.Err(); err != nil {
		log.Print("ERROR: Undefined row err: ", err.Error())
		return runs, err
	}

	return runs, nil
}

// Get runner info by id
func (db *Database) GetRun(run_id int) (*Run, error) {

	var id int
	var status string
	var runner_id int
	var repo int
	var batch int
	var logs string
	// Execute the query
	req := "SELECT * FROM run WHERE run.id='" + strconv.Itoa(run_id) + "'"
	err := db.conn.QueryRow(req).Scan(&id, &status, &runner_id, &repo, &batch, &logs)
	if err != nil {
		log.Print("ERROR: Unable to select run: ", err.Error())
		return nil, err
	}

	return &Run{id, status, runner_id, repo, batch, logs}, nil
}

// Launch run
func (db *Database) LaunchRun(id int, runner int) error {

	// Prepare statement for inserting data
	req := "update run set status='running', runner=" + strconv.Itoa(runner)
	req += " where id=" + strconv.Itoa(id)
	update, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare: ", err.Error())
		return err
	}
	defer update.Close()

	_, err = update.Exec()
	if err != nil {
		log.Print("ERROR: Unable to update run: ", err.Error())
		return err
	}

	return nil
}

// Update
func (db *Database) UpdateRunStatus(id int, state string) error {

	// Prepare statement for inserting data
	req := "update run set status='" + state + "' where id=" + strconv.Itoa(id)
	update, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare: ", err.Error())
		return err
	}
	defer update.Close()

	_, err = update.Exec()
	if err != nil {
		log.Print("ERROR: Unable to update status: ", err.Error())
		return err
	}

	return nil
}

func (db *Database) AddRun(depo int, batch int) (int, error) {
	req := "INSERT INTO run VALUES(NULL,'waiting',-1,?,?,'')"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add runner: ", err.Error())
		return -1, err
	}
	defer insert.Close()

	res, err := insert.Exec(depo, batch)
	if err != nil {
		log.Print("ERROR: Unable to insert run: ", err.Error())
		return -1, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (db *Database) UpdateRunLogs(run_id int, new_logs string) error {

	req := "UPDATE run SET logs=concat(logs, ?) WHERE run.id='" + strconv.Itoa(run_id) + "'"

	update, err := db.conn.Prepare(req)
	defer update.Close()

	if err != nil {
		log.Print("ERROR: Unable to prepare update logs: ", err.Error())
		return err
	}

	_, err = update.Exec(new_logs)
	if err != nil {
		log.Print("ERROR: Unable to update logs: ", err.Error())
		return err
	}

	return nil
}
