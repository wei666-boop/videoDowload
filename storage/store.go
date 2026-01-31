package storage

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	SerLog "videodowload/log"
	"videodowload/model"
	"videodowload/utils"
)

type Store struct {
	DataBase *sql.DB
}

type StorePath struct {
	Path utils.Path
}

var storePath StorePath

func NewStorePath(path utils.Path) {
	storePath.Path = path
}

func (s *Store) FindData() ([]model.DownLoadHis, error) {
	rows, err := s.DataBase.Query(`SELECT Id, URL, CreateAt FROM HISTORY ORDER BY createAt DESC`)
	if err != nil {
		SerLog.WriteLog(3, "无法查询历史记录: "+err.Error(), SerLog.GetLog(storePath.Path.LogPath.Store))
		return nil, err
	}
	defer rows.Close()

	var list []model.DownLoadHis
	for rows.Next() {
		var his model.DownLoadHis
		if err := rows.Scan(&his.Id, &his.URL, &his.CreateAt); err != nil {
			SerLog.WriteLog(3, "扫描历史记录失败: "+err.Error(), SerLog.GetLog(storePath.Path.LogPath.Store))
			return nil, err
		}
		list = append(list, his)
	}

	if err = rows.Err(); err != nil {
		SerLog.WriteLog(3, "遍历行时出错: "+err.Error(), SerLog.GetLog(storePath.Path.LogPath.Store))
		return nil, err
	}
	return list, nil
}

func (s *Store) InitTable() error {
	_, err := s.DataBase.Exec(`CREATE TABLE IF NOT EXISTS HISTORY (
        Id INTEGER PRIMARY KEY AUTOINCREMENT,
        URL TEXT NOT NULL,
        CreateAt DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CloseDB() {
	if s.DataBase != nil {
		s.DataBase.Close()
	}

}

func (s *Store) InsertData(uri string) error {
	if s == nil || s.DataBase == nil {
		return errors.New("store or database is nil")
	}
	_, err := s.DataBase.Exec(`INSERT INTO HISTORY(URL) VALUES (?)`, uri)
	if err != nil {
		SerLog.WriteLog(3, "无法新增数据"+err.Error(), SerLog.GetLog(storePath.Path.LogPath.Store))
		return err
	}
	return nil
}

func (s *Store) ClearData() error {
	// 开始事务

	//这里使用事务，因为删除计数器和删除表的数据要求同时成功

	tx, err := s.DataBase.Begin()
	if err != nil {
		SerLog.WriteLog(3, "开始事务失败"+err.Error(), SerLog.GetLog(storePath.Path.LogPath.Store))
		return err
	}

	// 删除所有数据
	_, err = tx.Exec(`DELETE FROM HISTORY`)
	if err != nil {
		tx.Rollback()
		SerLog.WriteLog(3, "清空数据失败"+err.Error(), SerLog.GetLog(storePath.Path.LogPath.Store))
		return err
	}

	// 重置AUTOINCREMENT计数器
	_, err = tx.Exec(`DELETE FROM sqlite_sequence WHERE name='HISTORY'`)
	if err != nil {
		tx.Rollback()
		SerLog.WriteLog(3, "重置ID计数器失败"+err.Error(), SerLog.GetLog(storePath.Path.LogPath.Store))
		return err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		SerLog.WriteLog(3, "提交事务失败"+err.Error(), SerLog.GetLog(storePath.Path.LogPath.Store))
		return err
	}

	return nil
}
