package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Bccp-Team/bccp-server/mysql"
	"github.com/Bccp-Team/bccp-server/scheduler"
	"github.com/gorilla/mux"

	. "github.com/Bccp-Team/bccp-server/proto/api"
)

type GitRequest struct {
	Repository struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
		HTTP        string `json:"git_http_url"`
		SSH         string `json:"git_ssh_url"`
		Visibility  int    `json:"visibility_level"`
	} `json:"repository"`
	Ref string
}

func PostCommitHandler(w http.ResponseWriter, r *http.Request) {
	var req GitRequest

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	err := decoder.Decode(&req)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		log.Printf("ERROR: api: ci: %v", err.Error())
		return
	}

	vars := mux.Vars(r)

	namespace := vars["namespace"]

	batch, err := mysql.Db.GetLastBatchFromNamespace(namespace)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		log.Printf("ERROR: api: ci: %v", err.Error())
		return
	}

	repo, err := mysql.Db.GetRepoFromName(req.Repository.Name, namespace)

	if !repo.Active {
		return
	}

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		log.Printf("ERROR: api: ci: %v", err.Error())
		return
	}

	repoId := strconv.FormatInt(repo.Id, 10)

	runs, err := mysql.Db.ListRuns(map[string]string{"repo": repoId,
		"status": "waiting"}, 0, 0)

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		log.Printf("ERROR: api: ci: %v", err.Error())
		return
	}

	if len(runs) > 0 {
		return
	}

	runID, err := mysql.Db.AddRun(repo.Id, batch.Id, 5)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		log.Printf("ERROR: api: ci: %v", err.Error())
		return
	}

	scheduler.DefaultScheduler.AddRun(&Run{Id: runID, Priority: 5})
}

// FIXME: ugly copy/past from the grpc api, because some peoples cant use grpc
func PushHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)

	encoder := json.NewEncoder(w)
	repoStr := vars["repo"]
	repos, err := mysql.Db.GetCiReposFromName(repoStr)

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	for _, repo := range repos {
		running, err := mysql.Db.ListRuns(map[string]string{
			"repo":   strconv.FormatInt(repo.Id, 10),
			"status": "waiting"},
			0, 0)

		if err != nil || len(running) > 0 {
			continue
		}

		batch, err := mysql.Db.GetLastBatchFromNamespace(repo.Namespace)

		if err != nil {
			//FIXME: log
			continue
		}

		runID, err := mysql.Db.AddRun(repo.Id, batch.Id, 5)

		if err != nil {
			//FIXME: log
			continue
		}

		run, err := mysql.Db.GetRun(runID)
		scheduler.DefaultScheduler.AddRun(run)
	}
}
