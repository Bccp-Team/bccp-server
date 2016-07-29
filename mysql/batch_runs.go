package mysql

import (
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) GetBatchFromRun(run_id int) (*Batch, error) {

	var id int
	var batch_id int
	var r_id int
	// Execute the query
	req := "SELECT * FROM batch_runs WHERE batch_runs.run='" + strconv.Itoa(run_id) + "'"
	err := db.conn.QueryRow(req).Scan(&id, &batch_id, &r_id)
	if err != nil {
		log.Print("ERROR: Unable to select batch: ", err.Error())
		return nil, err
	}

	return db.GetBatch(batch_id)
}

func (db *Database) AddBatchRun(batch_id int, run_id int) error {

	// Prepare statement for inserting data
	req := "INSERT INTO batch_runs VALUES(NULL,?,?)"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add batch: ", err.Error())
		return err
	}
	defer insert.Close()

	_, err = insert.Exec(batch_id, run_id)
	if err != nil {
		log.Print("ERROR: Unable to insert batch: ", err.Error())
		return err
	}

	return nil
}
