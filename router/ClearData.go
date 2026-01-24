package router

import (
	"net/http"
	"videodowload/storage"
)

func ClearData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := storage.ClearData()
	if err != nil {
		http.Error(w, "not clear the database", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database cleared successfully"))
}
