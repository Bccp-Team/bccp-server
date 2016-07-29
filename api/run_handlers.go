package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bccp-server/mysql"
	"github.com/gorilla/mux"
)

// List all run
func GetRunHandler(w http.ResponseWriter, r *http.Request) {
	runs, err := mysql.Db.ListRuns()

	if err != nil {
		w.Write([]byte("'error' : 'unable to list runs'"))
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(runs)
}

// Get information about given run
func GetRunByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.Write([]byte("'error': 'wrong id'"))
	}

	run, err := mysql.Db.GetRun(int(id))

	if err != nil {
		w.Write([]byte("'error': 'the run does not exist'"))
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(run)
}

// Add run
func PutRunHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Delete given runner
func DeleteRunHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}
