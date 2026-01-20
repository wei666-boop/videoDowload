package main

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path/filepath"
	SerLog "videodowload/log"
	"videodowload/model"
	"videodowload/router"
	"videodowload/storage"
	"videodowload/utils"
)

// 添加CORS中间件
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		//允许所有跨域请求 支持常用方法 允许自定义请求头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 调用下一个处理函数
		next.ServeHTTP(w, r)
	}
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		panic("读取配置失败")
	}
}

func main() {
	config, err := os.UserConfigDir()
	if err != nil {
		panic(err.Error())
	}
	temp := os.TempDir()
	fmt.Println(temp)
	//配置应用的根路径
	configRoot := filepath.Join(config, "videodowload")
	tempRoot := filepath.Join(temp, "videodowload")
	//创建这个目录
	err = os.MkdirAll(configRoot, 0755)
	if err != nil {
		panic("无法创建目录" + err.Error())
	}
	err = os.MkdirAll(tempRoot, 0755)
	if err != nil {
		panic("无法创建目录" + err.Error())
	}
	InitConfig()
	storeLogPath := utils.ResolvePath(configRoot, viper.GetString("backend.store_log_path"))
	serviceLogPath := utils.ResolvePath(configRoot, viper.GetString("backend.service_log_path"))
	tempPath := os.TempDir()
	fmt.Println(serviceLogPath)
	model.GlobalPath.LogPath.Store = storeLogPath
	model.GlobalPath.LogPath.Service = serviceLogPath
	model.GlobalPath.TempPath = tempPath

	//加载资源
	storage.GetDB()
	storage.InitTable()
	SerLog.WriteLog(1, "数据库初始化成功", SerLog.GetLog(model.GlobalPath.LogPath.Store))
	defer storage.CloseDB()

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

	http.HandleFunc("/dl/api", enableCORS(router.Download))

	http.HandleFunc("/center", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/center.html")
	})

	http.HandleFunc("/center/record", enableCORS(router.DownloadHistory))

	http.HandleFunc("/center/clear", enableCORS(router.ClearData))
	fmt.Println("已开始监听")
	http.ListenAndServe(":"+viper.GetString("service.port"), nil)
}

//测试使用地址:https://www.bilibili.com/video/BV157GHzeEV3?t=68.7
