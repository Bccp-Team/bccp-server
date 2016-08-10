package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bccp-server/mysql"
	"github.com/bccp-server/runners"
	"github.com/bccp-server/scheduler"
	"github.com/gorilla/mux"
)

// List all run
func GetRunHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Status string `json:"status"`
		Runner string `json:"runner"`
		Repo   string `json:"repo"`
	}

	var req request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		//FIXME error
	}

	filter := make(map[string]string)

	if len(req.Status) > 0 {
		filter["status"] = req.Status
	}
	if len(req.Runner) > 0 {
		filter["runner"] = req.Runner
	}
	if len(req.Repo) > 0 {
		filter["repo"] = req.Repo
	}

	runs, err := mysql.Db.ListRuns(filter)

	if err != nil {
		w.Write([]byte("{ \"error\" : \"unable to list runs\" }"))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(runs)
}

// Get information about given run
func GetRunByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.Write([]byte("{\"error\": \"wrong id\"}"))
		return
	}

	run, err := mysql.Db.GetRun(int(id))

	if err != nil {
		w.Write([]byte("{\"error\": \"the run does not exist\"}"))
		return
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
		w.Write([]byte("{\"error\": \"unable to parse request\"}"))
		return
	}

	if runReq.Namespace == "" || runReq.InitScript == "" || runReq.UpdateTime <= 0 || runReq.Timeout <= 0 {
		w.Write([]byte("{\"error\": \"missing fields\"}"))
		return
	}

	batch_id, err := mysql.Db.AddBatch(runReq.Namespace, runReq.InitScript,
		runReq.UpdateTime, runReq.Timeout)

	if err != nil {
		w.Write([]byte("{\"error\": \"unable to create batch\"}"))
		return
	}

	repos, err := mysql.Db.GetNamespaceRepos(&runReq.Namespace)

	if err != nil {
		w.Write([]byte("{\"error\": \"unable to list repos\"}"))
		return
	}

	for _, repo := range repos {
		runId, err := mysql.Db.AddRun(repo.Id, batch_id)
		if err != nil {
			w.Write([]byte("{\"error\": \"unable to add run\"}"))
			return
		}

		scheduler.DefaultScheduler.AddRun(runId)
	}
}

// Delete given runner
func DeleteRunHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.Write([]byte("{\"error\": \"wrong id\"}"))
		return
	}

	run, err := mysql.Db.GetRun(int(id))

	if err != nil {
		w.Write([]byte("{\"error\": \"the run does not exist\"}"))
		return
	}

	runners.KillRun(run.Runner_id, id)

	err = mysql.Db.UpdateRunStatus(id, "canceled")
}
