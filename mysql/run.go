package mysql

import (
	"database/sql"
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

func (db *Database) ListRuns() []Run {

	// Execute the query
	rows, err := db.conn.Query("SELECT * FROM run")
	if err != nil {
		log.Fatal("ERROR: Unable to select run: ", err.Error())
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal("ERROR: Unable to get columns: ", err.Error())
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var runs []Run

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal("ERROR: Unable to get next row: ", err.Error())
		}

		// Now fetch the data.
		id, err := strconv.Atoi(string(values[0]))
		if err != nil {
			log.Fatal("ERROR: Run id conversion error: ", err.Error())
		}
		status := string(values[1])
		runner_id, err := strconv.Atoi(string(values[2]))
		if err != nil {
			log.Fatal("ERROR: Runner id conversion error: ", err.Error())
		}
		repo := string(values[3])
		logs := string(values[4])
		runs = append(runs, Run{id, status, runner_id, repo, logs})
	}
	if err = rows.Err(); err != nil {
		log.Fatal("ERROR: Undefined row err: ", err.Error())
	}

	return runs
}

// Get runner info by id
// Return:
// - Run if succes
// - Run with id < 0 elsewhere
func (db *Database) GetRun(run_id int) Run {

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
		return Run{-1, "", 0, "", ""}
	}

	return Run{id, status, runner_id, repo, logs}
}

func (db *Database) UpdateRunLogs(run_id int, new_logs string) error {

	req := "UPDATE run SET logs=concat(logs, '?') WHERE run.id='" + strconv.Itoa(run_id) + "'"

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
