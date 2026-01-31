package router

import (
	"net/http"
	"videodowload/storage"
)

type ClearService struct {
	store *storage.Store
}

var clearService ClearService

func NewClearService(store *storage.Store) {
	clearService.store = store
}

func ClearData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := clearService.store.ClearData()
	if err != nil {
		http.Error(w, "not clear the database", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database cleared successfully"))
}
