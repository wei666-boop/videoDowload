package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"videodowload/App"
	SerLog "videodowload/log"
	"videodowload/storage"
	"videodowload/utils"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		panic("读取配置失败")
	}
}

func InitPath() utils.Path {
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
	storeLogPath := utils.ResolvePath(configRoot, viper.GetString("backend.store_log_path"))
	serviceLogPath := utils.ResolvePath(configRoot, viper.GetString("backend.service_log_path"))
	tempPath := os.TempDir()
	fmt.Println(serviceLogPath)
	path := utils.Path{
		TempPath: tempPath,
		LogPath: utils.LogFilePath{
			Service: serviceLogPath,
			Store:   storeLogPath,
		},
	}
	return path
}

func main() {
	InitConfig()
	path := InitPath()
	if &path == nil {
		panic("初始化路径失败")
	}
	//加载资源
	db, _ := utils.NewDB()
	if db == nil {
		SerLog.WriteLog(3, "数据库指针为空", SerLog.GetLog(path.LogPath.Store))
		panic("数据库指针为空")
	}
	defer db.Close()
	SerLog.WriteLog(1, "数据库初始化成功", SerLog.GetLog(path.LogPath.Store))
	store := storage.Store{DataBase: db}
	err := store.InitTable()
	if err != nil {
		SerLog.WriteLog(3, "初始化表"+err.Error(), SerLog.GetLog(path.LogPath.Store))
	}
	a := App.NewApp(&store, path)
	fmt.Println("已开始监听")
	a.Run(viper.GetString("service.port"))
}

//测试使用地址:https://www.bilibili.com/video/BV157GHzeEV3?t=68.7
