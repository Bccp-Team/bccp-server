package mysql

import (
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Run struct {
	Id        int
	Status    string
	Runner_id int
	Repo      string
	Logs      string
}

// List Runs
func (db *Database) ListRuns() ([]Run, error) {

	var id int
	var status string
	var runner_id int
	var repo string
	var logs string
	// Execute the query
	rows, err := db.conn.Query("SELECT * FROM run")
	if err != nil {
		log.Print("ERROR: Unable to select run: ", err.Error())
		return nil, err
	}

	var runs []Run

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(&id, &status, &runner_id, &repo, &logs)
		if err != nil {
			log.Print("ERROR: Unable to get next row: ", err.Error())
			return nil, err
		}

		runs = append(runs, Run{id, status, runner_id, repo, logs})
	}
	if err = rows.Err(); err != nil {
		log.Print("ERROR: Undefined row err: ", err.Error())
		return runs, err
	}

	return runs, nil
}

// Get runner info by id
func (db *Database) GetRun(run_id int) (Run, error) {

	var id int
	var status string
	var runner_id int
	var repo string
	var logs string
	// Execute the query
	req := "SELECT * FROM run WHERE run.id='" + strconv.Itoa(run_id) + "'"
	err := db.conn.QueryRow(req).Scan(&id, &status, &runner_id, &repo, &logs)
	if err != nil {
		log.Print("ERROR: Unable to select run: ", err.Error())
		return Run{-1, "", 0, "", ""}, err
	}

	return Run{id, status, runner_id, repo, logs}, nil
}

// Add run
func (db *Database) AddRun(repo string) (int, error) {

	// Prepare statement for inserting data
	req := "INSERT INTO run VALUES(NULL,'waiting',-1,'" + repo + "', '')"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add run: ", err.Error())
		return -1, err
	}
	defer insert.Close()

	res, err := insert.Exec()
	if err != nil {
		log.Print("ERROR: Unable to insert run: ", err.Error())
		return -1, err
	}

	id, _ := res.LastInsertId()
	return int(id), nil
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

func (db *Database) AddRun(depo string) (int, error) {
	req := "INSERT INTO run VALUES(NULL,'waiting',-1,?, '')"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add runner: ", err.Error())
		return -1, err
	}
	defer insert.Close()

	res, err := insert.Exec(depo)
	if err != nil {
		log.Print("ERROR: Unable to insert runner: ", err.Error())
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
		//FIXME error
		return err
	}

	_, err = update.Exec(new_logs)
	if err != nil {
		//FIXME error
		return err
	}

	return nil
}

func (db *Database) UpdateRunStatus(run_id int, status string) error {

	req := "UPDATE run SET status=? WHERE run.id='" + strconv.Itoa(run_id) + "'"

	update, err := db.conn.Prepare(req)
	defer update.Close()

	if err != nil {
		//FIXME error
		return err
	}

	_, err = update.Exec(status)
	if err != nil {
		//FIXME error
		return err
	}

	return nil
}

func (db *Database) UpdateRunRunner(run_id int, id int) error {

	req := "UPDATE run SET runner=? WHERE run.id='" + strconv.Itoa(run_id) + "'"

	update, err := db.conn.Prepare(req)
	defer update.Close()

	if err != nil {
		//FIXME error
		return err
	}

	_, err = update.Exec(id)
	if err != nil {
		//FIXME error
		return err
	}

	return nil
}
