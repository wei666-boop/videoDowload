package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"videodowload/utils"
)

type Config struct {
	Url       string `json:"url"`
	Type      string `json:"type"`
	Subtitle  string `json:"subtitle"`
	Thumbnail string `json:"thumbnail"`
}

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

// 下载函数
func download(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	var configStruct Config
	//将发送过来的json数据映射到Config结构体中
	err := json.NewDecoder(r.Body).Decode(&configStruct)

	if err != nil {
		http.Error(w, "invalid config", http.StatusBadRequest)
		return
	}

	fmt.Println(configStruct)
	//检查配置结构体中的数据是否符合要求
	if configStruct.Url == "" {
		http.Error(w, "no url", http.StatusBadRequest)
		return
	}
	url := configStruct.Url
	//解码得到url
	decordUrl, err := base64.StdEncoding.DecodeString(url)
	if err != nil {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}
	fmt.Println(configStruct.Type)
	//生成随机工作目录
	dir, err := utils.RandomID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(dir)

	var (
		VIDEOFILE = dir + "\\video.mp4"
		AUDIOFILE = dir + "\\audio.mp3"
		THUMBNAIL = dir + "\\video"
		SUBTITLE  = dir + "\\video.srt"
		MKVFILE   = dir + "\\output.mkv"
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
		args = append(args, "-o", AUDIOFILE)
		args = append(args, "-x", "--audio-format", "mp3")
		args = append(args, string(decordUrl))
		cmd := exec.Command("yt-dlp", args...)
		err = utils.AudioAndVideoStart(cmd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		faudio, err := os.Open(dir + "\\video,mp3")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "audio/mpeg")
		faudio.Seek(0, 0)
		io.Copy(w, faudio)
		return
	case "video":
		if configStruct.Subtitle == "true" {
			if configStruct.Thumbnail == "false" {
				args = append(args, "-o", VIDEOFILE)
				args = append(args, string(decordUrl))
				cmd1 := exec.Command("yt-dlp", args...)
				args = nil
				args = append(args, "--skip-download")
				args = append(args, "--write-subs", "--write-auto-subs")
				args = append(args, "-o", SUBTITLE, string(decordUrl))
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
				args = append(args, "-o", VIDEOFILE)
				args = append(args, string(decordUrl))
				cmd1 := exec.Command("yt-dlp", args...)
				args = nil
				args = append(args, "--skip-download")
				args = append(args, "--write-subs", "--write-auto-subs")
				args = append(args, "-o", SUBTITLE, string(decordUrl))
				cmd2 := exec.Command("yt-dlp", args...)
				args = nil
				args = append(args, "--skip-download")
				args = append(args, "--write-thumbnail")
				args = append(args, "-o", THUMBNAIL, string(decordUrl))
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
				args = append(args, "-o", VIDEOFILE)
				args = append(args, string(decordUrl))
				cmd1 := exec.Command("yt-dlp", args...)
				args = nil
				args = append(args, "--skip-download")
				args = append(args, "--write-thumbnail")
				args = append(args, "-o", THUMBNAIL, string(decordUrl))
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
				args = append(args, "-o", VIDEOFILE)
				args = append(args, string(decordUrl))
				cmd := exec.Command("yt-dlp", args...)
				err = utils.AudioAndVideoStart(cmd)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				fvideo, err := os.Open(dir + "\\vidio.mp4")
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "video/mp4")
				w.WriteHeader(http.StatusOK)
				fvideo.Seek(0, 0)
				io.Copy(w, fvideo)
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

	fmt.Println(subtitlePath, thumbnailPath)

	err = utils.GetMKV(VIDEOFILE, subtitlePath, thumbnailPath, MKVFILE)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//设置请求头
	w.Header().Set("Content-Type", "video/x-matroska")
	w.Header().Set("Content-Disposition", "attachment;filename=\"output.mkv\"")
	//向浏览器返回请求
	//http.ServeFile(w, r, dir+"\\output.mkv")
	out, _ := os.Open(dir + "\\output.mkv")
	if _, err = io.Copy(w, out); err != nil {
		http.Error(w, "下载失败", http.StatusBadRequest)
		utils.HandleTmp(dir, w)
		return
	}

	//处理临时文件
	utils.HandleTmp(dir, w)
	fmt.Println("已删除缓存文件")

	fmt.Println("下载完成")

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

	http.HandleFunc("/dl/api", enableCORS(download))

	http.ListenAndServe(":"+viper.GetString("service.port"), nil)
}

//测试使用地址:https://www.bilibili.com/video/BV157GHzeEV3?t=68.7
