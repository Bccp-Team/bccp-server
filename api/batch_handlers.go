package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bccp-server/mysql"
	"github.com/gorilla/mux"
)

type runInfo struct {
	Id   int
	Repo string
}

type batchInfo struct {
	Id          int
	Namespace   string
	Init_script string
	Update_time int
	Timeout     int
	Runs        map[string]([]runInfo)
}

// Get information about given batch
func GetBatchByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.Write([]byte("{\"error\": \"wrong id\"}"))
		return
	}

	batch, err := mysql.Db.GetBatch(int(id))

	if err != nil {
		w.Write([]byte("{\"error\": \"unable to get batch\"}"))
		return
	}

	res := &batchInfo{batch.Id, batch.Namespace, batch.Init_script,
		batch.Update_time, batch.Timeout, make(map[string]([]runInfo))}

	for _, kind := range []string{"waiting", "running", "canceled",
		"finished", "failed", "timeout"} {
		runs, err := mysql.Db.ListBatchRuns(id, kind)

		if err != nil {
			w.Write([]byte("{\"error\": \"unable to get batch run\"}"))
		}

		runs_array := make([]runInfo, len(runs))

		for i, r := range runs {
			runs_array[i] = runInfo{r.Id, r.Repo}
		}

		res.Runs[kind] = runs_array
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}

// Delete given batch
func DeleteBatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}
