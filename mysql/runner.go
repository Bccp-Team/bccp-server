package mysql

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Runner struct {
	Id     int
	Status string
}

func (db *Database) ListRunners() []Runner {

	// Execute the query
	rows, err := db.conn.Query("SELECT * FROM runner")
	if err != nil {
		log.Fatal("ERROR: Unable to select runner: ", err.Error())
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

	var runners []Runner

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
			log.Fatal("ERROR: Runner id conversion error: ", err.Error())
		}
		status := string(values[1])
		runners = append(runners, Runner{id, status})
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
	// Execute the query
	req := "SELECT * FROM runner WHERE runner.id='" + strconv.Itoa(runner_id) + "'"
	err := db.conn.QueryRow(req).Scan(&id, &status)
	if err != nil {
		log.Print("ERROR: Unable to select runner: ", err.Error())
		return Runner{-1, ""}
	}

	return Runner{id, status}
}
