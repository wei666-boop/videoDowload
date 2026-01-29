package utils

import (
	"net/http"
	"runtime"
	"sync"
	"time"
)

var once sync.Once

var RealUA string

var Client = &http.Client{
	Timeout: 10 * time.Second,
}

func NewUA() {
	switch runtime.GOOS {
	case "windows":
		RealUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)\nAppleWebKit/537.36 (KHTML, like Gecko)\nChrome/120.0.0.0 Safari/537.36\n"
	case "linux":
		RealUA = "Mozilla/5.0 (X11; Linux x86_64)\nAppleWebKit/537.36 (KHTML, like Gecko)\nChrome/120.0.0.0 Safari/537.36\n"
	case "darwin":
		RealUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)\nAppleWebKit/537.36 (KHTML, like Gecko)\nChrome/120.0.0.0 Safari/537.36\n"
	default:
		RealUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)\nAppleWebKit/537.36 (KHTML, like Gecko)\nChrome/120.0.0.0 Safari/537.36\n"
	}

}

func WarmUp(client *http.Client, url string) {
	once.Do(NewUA)
	client.Get("https://www.bilibili.com")
	time.Sleep(1 * time.Second)

	// 再访问视频页面

	//先访问主页->建立会话->在访问视频
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", RealUA)
	req.Header.Set("Referer", "https://www.bilibili.com") //模拟从b站主页访问视频
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, _ := client.Do(req)
	if resp != nil {
		resp.Body.Close()
	}
}
