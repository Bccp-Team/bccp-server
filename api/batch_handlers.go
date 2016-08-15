package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bccp-server/mysql"
	"github.com/gorilla/mux"
)

// Get information about given batch
func GetBatchsHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Namespace string `json:"namespace"`
	}

	var req request

	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	err := decoder.Decode(&req)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	var namespace *string
	if len(req.Namespace) == 0 {
		namespace = nil
	} else {
		namespace = &req.Namespace
	}

	batchs := mysql.Db.ListBatchs(namespace)
	encoder.Encode(batchs)
}

func GetActiveBatchsHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Namespace string `json:"namespace"`
	}

	var req request

	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	err := decoder.Decode(&req)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	var namespace *string
	if len(req.Namespace) == 0 {
		namespace = nil
	} else {
		namespace = &req.Namespace
	}

	batchs := mysql.Db.ListActiveBatchs(namespace)
	encoder.Encode(batchs)
}

func GetBatchByIdHandler(w http.ResponseWriter, r *http.Request) {
	type runInfo struct {
		Id   int `json:"id"`
		Repo int `json:"repo"`
	}

	type batchInfo struct {
		Id          int                    `json:"id"`
		Namespace   string                 `json:"namespace"`
		Init_script string                 `json:"init_script"`
		Update_time int                    `json:"update_time"`
		Timeout     int                    `json:"timeout"`
		Runs        map[string]([]runInfo) `json:"runs"`
	}

	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	batch, err := mysql.Db.GetBatch(int(id))

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	res := &batchInfo{batch.Id, batch.Namespace, batch.Init_script,
		batch.Update_time, batch.Timeout, make(map[string]([]runInfo))}

	for _, kind := range []string{"waiting", "running", "canceled",
		"finished", "failed", "timeout"} {
		runs, err := mysql.Db.ListRuns(map[string]string{"batch": string(id), "status": kind})

		if err != nil {
			encoder.Encode(map[string]string{"error": err.Error()})
			return
		}

		runs_array := make([]runInfo, len(runs))

		for i, r := range runs {
			runs_array[i] = runInfo{r.Id, r.Repo}
		}

		res.Runs[kind] = runs_array
	}
	encoder.Encode(res)
}

// Delete given batch
func DeleteBatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

func AddBatchHandler(w http.ResponseWriter, r *http.Request) {
	type runRequest struct {
		Namespace  string `json:"namespace"`
		InitScript string `json:"init_script"`
		UpdateTime int    `json:"update_time"`
		Timeout    int    `json:"timeout"`
	}

	var runReq runRequest

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

	_, err = mysql.Db.AddBatch(runReq.Namespace, runReq.InitScript,
		runReq.UpdateTime, runReq.Timeout)

	if err != nil {
		encoder.Encode(map[string]string{"error": "missing fields"})
		return
	}

	encoder.Encode(map[string]string{"ok": "created"})
}
