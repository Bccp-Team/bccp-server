package rest

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

func GetPingHandler(w http.ResponseWriter, r *http.Request) {
}

func SetupRestAPI(wait *sync.WaitGroup, port string, crtFile string, keyFile string) {
	r := mux.NewRouter()

	// Define supported methods
	pstAPI := r.Methods("POST").Subrouter()

	/*
		getAPI := r.Methods("GET").Subrouter()
		putAPI := r.Methods("PUT").Subrouter()
		delAPI := r.Methods("DELETE").Subrouter()

		// Define routes
		getAPI.HandleFunc("/ping", GetPingHandler)

		getAPI.HandleFunc("/runner", GetRunnerHandler)
		getAPI.HandleFunc("/runner/stats", GetRunnerStatHandler)
		getAPI.HandleFunc("/runner/{id:[0-9]+}", GetRunnerByIDHandler)
		delAPI.HandleFunc("/runner/{id:[0-9]+}", DeleteRunnerHandler)
		pstAPI.HandleFunc("/runner/{id:[0-9]+}/enable", PostEnableRunnerHandler)
		pstAPI.HandleFunc("/runner/{id:[0-9]+}/disable", PostDisableRunnerHandler)

		getAPI.HandleFunc("/run", GetRunHandler)
		getAPI.HandleFunc("/run/stats", GetRunStatHandler)
		getAPI.HandleFunc("/run/{id:[0-9]+}", GetRunByIDHandler)
		putAPI.HandleFunc("/run/{batch_id:[0-9]+}/{repo_id:[0-9]+}", PutRunRepoHandler)
		putAPI.HandleFunc("/run", PutRunHandler)
		delAPI.HandleFunc("/run/{id:[0-9]+}", DeleteRunHandler)

		putAPI.HandleFunc("/batch", AddBatchHandler)
		getAPI.HandleFunc("/batch", GetBatchsHandler)
		getAPI.HandleFunc("/batch/stats", GetBatchStatHandler)
		getAPI.HandleFunc("/batch/active", GetActiveBatchsHandler)
		getAPI.HandleFunc("/batch/{id:[0-9]+}", GetBatchByIDHandler)
		delAPI.HandleFunc("/batch/{id:[0-9]+}", DeleteBatchHandler)

		getAPI.HandleFunc("/namespace", GetNamespaceHandler)
		getAPI.HandleFunc("/namespace/{name:[--~]+}", GetNamespaceByNameHandler)
		putAPI.HandleFunc("/namespace", PutNamespaceHandler)
		delAPI.HandleFunc("/namespace/{name:[--~]+}", AddRepoHandler)
		delAPI.HandleFunc("/namespace/{name:[--~]+}", DeleteNamespaceHandler)
	*/

	pstAPI.HandleFunc("/ci/{namespace:[--~]+}", PostCommitHandler)

	// Launch async server with router
	var err error
	(*wait).Add(1)
	go func() {
		http.Handle("/", r)
		log.Print("INFO: Launching http server with crt: '" + crtFile + "'")
		log.Print("      and key: '" + keyFile + "'")
		err = http.ListenAndServeTLS(":"+port, crtFile, keyFile, nil)
		defer (*wait).Done()
	}()
	time.Sleep(1 * time.Second)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
