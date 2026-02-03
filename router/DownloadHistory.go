package router

import (
	"encoding/json"
	"errors"
	"net/http"
	"videodowload/model"
	"videodowload/storage"
)

type HistoryService struct {
	Store *storage.Store
}

var historyService HistoryService

func NewHistoryService(store *storage.Store) {
	historyService.Store = store
}

// @Summary 获取下载历史
// @Description 获取所有视频下载的历史记录列表
// @Tags 历史记录
// @Accept json
// @Produce json
// @Success 200 {array} model.DownLoadHis "下载历史记录列表"
// @Failure 500 {object} string "服务器内部错误"
// @Router /history [get]
func DownloadHistory(w http.ResponseWriter, r *http.Request) {
	var list []model.DownLoadHis
	list, err := historyService.Store.FindData()
	if err != nil {
		//因为在查询数据库阶段就失败的话就会给用户提示,所以说在网络发送阶段就无需给予提示
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if list == nil {
		w.Write([]byte("[]"))
		return
	}
	j, err := json.Marshal(list)
	if err != nil {
		http.Error(w, errors.New("获取历史下载记录失败").Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(j); err != nil {
		http.Error(w, "写入响应失败", http.StatusInternalServerError)
	}
}
