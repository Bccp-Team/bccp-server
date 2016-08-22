package mysql

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"
)

type Run struct {
	ID         int       `json:"id"`
	Status     string    `json:"status"`
	RunnerID   int       `json:"runner_id"`
	RunnerName string    `json:"runner_name"`
	Repo       int       `json:"repo"`
	RepoName   string    `json:"repo_name"`
	Batch      int       `json:"batch"`
	Namespace  string    `json:"namespace"`
	Logs       string    `json:"logs"`
	Creation   time.Time `json:"creation"`
	LastUpdate time.Time `json:"last_update"`
}

// List Runs
func (db *Database) ListRuns(filter map[string]string, limit, offset int) ([]Run, error) {
	var id int
	var status string
	var runnerID int
	var runnerName sql.NullString
	var repo int
	var repoName string
	var batch int
	var namespace string
	var logs string
	var creation time.Time
	var lastUpdate time.Time

	var rows *sql.Rows
	var err error

	limitReq := " ORDER BY run.last_update DESC"

	if limit > 0 {
		limitReq += " LIMIT " + strconv.Itoa(limit)
	}

	if offset > 0 {
		limitReq += " OFFSET " + strconv.Itoa(offset)
	}

	// Execute the query
	if len(filter) == 0 {
		rows, err = db.conn.Query("SELECT run.id, run.status, run.runner, runner.name, run.repo, namespace_repos.repo, run.batch, batch.namespace, run.logs, run.creation, run.last_update FROM run LEFT JOIN runner ON runner.id = run.runner JOIN namespace_repos ON run.repo = namespace_repos.id JOIN batch ON run.batch = batch.id" + limitReq)
	} else {
		req := "SELECT run.id, run.status, run.runner, runner.name, run.repo, namespace_repos.repo, run.batch, batch.namespace, run.logs, run.creation, run.last_update FROM run LEFT JOIN runner ON runner.id = run.runner JOIN namespace_repos ON run.repo = namespace_repos.id JOIN batch ON run.batch = batch.id WHERE "
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

	var runs []Run

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(&id, &status, &runnerID, &runnerName, &repo, &repoName, &batch, &namespace, &logs, &creation, &lastUpdate)
		if err != nil {
			log.Print("ERROR: Unable to get next row: ", err.Error())
			return nil, err
		}

		if !runnerName.Valid {
			runs = append(runs, Run{id, status, runnerID, "", repo, repoName, batch, namespace, logs, creation, lastUpdate})
		} else {
			runs = append(runs, Run{id, status, runnerID, runnerName.String, repo, repoName, batch, namespace, logs, creation, lastUpdate})
		}
	}
	if err = rows.Err(); err != nil {
		log.Print("ERROR: Undefined row err: ", err.Error())
		return runs, err
	}

	return runs, nil
}

// Get runner info by id
func (db *Database) GetRun(runID int) (*Run, error) {
	var id int
	var status string
	var runnerID int
	var runnerName sql.NullString
	var repo int
	var repoName string
	var batch int
	var namespace string
	var logs string
	var creation time.Time
	var lastUpdate time.Time

	// Execute the query
	req := "SELECT run.id, run.status, run.runner, runner.name, run.repo, namespace_repos.repo, run.batch, batch.namespace, run.logs, run.creation, run.last_update FROM run LEFT JOIN runner ON runner.id = run.runner JOIN namespace_repos ON run.repo = namespace_repos.id JOIN batch ON run.batch = batch.id WHERE run.id='" + strconv.Itoa(runID) + "'"
	err := db.conn.QueryRow(req).Scan(&id, &status, &runnerID, &runnerName, &repo, &repoName, &batch, &namespace, &logs, &creation, &lastUpdate)
	if err != nil {
		log.Print("ERROR: Unable to select run: ", err.Error())
		return nil, err
	}

	if !runnerName.Valid {
		return &Run{id, status, runnerID, "", repo, repoName, batch, namespace, logs, creation, lastUpdate}, nil
	}

	return &Run{id, status, runnerID, runnerName.String, repo, repoName, batch, namespace, logs, creation, lastUpdate}, nil
}

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
	req := "INSERT INTO run VALUES(NULL,'waiting',-1,?,?,'',NULL,NULL)"
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

func (db *Database) UpdateRunLogs(runID int, newLogs string) error {
	req := "UPDATE run SET logs=concat(logs, ?) WHERE run.id='" + strconv.Itoa(runID) + "'"

	update, err := db.conn.Prepare(req)
	defer update.Close()

	if err != nil {
		log.Print("ERROR: Unable to prepare update logs: ", err.Error())
		return err
	}

	_, err = update.Exec(newLogs)
	if err != nil {
		log.Print("ERROR: Unable to update logs: ", err.Error())
		return err
	}

	return nil
}
func (db *Database) StatRun(filter map[string]string) (stats map[string]int64, err error) {
	var total int64
	var waiting, running, canceled, finished, failed, timeout sql.NullInt64

	stats = make(map[string]int64)

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
		return
	}

	//FIXME refactor

	stats["all"] = total

	if waiting.Valid {
		stats["waiting"] = waiting.Int64
	} else {
		stats["waiting"] = 0
	}

	if running.Valid {
		stats["running"] = running.Int64
	} else {
		stats["running"] = 0
	}

	if canceled.Valid {
		stats["canceled"] = canceled.Int64
	} else {
		stats["canceled"] = 0
	}

	if finished.Valid {
		stats["finished"] = finished.Int64
	} else {
		stats["finished"] = 0
	}

	if failed.Valid {
		stats["failed"] = failed.Int64
	} else {
		stats["failed"] = 0
	}

	if timeout.Valid {
		stats["timeout"] = timeout.Int64
	} else {
		stats["timeout"] = 0
	}

	return
}
