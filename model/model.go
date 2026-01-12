package model

import "time"

type Config struct {
	Url       string `json:"url"`
	Type      string `json:"type"`
	Subtitle  string `json:"subtitle"`
	Thumbnail string `json:"thumbnail"`
}

type DownLoadHis struct {
	URL  string
	Time time.Time
}

var DownloadList []DownLoadHis
