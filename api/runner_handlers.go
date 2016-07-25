package api

import "net/http"

// List all runners
func GetRunnerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Get information about given runner
func GetRunnerByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Delete given runner
func DeleteRunnerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Enable given runner
func PostEnableRunnerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Disable given runner
func PostDisableRunnerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// Add runner
func PutRunnerRegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}
