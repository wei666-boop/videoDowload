package router

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"videodowload/model"
	"videodowload/router"
)

// 测试传入有效值的时候能否正确解析
func TestParseConfig_ValidInput(t *testing.T) {
	config := model.Config{
		Url: "aHR0cHM6Ly93d3cuYmlsaWJpbGkuY29tL3ZpZGVvL0JWMTU3R0h6ZUVWMz90PTY4Ljc=",
	}
	body, _ := json.Marshal(config)
	req := httptest.NewRequest("POST", "/dl/api", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	result, err := router.ParseConfig(w, req)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if result.Url != "https://www.bilibili.com/video/BV157GHzeEV3?t=68.7" {
		t.Errorf("you don't got correct url:%v", result.Url)
	}
}

// 传入空url时候的情况
func TestParConfig_EmptyInput(t *testing.T) {
	config := model.Config{
		Url: "",
	}
	body, _ := json.Marshal(config)
	req := httptest.NewRequest("POST", "/dl/api", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	_, err := router.ParseConfig(w, req)
	if err == nil {
		t.Errorf("you should get a error,but don't get")
	}
}
