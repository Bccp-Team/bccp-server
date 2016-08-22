package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Bccp-Team/bccp-server/mysql"
	"github.com/Bccp-Team/bccp-server/scheduler"
	"github.com/gorilla/mux"
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
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		log.Printf("ERROR: api: ci: %v", err.Error())
		return
	}

	runID, err := mysql.Db.AddRun(repo.ID, batch.ID)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		log.Printf("ERROR: api: ci: %v", err.Error())
		return
	}

	scheduler.DefaultScheduler.AddRun(runID)
}
