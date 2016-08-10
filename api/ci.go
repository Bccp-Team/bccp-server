package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bccp-server/mysql"
	"github.com/bccp-server/scheduler"
	"github.com/gorilla/mux"
)

type GitRequest struct {
	Repository struct {
		Name        string `json:name`
		Url         string `json:url`
		Description string `json:description`
		Homepage    string `json:homepage`
		Http        string `json:git_http_url`
		Ssh         string `json:git_ssh_url`
		Visibility  int    `json:visibility_level`
	} `json:repository`
	Ref string
}

func PostCommitHandler(w http.ResponseWriter, r *http.Request) {
	var req GitRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		//FIXME
		log.Printf(err.Error())
		return
	}

	vars := mux.Vars(r)

	namespace := vars["namespace"]

	batch, err := mysql.Db.GetLastBatchFromNamespace(namespace)

	if err != nil {
		//FIXME
		log.Printf(err.Error())
		return
	}

	repo, err := mysql.Db.GetRepoFromName(req.Repository.Name, namespace)

	if err != nil {
		//FIXME
		log.Printf(err.Error())
		return
	}

	runId, err := mysql.Db.AddRun(repo.Id, batch.Id)

	if err != nil {
		//FIXME
		log.Printf(err.Error())
		return
	}
	scheduler.DefaultScheduler.AddRun(runId)
}
