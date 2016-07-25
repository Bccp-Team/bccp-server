package api

import "net/http"

// List all run
func GetNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Get information about given run
func GetNamespaceByNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Add run
func PutNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Delete given runner
func DeleteNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}
