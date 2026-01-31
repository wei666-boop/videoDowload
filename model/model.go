package model

import "time"

type Config struct {
	Url       string `json:"url"`
	Type      string `json:"type"`
	Subtitle  string `json:"subtitle"`
	Thumbnail string `json:"thumbnail"`
}

type DownLoadHis struct {
	Id       int       `json:"id"`
	URL      string    `json:"uri"`
	CreateAt time.Time `json:"time"`
}
