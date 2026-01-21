package router

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	SerLog "videodowload/log"
	"videodowload/model"
	"videodowload/storage"
	"videodowload/utils"
)

func parseConfig(w http.ResponseWriter, r *http.Request) (model.Config, error) {
	var configStruct model.Config
	//将发送过来的json数据映射到Config结构体中
	err := json.NewDecoder(r.Body).Decode(&configStruct)
	if err != nil {
		return model.Config{}, errors.New("数据解析失败")
	}
	if configStruct.Url == "" {
		return model.Config{}, errors.New("没有该资源")
	}
	url := configStruct.Url
	decodeUrl, err := base64.StdEncoding.DecodeString(url)
	if err != nil {
		return model.Config{}, errors.New("解码失败")
	}
	configStruct.Url = string(decodeUrl)
	return configStruct, nil
}

//ToDo重构这段代码

func Download(w http.ResponseWriter, r *http.Request) {
	SerLog.WriteLog(1, r.Body, SerLog.GetLog(model.GlobalPath.LogPath.Service))
	var configStruct model.Config

	configStruct, err := parseConfig(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	decodeUrl := configStruct.Url
	//将现在记录存储进数据库
	if err = storage.InsertData(decodeUrl); err != nil {
		panic("添加数据失败" + err.Error())
	}
	SerLog.WriteLog(1, configStruct.Type, SerLog.GetLog(model.GlobalPath.LogPath.Service))
	//生成随机工作目录
	dir, err := utils.RandomID(model.GlobalPath.TempPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(dir)

	var (
		VIDEOFILE = filepath.Join(dir, "./video.mp3")
		AUDIOFILE = filepath.Join(dir, "audio.mp4")
		THUMBNAIL = filepath.Join(dir, "./video")
		SUBTITLE  = filepath.Join(dir, "./video.srt")
		MKVFILE   = filepath.Join(dir, "./output.mkv")
	)

	var (
		subtitlePath  = ""
		thumbnailPath = ""
	)

	var args []string

	//主要功能为不同的配置执行不同功能

	//ToDo重构这个写的太复杂了

	switch configStruct.Type {
	case "audio":
		args = utils.Audio(AUDIOFILE, decodeUrl)
		cmd := exec.Command("yt-dlp", args...)
		err = utils.AudioAndVideoStart(cmd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		faudio, err := os.Open(filepath.Join(dir, "./video.mp3"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Disposition", "attachment;filename=\"audio.mp4\"")
		faudio.Seek(0, 0)
		io.Copy(w, faudio)
		return
	case "video":
		if configStruct.Subtitle == "true" {
			if configStruct.Thumbnail == "false" {
				args = utils.Video(VIDEOFILE, decodeUrl)
				cmd1 := exec.Command("yt-dlp", args...)
				args = utils.Subtitle(SUBTITLE, decodeUrl)
				cmd2 := exec.Command("yt-dlp", args...)
				err = utils.AudioAndVideoStart(cmd1)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				err = utils.ThumbnailORSubtitleStart(cmd2)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				utils.Video(VIDEOFILE, decodeUrl)
				cmd1 := exec.Command("yt-dlp", args...)
				args = utils.Subtitle(SUBTITLE, decodeUrl)
				cmd2 := exec.Command("yt-dlp", args...)
				args = utils.Thumbnail(THUMBNAIL, decodeUrl)
				cmd3 := exec.Command("yt-dlp", args...)
				err = utils.AudioAndVideoStart(cmd1) // yt-dlp --write-subs --write-auto-subs --convert-subs srt url
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				err = utils.ThumbnailORSubtitleStart(cmd2)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				err = utils.ThumbnailORSubtitleStart(cmd3)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		} else {
			if configStruct.Thumbnail == "true" {
				args = utils.Video(VIDEOFILE, decodeUrl)
				cmd1 := exec.Command("yt-dlp", args...)
				args = utils.Subtitle(SUBTITLE, decodeUrl)
				cmd2 := exec.Command("yt-dlp", args...)
				err = utils.AudioAndVideoStart(cmd1)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				err = utils.ThumbnailORSubtitleStart(cmd2)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				args = utils.Video(VIDEOFILE, decodeUrl)
				cmd := exec.Command("yt-dlp", args...)
				err = utils.AudioAndVideoStart(cmd)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				fvideo, err := os.Open(filepath.Join(dir, "./video.mp4"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Disposition", "attachment;filename=\"video.mp4\"")
				w.Header().Set("Content-Type", "video/mp4")
				fvideo.Seek(0, 0)
				if _, err = io.Copy(w, fvideo); err != nil {
					http.Error(w, "下载失败", http.StatusBadRequest)
					return
				}
				return
			}
		}

	default:
		http.Error(w, "invalid format", http.StatusBadRequest)
		return
	}

	//检查文件完整性(因为有一些视频受到由于平台的原因可能未提供完整资源)
	//防止拓展名不一样
	files, _ := os.ReadDir(dir)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "video.") &&
			(strings.HasSuffix(f.Name(), ".png") ||
				strings.HasSuffix(f.Name(), ".jpg")) {
			thumbnailPath = filepath.Join(dir, f.Name())
			break
		}
	}

	if _, err = os.Stat(SUBTITLE); err == nil {
		subtitlePath = SUBTITLE
	}

	//处理视频以及附属文件
	//TODO问题似乎出现在这里
	if thumbnailPath == "" && subtitlePath == "" {
		http.Error(w, errors.New("there is a unknown error").Error(), http.StatusBadRequest)
		return
	}

	err = utils.GetMKV(VIDEOFILE, subtitlePath, thumbnailPath, MKVFILE)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//设置请求头
	w.Header().Set("Content-Type", "video/x-matroska")
	w.Header().Set("Content-Disposition", "attachment;filename=\"output.mkv\"")
	out, _ := os.Open(filepath.Join(dir, "./output.mkv"))
	if _, err = io.Copy(w, out); err != nil {
		http.Error(w, "下载失败", http.StatusBadRequest)
		utils.HandleTmp(dir, w)
		return
	}

	//处理临时文件
	os.Remove(dir + "./output")
	os.RemoveAll(dir)

	SerLog.WriteLog(1, "下载完成", SerLog.GetLog(model.GlobalPath.LogPath.Service))
}

/*
如何重构这个Download函数
拆分功能
这个函数应该是一个下载的函数，不应该直接在里面写一些不属于它范畴的事情
列举这个函数目前干了什么:
1.获取前端传来的请求体并解析
2.根据请求信息拼接不同的命令
3.执行命令
4.后续处理
将这些功能拆散到不同函数中去，尤其是根据请求信息拼接命令
信息拼接这一块嵌套了复杂的if-else并且里面命令拼接和执行也没有抽象为函数
*/
