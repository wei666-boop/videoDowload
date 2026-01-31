package App

import (
	"net/http"
	"videodowload/router"
	"videodowload/storage"
	"videodowload/utils"
)

type App struct {
	Mux   *http.ServeMux
	Store *storage.Store
	Path  utils.Path
}

func NewApp(store *storage.Store, path utils.Path) *App {
	mux := http.NewServeMux()

	return &App{
		Mux:   mux,
		Store: store,
		Path:  path,
	}
}

func (a *App) Run(port string) {
	storage.NewStorePath(a.Path)
	router.NewDownloadService(a.Store, a.Path)
	router.NewHistoryService(a.Store)
	router.NewClearService(a.Store)
	router.RouterRegister(a.Mux)
	http.ListenAndServe(":"+port, a.Mux)
}
