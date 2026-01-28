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

type operation func(w http.ResponseWriter, dir string, decodeUrl string, params model.Param)

type key struct {
	Type      string
	Subtitle  string
	Thumbnail string
}

func ParseConfig(w http.ResponseWriter, r *http.Request) (model.Config, error) {
	var configStruct model.Config
	//将发送过来的json数据映射到Config结构体中
	fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&configStruct)

	if err != nil {
		return model.Config{}, errors.New("数据解析失败" + err.Error())
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

func downloadAudioOnly(w http.ResponseWriter, dir string, decodeUrl string, param model.Param) {
	utils.WarmUp(utils.Client, decodeUrl)
	var args []string
	args = utils.Audio(param.AudioFile, decodeUrl)
	cmd := exec.Command("yt-dlp", args...)
	err := utils.AudioAndVideoStart(cmd)
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
}

func downloadVideoOnly(w http.ResponseWriter, dir string, decodeUrl string, param model.Param) {
	utils.WarmUp(utils.Client, decodeUrl)
	var args []string
	args = utils.Video(param.VideoFile, decodeUrl)
	cmd := exec.Command("yt-dlp", args...)
	err := utils.AudioAndVideoStart(cmd)
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
}

func downloadVideoWithSubtitle(w http.ResponseWriter, dir string, decodeUrl string, param model.Param) {
	utils.WarmUp(utils.Client, decodeUrl)
	var args []string
	args = utils.Video(param.VideoFile, decodeUrl)
	cmd1 := exec.Command("yt-dlp", args...)
	args = utils.Subtitle(param.Subtitle, decodeUrl)
	cmd2 := exec.Command("yt-dlp", args...)
	err := utils.AudioAndVideoStart(cmd1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = utils.ThumbnailORSubtitleStart(cmd2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func downloadVideoComplete(w http.ResponseWriter, dir string, decodeUrl string, param model.Param) {
	utils.WarmUp(utils.Client, decodeUrl)
	var args []string
	args = utils.Video(param.VideoFile, decodeUrl)
	cmd1 := exec.Command("yt-dlp", args...)
	args = utils.Subtitle(param.Subtitle, decodeUrl)
	cmd2 := exec.Command("yt-dlp", args...)
	args = utils.Thumbnail(param.Thumbnail, decodeUrl)
	cmd3 := exec.Command("yt-dlp", args...)
	err := utils.AudioAndVideoStart(cmd1) // yt-dlp --write-subs --write-auto-subs --convert-subs srt url
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

func downloadVideoWithThumbnail(w http.ResponseWriter, dir string, decodeUrl string, param model.Param) {
	utils.WarmUp(utils.Client, decodeUrl)
	var args []string
	args = utils.Video(param.VideoFile, decodeUrl)
	cmd1 := exec.Command("yt-dlp", args...)
	args = utils.Thumbnail(param.Thumbnail, decodeUrl)
	cmd2 := exec.Command("yt-dlp", args...)
	err := utils.AudioAndVideoStart(cmd1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = utils.ThumbnailORSubtitleStart(cmd2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

var operations = make(map[key]operation, 5)

func init() {
	operations[key{
		Type:      "audio",
		Subtitle:  "false",
		Thumbnail: "false",
	}] = downloadAudioOnly
	operations[key{
		Type:      "video",
		Subtitle:  "false",
		Thumbnail: "false",
	}] = downloadVideoOnly
	operations[key{
		Type:      "video",
		Subtitle:  "false",
		Thumbnail: "true",
	}] = downloadVideoWithThumbnail
	operations[key{
		Type:      "video",
		Subtitle:  "true",
		Thumbnail: "false",
	}] = downloadVideoWithSubtitle
	operations[key{
		Type:      "video",
		Subtitle:  "true",
		Thumbnail: "true",
	}] = downloadVideoComplete
}

func Download(w http.ResponseWriter, r *http.Request) {
	SerLog.WriteLog(1, r.Body, SerLog.GetLog(model.GlobalPath.LogPath.Service))
	var configStruct model.Config

	configStruct, err := ParseConfig(w, r)
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

	var params model.Param

	params.VideoFile = filepath.Join(dir, "./video.mp4")

	params.Thumbnail = filepath.Join(dir, "./video")
	params.Subtitle = filepath.Join(dir, "./video.srt")
	params.MkvFile = filepath.Join(dir, "./output.mkv")
	fmt.Println(params.VideoFile)
	var (
		subtitlePath  = ""
		thumbnailPath = ""
	)

	//主要功能为不同的配置执行不同功能

	Key := key{
		Type:      configStruct.Type,
		Subtitle:  configStruct.Subtitle,
		Thumbnail: configStruct.Thumbnail,
	}
	download, ok := operations[Key]
	if !ok {
		panic("未找到相应的函数")
	}
	download(w, dir, decodeUrl, params)

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

	if _, err = os.Stat(params.Subtitle); err == nil {
		subtitlePath = params.Subtitle
	}

	//处理视频以及附属文件
	if thumbnailPath == "" && subtitlePath == "" {
		http.Error(w, errors.New("there is a unknown error").Error(), http.StatusBadRequest)
		return
	}

	err = utils.GetMKV(params.VideoFile, subtitlePath, thumbnailPath, params.MkvFile)

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
