package router

import (
	"net/http"
	"videodowload/storage"
)

func ClearData(w http.ResponseWriter, r *http.Request) {
	err := storage.ClearData()
	if err != nil {
		http.Error(w, "not clear the database", http.StatusBadRequest)
	}
}
