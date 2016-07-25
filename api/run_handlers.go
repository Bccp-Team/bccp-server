package api

import "net/http"

// List all run
func GetRunHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Get information about given run
func GetRunByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Add run
func PutRunHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Delete given runner
func DeleteRunHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Update run logs
func PostRunUpdateByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}
