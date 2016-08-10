package mysql

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Batch struct {
	Id          int
	Namespace   string
	Init_script string
	Update_time int
	Timeout     int
}

func (db *Database) ListBatchs(namespace *string) []Batch {
	var rows *sql.Rows
	var err error

	if namespace == nil {
		rows, err = db.conn.Query("SELECT * FROM batch")
	} else {
		rows, err = db.conn.Query("SELECT * FROM batch where namespace=?", namespace)
	}
	// Execute the query
	if err != nil {
		log.Fatal("ERROR: Unable to select batch: ", err.Error())
	}

	var batchs []Batch

	// Fetch rows
	for rows.Next() {
		var id int
		var namespace string
		var init_script string
		var update_time int
		var timeout int
		// get RawBytes from data
		err = rows.Scan(&id, &namespace, &init_script, &update_time, &timeout)
		if err != nil {
			log.Fatal("ERROR: Unable to get next row: ", err.Error())
		}

		batchs = append(batchs, Batch{id, namespace, init_script, update_time, timeout})
	}
	if err = rows.Err(); err != nil {
		log.Fatal("ERROR: Undefined row err: ", err.Error())
	}

	return batchs
}

func (db *Database) GetBatch(id int) (*Batch, error) {

	var b_id int
	var namespace string
	var init_script string
	var update_time int
	var timeout int
	// Execute the query
	req := "SELECT * FROM batch WHERE batch.id='" + strconv.Itoa(id) + "'"
	err := db.conn.QueryRow(req).Scan(&b_id, &namespace, &init_script, &update_time, &timeout)
	if err != nil {
		log.Print("ERROR: Unable to select batch: ", err.Error())
		return nil, err
	}

	return &Batch{id, namespace, init_script, update_time, timeout}, nil
}

func (db *Database) GetLastBatchFromNamespace(n string) (*Batch, error) {

	var b_id int
	var namespace string
	var init_script string
	var update_time int
	var timeout int
	// Execute the query
	req := "SELECT * FROM batch WHERE batch.namespace=? ORDER BY id DESC"
	err := db.conn.QueryRow(req, n).Scan(&b_id, &namespace, &init_script, &update_time, &timeout)
	if err != nil {
		log.Print("ERROR: Unable to select batch: ", err.Error())
		return nil, err
	}

	return &Batch{b_id, namespace, init_script, update_time, timeout}, nil
}

func (db *Database) AddBatch(namespace string, init_script string, update_time int, timeout int) (int, error) {

	// Prepare statement for inserting data
	req := "INSERT INTO batch VALUES(NULL,?,?,?,?)"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add batch: ", err.Error())
		return -1, err
	}
	defer insert.Close()

	res, err := insert.Exec(namespace, init_script, update_time, timeout)
	if err != nil {
		log.Print("ERROR: Unable to insert batch: ", err.Error())
		return -1, err
	}

	id, _ := res.LastInsertId()
	return int(id), nil
}
