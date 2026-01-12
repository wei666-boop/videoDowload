package router

import (
	"encoding/json"
	"errors"
	"net/http"
	"videodowload/model"
)

func DownloadHistory(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(model.DownloadList)
	if err != nil {
		http.Error(w, errors.New("获取历史下载记录失败").Error(), http.StatusUnauthorized)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
