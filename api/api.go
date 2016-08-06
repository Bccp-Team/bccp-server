package api

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

func SetupRestAPI(wait *sync.WaitGroup, port string, crt_file string, key_file string) {
	r := mux.NewRouter()

	// Define supported methods
	get_api := r.Methods("GET").Subrouter()
	pst_api := r.Methods("POST").Subrouter()
	put_api := r.Methods("PUT").Subrouter()
	del_api := r.Methods("DELETE").Subrouter()

	// Define routes
	get_api.HandleFunc("/runner", GetRunnerHandler)
	get_api.HandleFunc("/runner/{id:[0-9]+}", GetRunnerByIdHandler)
	del_api.HandleFunc("/runner/{id:[0-9]+}", DeleteRunnerHandler)
	pst_api.HandleFunc("/runner/{id:[0-9]+/enable}", PostEnableRunnerHandler)
	pst_api.HandleFunc("/runner/{id:[0-9]+/disable}", PostDisableRunnerHandler)

	get_api.HandleFunc("/run", GetRunHandler)
	get_api.HandleFunc("/run/active", GetActiveRunHandler)
	get_api.HandleFunc("/run/{id:[0-9]+}", GetRunByIdHandler)
	put_api.HandleFunc("/run", PutRunHandler)
	del_api.HandleFunc("/run/{id:[0-9]+}", DeleteRunHandler)

	get_api.HandleFunc("/batch", GetBatchsHandler)
	get_api.HandleFunc("/batch/{id:[0-9]+}", GetBatchByIdHandler)
	del_api.HandleFunc("/batch/{id:[0-9]+}", DeleteBatchHandler)

	get_api.HandleFunc("/namespace", GetNamespaceHandler)
	get_api.HandleFunc("/namespace/{name:[--~]+}", GetNamespaceByNameHandler)
	put_api.HandleFunc("/namespace", PutNamespaceHandler)
	del_api.HandleFunc("/namespace/{name:[--~]+}", DeleteNamespaceHandler)

	pst_api.HandleFunc("/ci/{namespace:[--~]+}", PostCommitHandler)

	// Launch async server with router
	var err error
	(*wait).Add(1)
	go func() {
		http.Handle("/", r)
		log.Print("INFO: Launching http server with crt: '" + crt_file + "'")
		log.Print("      and key: '" + key_file + "'")
		err = http.ListenAndServeTLS(":"+port, crt_file, key_file, nil)
		defer (*wait).Done()
	}()
	time.Sleep(1 * time.Second)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
