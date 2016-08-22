package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bccp-Team/bccp-server/mysql"
	"github.com/Bccp-Team/bccp-server/runners"
	"github.com/gorilla/mux"
)

// List all runners
func GetRunnerHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Status string `json:"status"`
		Name   string `json:"name"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
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
	if len(req.Name) > 0 {
		filter["name"] = req.Name
	}

	dbRunners := mysql.Db.ListRunners(filter, req.Limit, req.Offset)
	encoder.Encode(dbRunners)
}

func GetRunnerStatHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	total, waiting, dead, _ := mysql.Db.StatRunners()
	encoder.Encode(map[string]int64{"all": total, "waiting": waiting, "dead": dead})
}

// Get information about given runner
func GetRunnerByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	runner, err := mysql.Db.GetRunner(int(id))
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	encoder.Encode(runner)
}

// Delete given runner
func DeleteRunnerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	runners.KillRunner(id)
	encoder.Encode(map[string]string{"ok": "killed"})
}

// FIXME: we should do something clever to avoid race conditions

// Enable given runner
func PostEnableRunnerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Disable given runner
func PostDisableRunnerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}
