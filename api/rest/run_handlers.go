package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bccp-Team/bccp-server/mysql"
	"github.com/Bccp-Team/bccp-server/runners"
	"github.com/Bccp-Team/bccp-server/scheduler"
	"github.com/gorilla/mux"
)

// List all run
func GetRunHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Status string `json:"status"`
		Runner string `json:"runner"`
		Repo   string `json:"repo"`
		Batch  string `json:"batch"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
	}

	var req request

	defer r.Body.Close()
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
	if len(req.Batch) > 0 {
		filter["batch"] = req.Batch
	}

	runs, err := mysql.Db.ListRuns(filter, req.Limit, req.Offset)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	encoder.Encode(runs)
}

func GetRunStatHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Status string `json:"status"`
		Runner string `json:"runner"`
		Repo   string `json:"repo"`
		Batch  string `json:"batch"`
	}

	var req request

	defer r.Body.Close()
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
	if len(req.Batch) > 0 {
		filter["batch"] = req.Batch
	}

	stats, err := mysql.Db.StatRun(filter)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	encoder.Encode(stats)
}

// Get information about given run
func GetRunByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	run, err := mysql.Db.GetRun(int64(id))
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	encoder.Encode(run)
}

func PutRunRepoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	batchID, _ := strconv.Atoi(vars["batch_id"])
	repoID, _ := strconv.Atoi(vars["repo_id"])

	runID, err := mysql.Db.AddRun(repoID, batchID)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	scheduler.DefaultScheduler.AddRun(runID)
	encoder.Encode(map[string]string{"ok": strconv.Itoa(runID)})
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
		ID   int            `json:"id"`
		Runs map[int]string `json:"runs"`
	}

	var runReq runRequest
	var runRes runResult

	defer r.Body.Close()
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

	batchID, err := mysql.Db.AddBatch(runReq.Namespace, runReq.InitScript, runReq.UpdateTime, runReq.Timeout)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	runRes.ID = batchID
	repos, err := mysql.Db.GetNamespaceRepos(&runReq.Namespace)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	runRes.Runs = make(map[int]string)

	for _, repo := range repos {
		runID, err := mysql.Db.AddRun(repo.ID, batchID)
		if err != nil {
			encoder.Encode(map[string]string{"error": err.Error()})
			return
		}

		scheduler.DefaultScheduler.AddRun(runID)
		runRes.Runs[runID] = repo.Repo
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

	if run.RunnerID != 0 {
		runners.KillRun(run.RunnerID, id)
	}

	err = mysql.Db.UpdateRunStatus(id, "canceled")
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	encoder.Encode(map[string]string{"ok": "canceled"})
}
