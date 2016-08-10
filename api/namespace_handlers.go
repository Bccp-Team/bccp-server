package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bccp-server/mysql"
	"github.com/gorilla/mux"
)

// List all namespaces
func GetNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	runs, err := mysql.Db.ListNamespaces()

	if err != nil {
		w.Write([]byte("{ \"error\" : \"unable to list namespaces\" }"))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(runs)
}

// Get information about given namespace
func GetNamespaceByNameHandler(w http.ResponseWriter, r *http.Request) {
	namespace := strings.Split(r.URL.Path, "/")[2]
	vars := mux.Vars(r)

	namespace = vars["name"]

	repos, err := mysql.Db.GetNamespaceRepos(&namespace)

	if err != nil {
		w.Write([]byte("{ 'error' : 'unable to list namespaces' }"))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(repos)
}

type namespace struct {
	Name  string
	Repos []struct {
		Repo string
		Ssh  string
	}
}

// Add namespace
func PutNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	var n namespace
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&n)

	if err != nil {
		w.Write([]byte("{ \"error\" : \"unable to decode namespace\" }"))
		return
	}

	err = mysql.Db.AddNamespace(n.Name)

	if err != nil {
		w.Write([]byte("{ \"error\" : \"unable to create namespace\" }"))
		return
	}

	for _, repo := range n.Repos {
		_, err = mysql.Db.AddRepoToNamespace(n.Name, repo.Repo, repo.Ssh)
		if err != nil {
			w.Write([]byte("{ \"error\" : \"unable to create repo\" }"))
			return
		}
	}
}

// Delete given namespace
func DeleteNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}
