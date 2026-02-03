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

// @Summary 清空数据库
// @Description 删除数据库中的所有下载历史记录
// @Tags 数据管理
// @Accept json
// @Produce json
// @Success 200 {string} string "Database cleared successfully"
// @Failure 400 {object} string "清空数据库失败"
// @Failure 405 {object} string "请求方法不允许"
// @Router /clear [delete]
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
