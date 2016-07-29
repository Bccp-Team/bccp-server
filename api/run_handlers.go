package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bccp-server/mysql"
	"github.com/bccp-server/scheduler"
	"github.com/gorilla/mux"
)

// List all run
func GetRunHandler(w http.ResponseWriter, r *http.Request) {
	runs, err := mysql.Db.ListRuns()

	if err != nil {
		w.Write([]byte("{ 'error' : 'unable to list runs' }"))
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(runs)
}

// Get information about given run
func GetRunByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.Write([]byte("{'error': 'wrong id'}"))
	}

	run, err := mysql.Db.GetRun(int(id))

	if err != nil {
		w.Write([]byte("{'error': 'the run does not exist'}"))
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(run)
}

// Add run

type runRequest struct {
	Namespace  string
	InitScript string
	UpdateTime int
	Timeout    int
}

func PutRunHandler(w http.ResponseWriter, r *http.Request) {
	var runReq runRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&runReq)

	if err != nil {
		w.Write([]byte("{'error': 'unable to parse request'}"))
		return
	}

	if runReq.Namespace == "" || runReq.InitScript == "" || runReq.UpdateTime <= 0 || runReq.Timeout <= 0 {
		w.Write([]byte("{'error': 'missing fields'}"))
		return
	}

	id, err := mysql.Db.AddBatch(runReq.Namespace, runReq.InitScript,
		runReq.UpdateTime, runReq.Timeout)

	if err != nil {
		w.Write([]byte("{'error': 'unable to create batch'}"))
		return
	}

	repos, err := mysql.Db.GetNamespaceRepos(runReq.Namespace)

	if err != nil {
		w.Write([]byte("{'error': 'unable to list repos'}"))
		return
	}

	for _, repo := range repos {
		runId, err := mysql.Db.AddRun(repo)
		if err != nil {
			w.Write([]byte("{'error': 'unable to add run'}"))
			return
		}

		err = mysql.Db.AddBatchRun(id, runId)

		if err != nil {
			w.Write([]byte("{'error': 'unable to add batch/run'}"))
			return
		}

		scheduler.DefaultScheduler.AddRun(runId)
	}
}

// Delete given runner
func DeleteRunHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}
