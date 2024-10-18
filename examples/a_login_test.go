package examples

import (
	"encoding/json"
	"github.com/HHU-47133/qzone"
	"os"
	"testing"
)

type Config struct {
	Cookie   string   `json:"cookie"`
	Tid      string   `json:"tid"`
	FriendQQ int64    `json:"friendQQ"`
	GroupQQ  int64    `json:"groupQQ"`
	ImgPath  []string `json:"imgPath"`
}

var Cfg Config

// 登录测试
func TestLogin(t *testing.T) {
	// 读取测试配置文件
	data, err := os.ReadFile("config.json")
	if err != nil {
		t.Fatal("读取json配置失败:", err)
	}
	// 调用登录接口
	m, err := qzone.QzoneLogin("qrcode.png", nil, 2)
	if err != nil {
		t.Fatal(err)
		t.Skip("[登录失败]请重新开启测试")
	}
	err = json.Unmarshal(data, &Cfg)
	if err != nil {
		t.Fatal(err)
	}
	Cfg.Cookie = m.Cookie
}
