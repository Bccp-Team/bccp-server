package mysql

import (
	"database/sql"
	"log"
	"strconv"
	"time"
)

type Batch struct {
	ID         int       `json:"id"`
	Namespace  string    `json:"namespace"`
	InitScript string    `json:"init_script"`
	UpdateTime int       `json:"update_time"`
	Timeout    int       `json:"timeout"`
	Creation   time.Time `json:"time"`
}

func (db *Database) ListBatchs(namespace *string, limit, offset int) []Batch {
	var rows *sql.Rows
	var err error

	limitReq := " ORDER BY creation DESC"

	if limit > 0 {
		limitReq += " LIMIT " + strconv.Itoa(limit)
	}

	if offset > 0 {
		limitReq += " OFFSET " + strconv.Itoa(offset)
	}

	if namespace == nil {
		rows, err = db.conn.Query("SELECT * FROM batch" + limitReq)
	} else {
		rows, err = db.conn.Query("SELECT * FROM batch where namespace=?"+limitReq, namespace)
	}
	// Execute the query
	if err != nil {
		log.Fatal("ERROR: Unable to select batch: ", err.Error())
	}

	var batches []Batch

	// Fetch rows
	for rows.Next() {
		var id int
		var namespace string
		var initScript string
		var updateTime int
		var timeout int
		var creation time.Time

		// get RawBytes from data
		err = rows.Scan(&id, &namespace, &initScript, &updateTime, &timeout, &creation)
		if err != nil {
			log.Fatal("ERROR: Unable to get next row: ", err.Error())
		}

		batches = append(batches, Batch{id, namespace, initScript, updateTime, timeout, creation})
	}

	if err = rows.Err(); err != nil {
		log.Fatal("ERROR: Undefined row err: ", err.Error())
	}

	return batches
}

func (db *Database) ListActiveBatches(namespace *string, limit, offset int) []Batch {
	var rows *sql.Rows
	var err error

	limitReq := " ORDER BY creation DESC"

	if limit > 0 {
		limitReq += " LIMIT " + strconv.Itoa(limit)
	}

	if offset > 0 {
		limitReq += " OFFSET " + strconv.Itoa(offset)
	}

	if namespace == nil {
		rows, err = db.conn.Query("SELECT * FROM batch WHERE EXISTS(SELECT * FROM run WHERE run.batch = batch.id AND run.status IN ('waiting', 'running'))" + limitReq)
	} else {
		rows, err = db.conn.Query("SELECT * FROM batch WHERE namespace=? AND EXISTS(SELECT * FROM run WHERE run.batch = batch.id AND run.status IN ('waiting', 'running'))"+limitReq, namespace)
	}
	// Execute the query
	if err != nil {
		log.Fatal("ERROR: Unable to select batch: ", err.Error())
	}

	var batches []Batch

	// Fetch rows
	for rows.Next() {
		var id int
		var namespace string
		var initScript string
		var updateTime int
		var timeout int
		var creation time.Time

		// get RawBytes from data
		err = rows.Scan(&id, &namespace, &initScript, &updateTime, &timeout, &creation)
		if err != nil {
			log.Fatal("ERROR: Unable to get next row: ", err.Error())
		}

		batches = append(batches, Batch{id, namespace, initScript, updateTime, timeout, creation})
	}
	if err = rows.Err(); err != nil {
		log.Fatal("ERROR: Undefined row err: ", err.Error())
	}

	return batches
}

func (db *Database) GetBatch(id int) (*Batch, error) {
	var bID int
	var namespace string
	var initScript string
	var updateTime int
	var timeout int
	var creation time.Time

	// Execute the query
	req := "SELECT * FROM batch WHERE batch.id='" + strconv.Itoa(id) + "'"
	err := db.conn.QueryRow(req).Scan(&bID, &namespace, &initScript, &updateTime, &timeout, &creation)
	if err != nil {
		log.Print("ERROR: Unable to select batch: ", err.Error())
		return nil, err
	}

	return &Batch{id, namespace, initScript, updateTime, timeout, creation}, nil
}

func (db *Database) GetLastBatchFromNamespace(n string) (*Batch, error) {
	var bID int
	var namespace string
	var initScript string
	var updateTime int
	var timeout int
	var creation time.Time

	// Execute the query
	req := "SELECT * FROM batch WHERE batch.namespace=? ORDER BY creation DESC"
	err := db.conn.QueryRow(req, n).Scan(&bID, &namespace, &initScript, &updateTime, &timeout, &creation)
	if err != nil {
		log.Print("ERROR: Unable to select batch: ", err.Error())
		return nil, err
	}

	return &Batch{bID, namespace, initScript, updateTime, timeout, creation}, nil
}

func (db *Database) AddBatch(namespace string, initScript string, updateTime int, timeout int) (int, error) {
	// Prepare statement for inserting data
	req := "INSERT INTO batch VALUES(NULL,?,?,?,?,NULL)"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add batch: ", err.Error())
		return -1, err
	}
	defer insert.Close()

	res, err := insert.Exec(namespace, initScript, updateTime, timeout)
	if err != nil {
		log.Print("ERROR: Unable to insert batch: ", err.Error())
		return -1, err
	}

	id, _ := res.LastInsertId()
	return int(id), nil
}

func (db *Database) StatBatch(namespace *string) (stats map[string]int64, err error) {
	stats = make(map[string]int64)
	var total int64
	var active sql.NullInt64

	req := `SELECT COUNT(*) total,
		SUM(CASE WHEN EXISTS(SELECT * FROM run
					WHERE run.batch = batch.id
					AND run.status IN ('waiting', 'running'))
			THEN 1
			ELSE 0
		    END) active
		FROM batch`

	if namespace == nil {
		err = db.conn.QueryRow(req).Scan(&total, &active)
	} else {
		req += "WHERE namespace=?"
		err = db.conn.QueryRow(req, *namespace).Scan(&total, &active)
	}

	if err != nil {
		log.Print("ERROR: Unable to stat batch: ", err.Error())
		return
	}

	stats["all"] = total

	if active.Valid {
		stats["active"] = active.Int64
	} else {
		stats["active"] = 0
	}

	return
}
