package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bccp-Team/bccp-server/mysql"
	"github.com/gorilla/mux"
)

// Get information about given batch
func GetBatchsHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Namespace string `json:"namespace"`
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
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

	var namespace *string
	if len(req.Namespace) == 0 {
		namespace = nil
	} else {
		namespace = &req.Namespace
	}

	batches := mysql.Db.ListBatchs(namespace, req.Limit, req.Offset)
	encoder.Encode(batches)
}

func GetActiveBatchsHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Namespace string `json:"namespace"`
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
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

	var namespace *string
	if len(req.Namespace) == 0 {
		namespace = nil
	} else {
		namespace = &req.Namespace
	}

	batches := mysql.Db.ListActiveBatches(namespace, req.Limit, req.Offset)
	encoder.Encode(batches)
}

func GetBatchByIDHandler(w http.ResponseWriter, r *http.Request) {
	type batchInfo struct {
		ID         int    `json:"id"`
		Namespace  string `json:"namespace"`
		InitScript string `json:"init_script"`
		UpdateTime int    `json:"update_time"`
		Timeout    int    `json:"timeout"`
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

	res := &batchInfo{batch.ID, batch.Namespace, batch.InitScript,
		batch.UpdateTime, batch.Timeout}
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

	_, err = mysql.Db.AddBatch(runReq.Namespace, runReq.InitScript, runReq.UpdateTime, runReq.Timeout)
	if err != nil {
		encoder.Encode(map[string]string{"error": "missing fields"})
		return
	}

	encoder.Encode(map[string]string{"ok": "created"})
}

func GetBatchStatHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Namespace string `json:"namespace"`
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

	var namespace *string

	if len(req.Namespace) > 0 {
		namespace = &req.Namespace
	}

	stats, err := mysql.Db.StatBatch(namespace)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	encoder.Encode(stats)
}
