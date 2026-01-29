package router

import (
	"net/http"
	"videodowload/middle"
)

func RouterRegister() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	http.HandleFunc("/setting", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/setting.html")
	})

	http.HandleFunc("/dl", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/page.html")
	})

	http.HandleFunc("/dl/api", middle.EnableCORS(Download))

	http.HandleFunc("/center", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/center.html")
	})

	http.HandleFunc("/center/record", middle.EnableCORS(DownloadHistory))

	http.HandleFunc("/center/clear", middle.EnableCORS(ClearData))
}
