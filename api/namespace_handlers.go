package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bccp-Team/bccp-server/mysql"
	"github.com/gorilla/mux"
)

// List all namespaces
func GetNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	runs, err := mysql.Db.ListNamespaces()
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	encoder.Encode(runs)
}

// Get information about given namespace
func GetNamespaceByNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)
	namespace := vars["name"]

	repos, err := mysql.Db.GetNamespaceRepos(&namespace)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	encoder.Encode(repos)
}

type repo struct {
	Repo string `json:"repo"`
	SSH  string `json:"ssh"`
}

// Add namespace
func PutNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	type namespace struct {
		Name  string `json:"namespace"`
		Repos []repo `json:"repos"`
	}

	var n namespace
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	err := decoder.Decode(&n)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	err = mysql.Db.AddNamespace(n.Name)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	for _, repo := range n.Repos {
		_, err = mysql.Db.AddRepoToNamespace(n.Name, repo.Repo, repo.SSH)
		if err != nil {
			encoder.Encode(map[string]string{"error": err.Error()})
			return
		}
	}

	encoder.Encode(map[string]string{"ok": "created"})
}
func AddRepoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	namespace := vars["name"]

	var rep repo

	err := decoder.Decode(&rep)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	id, err := mysql.Db.AddRepoToNamespace(namespace, rep.Repo, rep.SSH)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	encoder.Encode(map[string]string{"ok": strconv.Itoa(id)})
}

// Delete given namespace
func DeleteNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}
