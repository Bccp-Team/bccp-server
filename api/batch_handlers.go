package api

import "net/http"

// Get information about given batch
func GetBatchByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Delete given batch
func DeleteBatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}
