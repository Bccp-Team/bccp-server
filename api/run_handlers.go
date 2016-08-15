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
	encoder := json.NewEncoder(w)

	err := decoder.Decode(&req)

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
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
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	encoder.Encode(runs)
}

// Get information about given run
func GetRunByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	run, err := mysql.Db.GetRun(int(id))

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	encoder.Encode(run)
}

func PutRunRepoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	batch_id, _ := strconv.Atoi(vars["batch_id"])
	repo_id, _ := strconv.Atoi(vars["repo_id"])

	runId, err := mysql.Db.AddRun(repo_id, batch_id)

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	scheduler.DefaultScheduler.AddRun(runId)
	encoder.Encode(map[string]string{"ok": string(runId)})
}

// Add run

func PutRunHandler(w http.ResponseWriter, r *http.Request) {
	type runRequest struct {
		Namespace  string `json:"namespace"`
		InitScript string `json:"init_script"`
		UpdateTime int    `json:"update_time"`
		Timeout    int    `json:"timeout"`
	}

	type runResult struct {
		Id   int            `json:"id"`
		Runs map[int]string `json:"runs"`
	}

	var runReq runRequest
	var runRes runResult

	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	err := decoder.Decode(&runReq)

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	if runReq.Namespace == "" || runReq.InitScript == "" || runReq.UpdateTime <= 0 || runReq.Timeout <= 0 {
		encoder.Encode(map[string]string{"error": "missing fields"})
		return
	}

	batch_id, err := mysql.Db.AddBatch(runReq.Namespace, runReq.InitScript,
		runReq.UpdateTime, runReq.Timeout)

	runRes.Id = batch_id

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	repos, err := mysql.Db.GetNamespaceRepos(&runReq.Namespace)

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	runRes.Runs = make(map[int]string)

	for _, repo := range repos {
		runId, err := mysql.Db.AddRun(repo.Id, batch_id)
		if err != nil {
			encoder.Encode(map[string]string{"error": err.Error()})
			return
		}

		scheduler.DefaultScheduler.AddRun(runId)
		runRes.Runs[runId] = repo.Repo
	}

	encoder.Encode(runRes)
}

// Delete given runner
func DeleteRunHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	run, err := mysql.Db.GetRun(int(id))

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	runners.KillRun(run.Runner_id, id)

	err = mysql.Db.UpdateRunStatus(id, "canceled")
	encoder.Encode(map[string]string{"ok": "canceled"})
}
