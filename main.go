package main

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"videodowload/router"
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

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return
	}

}

func main() {

	//加载资源

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

	http.HandleFunc("/center/record", func(w http.ResponseWriter, r *http.Request) {
		router.DownloadHistory(w, r)
	})

	http.ListenAndServe(":"+viper.GetString("service.port"), nil)
}

//测试使用地址:https://www.bilibili.com/video/BV157GHzeEV3?t=68.7
