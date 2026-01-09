package utils

import (
	"net/http"
	"os"
	"path/filepath"
)

func HandleTmp(dir string, w http.ResponseWriter) {
	files, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, "not exist temp", http.StatusBadRequest)
		return
	}

	for _, file := range files {
		if file.Name() == "output.mkv" {
			continue
		} else {
			os.Remove(filepath.Join(dir, file.Name()))
		}
	}

	//最后删除临时目录
	//os.RemoveAll(dir)
}
