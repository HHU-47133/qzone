package examples

import (
	"fmt"
	"github.com/HHU-47133/qzone"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	QrcodeName = "ptqrcode.png"
)

// 登录测试
func TestLogin(t *testing.T) {
	var m qzone.Manager
	// 1. 获取二维码信息（data），取出cookie重要参数（qrsig、ptqrtoken）
	data, qrsig, ptqrtoken, err := qzone.Ptqrshow()
	if err != nil {
		t.Fatal(err)
	}
	// 2. 保存二维码
	err = os.WriteFile(QrcodeName, data, 0666)
	if err != nil {
		t.Fatal(err)
	}
	// 3. 查询登录回调，检测登录状态
LOOP:
	for {
		data, ptqrloginCookie, err := qzone.Ptqrlogin(qrsig, ptqrtoken)
		if err != nil {
			t.Fatal(err)
		}
		text := string(data)
		fmt.Printf("%#v\n", text)
		switch {
		case strings.Contains(text, "二维码已失效"):
			t.Fatal("二维码已失效, 登录失败")
			return
		case strings.Contains(text, "登录成功"):
			_ = os.Remove(QrcodeName)
			dealedCheckText := strings.ReplaceAll(text, "'", "")
			redirectURL := strings.Split(dealedCheckText, ",")[2]
			// 4. 成功登录后，获取登录重定向URL
			redirectCookie, err := qzone.LoginRedirect(redirectURL)
			if err != nil {
				t.Fatal(err)
			}
			// 5. 创建信息管理结构，携带登录回调cookie和重定向页面cookie
			m = qzone.NewManager(ptqrloginCookie + redirectCookie)
			cookie = ptqrloginCookie + redirectCookie
			t.Log("cookie: ", cookie)
			break LOOP
		}
		time.Sleep(2 * time.Second)
	}

	// 6. 执行其它接口操作
	fmt.Println(m.Uin, m.QQ, m.Gtk2, m.Cookie)
}
