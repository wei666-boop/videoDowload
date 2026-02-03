package router

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "videodowload/docs"
	"videodowload/middle"
)

func RouterRegister(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	mux.HandleFunc("/setting", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/setting.html")
	})

	mux.HandleFunc("/dl", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/page.html")
	})

	mux.HandleFunc("/dl/api", middle.EnableCORS(Download))

	mux.HandleFunc("/center", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/center.html")
	})

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/center/record", middle.EnableCORS(DownloadHistory))

	mux.HandleFunc("/center/clear", middle.EnableCORS(ClearData))
}
