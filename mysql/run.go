package mysql

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	. "github.com/Bccp-Team/bccp-server/proto/api"
)

// List Runs
func (db *Database) ListRuns(filter map[string]string, limit, offset int64) ([]*Run, error) {
	var id int64
	var status string
	var runnerID int64
	var runnerName sql.NullString
	var repo int64
	var repoName string
	var batch int64
	var namespace string
	var logs string
	var creation time.Time
	var lastUpdate time.Time
	var startTime time.Time

	var rows *sql.Rows
	var err error

	limitReq := " ORDER BY run.last_update DESC"

	if limit > 0 {
		limitReq += " LIMIT " + strconv.FormatInt(limit, 10)
	}

	if offset > 0 {
		limitReq += " OFFSET " + strconv.FormatInt(offset, 10)
	}

	// Execute the query
	if len(filter) == 0 {
		rows, err = db.conn.Query("SELECT run.id, run.status, run.runner, runner.name, run.repo, namespace_repos.repo, run.batch, batch.namespace, run.creation, run.last_update, run.start_time FROM run LEFT JOIN runner ON runner.id = run.runner JOIN namespace_repos ON run.repo = namespace_repos.id JOIN batch ON run.batch = batch.id" + limitReq)
	} else {
		req := "SELECT run.id, run.status, run.runner, runner.name, run.repo, namespace_repos.repo, run.batch, batch.namespace, run.creation, run.last_update, run.start_time FROM run LEFT JOIN runner ON runner.id = run.runner JOIN namespace_repos ON run.repo = namespace_repos.id JOIN batch ON run.batch = batch.id WHERE "
		f := make([]string, len(filter))
		i := 0
		l := make([]interface{}, len(filter))
		for key, value := range filter {
			//Here we trust that keys are legit
			f[i] = "run." + key + "=?"
			l[i] = value
			i = i + 1
		}
		rows, err = db.conn.Query(req+strings.Join(f, " AND ")+limitReq, l...)
	}
	if err != nil {
		log.Print("ERROR: Unable to select run: ", err.Error())
		return nil, err
	}

	var runs []*Run

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(&id, &status, &runnerID, &runnerName, &repo, &repoName, &batch, &namespace, &creation, &lastUpdate, &startTime)
		if err != nil {
			log.Print("ERROR: Unable to get next row: ", err.Error())
			return nil, err
		}

		runs = append(runs, &Run{id, status, runnerID, runnerName.String, repo, repoName, batch, namespace, logs, creation.String(), lastUpdate.String(), startTime.String(), lastUpdate.Sub(startTime).String()})
	}
	if err = rows.Err(); err != nil {
		log.Print("ERROR: Undefined row err: ", err.Error())
		return runs, err
	}

	return runs, nil
}

// Get runner info by id
func (db *Database) GetRun(runID int64) (*Run, error) {
	var id int64
	var status string
	var runnerID int64
	var runnerName sql.NullString
	var repo int64
	var repoName string
	var batch int64
	var namespace string
	var logs string
	var creation time.Time
	var lastUpdate time.Time
	var startTime time.Time

	// Execute the query
	req := "SELECT run.id, run.status, run.runner, runner.name, run.repo, namespace_repos.repo, run.batch, batch.namespace, run.logs, run.creation, run.last_update, run.start_time FROM run LEFT JOIN runner ON runner.id = run.runner JOIN namespace_repos ON run.repo = namespace_repos.id JOIN batch ON run.batch = batch.id WHERE run.id='" + strconv.FormatInt(runID, 10) + "'"
	err := db.conn.QueryRow(req).Scan(&id, &status, &runnerID, &runnerName, &repo, &repoName, &batch, &namespace, &logs, &creation, &lastUpdate, &startTime)
	if err != nil {
		log.Print("ERROR: Unable to select run: ", err.Error())
		return nil, err
	}

	return &Run{id, status, runnerID, runnerName.String, repo, repoName, batch, namespace, logs, creation.String(), lastUpdate.String(), startTime.String(), lastUpdate.Sub(startTime).String()}, nil
}

func (db *Database) LaunchRun(id int64, runner int64) error {
	// Prepare statement for inserting data
	req := "update run set status='running', runner=?, start_time=now() where id=?"
	update, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare: ", err.Error())
		return err
	}
	defer update.Close()

	_, err = update.Exec(runner, id)
	if err != nil {
		log.Print("ERROR: Unable to update run: ", err.Error())
		return err
	}

	return nil
}

func (db *Database) UpdateRunStatus(id int64, state string) error {
	// Prepare statement for inserting data
	req := "update run set status=? where id=?"
	update, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare: ", err.Error())
		return err
	}
	defer update.Close()

	_, err = update.Exec(state, id)
	if err != nil {
		log.Print("ERROR: Unable to update status: ", err.Error())
		return err
	}

	return nil
}

func (db *Database) AddRun(depo int64, batch int64) (int64, error) {
	req := "INSERT INTO run VALUES(NULL,'waiting',0,?,?,'',NULL,NULL,NULL)"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add run: ", err.Error())
		return -1, err
	}
	defer insert.Close()

	res, err := insert.Exec(depo, batch)
	if err != nil {
		log.Print("ERROR: Unable to insert run: ", err.Error())
		return -1, err
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func (db *Database) UpdateRunLogs(runID int64, newLogs string) error {
	req := "UPDATE run SET logs=concat(logs, ?) WHERE run.id=?"

	update, err := db.conn.Prepare(req)
	defer update.Close()

	if err != nil {
		log.Print("ERROR: Unable to prepare update logs: ", err.Error())
		return err
	}

	_, err = update.Exec(newLogs, runID)
	if err != nil {
		log.Print("ERROR: Unable to update logs: ", err.Error())
		return err
	}

	return nil
}
func (db *Database) StatRun(filter map[string]string) (*RunStats, error) {
	var total int64
	var err error
	var waiting, running, canceled, finished, failed, timeout sql.NullInt64

	req := `SELECT COUNT(*) total,
		SUM(CASE WHEN status = 'waiting' then 1 else 0 end) waiting,
		SUM(CASE WHEN status = 'running' then 1 else 0 end) running,
		SUM(CASE WHEN status = 'canceled' then 1 else 0 end) canceled,
		SUM(CASE WHEN status = 'finished' then 1 else 0 end) finished,
		SUM(CASE WHEN status = 'failed' then 1 else 0 end) failed,
		SUM(CASE WHEN status = 'timeout' then 1 else 0 end) timeout
		FROM run`

	if len(filter) == 0 {
		err = db.conn.QueryRow(req).Scan(&total, &waiting, &running, &canceled, &finished, &failed, &timeout)
	} else {
		req += " WHERE "
		f := make([]string, len(filter))
		i := 0
		l := make([]interface{}, len(filter))
		for key, value := range filter {
			//Here we trust that keys are legit
			f[i] = "run." + key + "=?"
			l[i] = value
			i = i + 1
		}
		err = db.conn.QueryRow(req+strings.Join(f, " AND "), l...).Scan(&total, &waiting, &running, &canceled, &finished, &failed, &timeout)
	}

	if err != nil {
		log.Print("ERROR: Unable to select run: ", err.Error())
		return nil, err
	}

	//FIXME refactor

	return &RunStats{
		total,
		waiting.Int64,
		running.Int64,
		canceled.Int64,
		finished.Int64,
		failed.Int64,
		timeout.Int64,
	}, nil
}
