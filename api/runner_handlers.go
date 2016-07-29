package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bccp-server/mysql"
	"github.com/bccp-server/runners"
	"github.com/gorilla/mux"
)

// List all runners
func GetRunnerHandler(w http.ResponseWriter, r *http.Request) {
	runners := mysql.Db.ListRunners()
	encoder := json.NewEncoder(w)
	encoder.Encode(runners)
}

// Get information about given runner
func GetRunnerByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.Write([]byte("'error': 'wrong id'"))
	}

	runner, err := mysql.Db.GetRunner(int(id))

	if err != nil {
		w.Write([]byte("'error': 'the runner does not exist'"))
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(runner)
}

// Delete given runner
func DeleteRunnerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.Write([]byte("'error': 'wrong id'"))
	}

	runners.KillRunner(id)

	w.Write([]byte("'ok': 'killed'"))
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
