package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sync"
	SerLog "videodowload/log"
	"videodowload/model"
)

var once sync.Once
var DB *sql.DB

func NewDB() {
	db, err := sql.Open("sqlite3", "./history.db")
	if err != nil {
		SerLog.WriteLog(3, "无法打开初始化数据库"+err.Error(), SerLog.GetLog(model.GlobalPath.LogPath.Store))
		panic("无法打开数据库" + err.Error())
	}
	err = db.Ping()
	if err != nil {
		SerLog.WriteLog(3, "无法打开链接数据库"+err.Error(), SerLog.GetLog(model.GlobalPath.LogPath.Store))
		panic("无法连接数据库" + err.Error())
	}
	DB = db
}

func GetDB() *sql.DB {
	once.Do(NewDB)
	return DB
}

func FindData() ([]model.DownLoadHis, error) {
	rows, err := GetDB().Query(`SELECT Id, URL, CreateAt FROM HISTORY ORDER BY createAt DESC`)
	if err != nil {
		SerLog.WriteLog(3, "无法查询历史记录: "+err.Error(), SerLog.GetLog(model.GlobalPath.LogPath.Store))
		return nil, err
	}
	defer rows.Close()

	var list []model.DownLoadHis
	for rows.Next() {
		var his model.DownLoadHis
		if err := rows.Scan(&his.Id, &his.URL, &his.CreateAt); err != nil {
			SerLog.WriteLog(3, "扫描历史记录失败: "+err.Error(), SerLog.GetLog(model.GlobalPath.LogPath.Store))
			return nil, err
		}
		list = append(list, his)
	}

	if err = rows.Err(); err != nil {
		SerLog.WriteLog(3, "遍历行时出错: "+err.Error(), SerLog.GetLog(model.GlobalPath.LogPath.Store))
		return nil, err
	}
	return list, nil
}

func InitTable() {
	_, err := GetDB().Exec(`CREATE TABLE IF NOT EXISTS HISTORY (
        Id INTEGER PRIMARY KEY AUTOINCREMENT,
        URL TEXT NOT NULL,
        CreateAt DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)
	if err != nil {
		SerLog.WriteLog(3, "初始化表"+err.Error(), SerLog.GetLog(model.GlobalPath.LogPath.Store))
		panic("初始化表失败" + err.Error())
	}
}

func CloseDB() {
	if GetDB() != nil {
		GetDB().Close()
	}

}

func InsertData(uri string) error {
	_, err := GetDB().Exec(`INSERT INTO HISTORY(URL) VALUES (?)`, uri)
	if err != nil {
		SerLog.WriteLog(3, "无法新增数据"+err.Error(), SerLog.GetLog(model.GlobalPath.LogPath.Store))
		return err
	}
	return nil
}

func ClearData() error {
	// 开始事务

	//这里使用事务，因为删除计数器和删除表的数据要求同时成功

	tx, err := GetDB().Begin()
	if err != nil {
		SerLog.WriteLog(3, "开始事务失败"+err.Error(), SerLog.GetLog(model.GlobalPath.LogPath.Store))
		return err
	}

	// 删除所有数据
	_, err = tx.Exec(`DELETE FROM HISTORY`)
	if err != nil {
		tx.Rollback()
		SerLog.WriteLog(3, "清空数据失败"+err.Error(), SerLog.GetLog(model.GlobalPath.LogPath.Store))
		return err
	}

	// 重置AUTOINCREMENT计数器
	_, err = tx.Exec(`DELETE FROM sqlite_sequence WHERE name='HISTORY'`)
	if err != nil {
		tx.Rollback()
		SerLog.WriteLog(3, "重置ID计数器失败"+err.Error(), SerLog.GetLog(model.GlobalPath.LogPath.Store))
		return err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		SerLog.WriteLog(3, "提交事务失败"+err.Error(), SerLog.GetLog(model.GlobalPath.LogPath.Store))
		return err
	}

	return nil
}
