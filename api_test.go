package qzone

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	// cookie 登录成功后自动更新
	cookie = ""
)

func TestLogin(t *testing.T) {
	var (
		qrsig           string
		ptqrtoken       string
		ptqrloginCookie string
		redirectCookie  string
		data            []byte
		err             error
	)
	data, qrsig, ptqrtoken, err = Ptqrshow()
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("ptqrcode.png", data, 0666)
	if err != nil {
		t.Fatal(err)
	}
LOOP:
	for {
		time.Sleep(2 * time.Second)
		data, ptqrloginCookie, err = Ptqrlogin(qrsig, ptqrtoken)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("ptqrloginCookie:%v\n", ptqrloginCookie)
		text := string(data)
		t.Logf("text:%v\n", text)
		switch {
		case strings.Contains(text, "二维码已失效"):
			t.Fatal("二维码已失效, 登录失败")
			return
		case strings.Contains(text, "登录成功"):
			_ = os.Remove("ptqrcode.png")
			dealedCheckText := strings.ReplaceAll(text, "'", "")
			redirectURL := strings.Split(dealedCheckText, ",")[2]
			redirectCookie, err = LoginRedirect(redirectURL)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("ptqrloginCookie:%v\n", redirectCookie)
			break LOOP
		}
	}
	m := NewManager(ptqrloginCookie + redirectCookie)

	path := `D:\1.png`
	srcByte, err := os.ReadFile(path)
	if err != nil {
		return
	}
	picBase64 := base64.StdEncoding.EncodeToString(srcByte)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("m:%#v\n", m)
	err = os.WriteFile("cookie.txt", []byte(ptqrloginCookie+redirectCookie), 0666)
	if err != nil {
		t.Fatal(err)
	}

	//读取cookie
	cookie = ptqrloginCookie + redirectCookie
	t.Logf("ptqrloginCookie:%v\n", ptqrloginCookie+redirectCookie)

	gotResult, err := m.EmotionPublish("真好", []string{picBase64})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("gotResult:%#v\n", gotResult)
}

func TestManager_PublishEmotion(t *testing.T) {
	type args struct {
		Content string
	}
	m := NewManager(cookie)
	gotResult, err := m.EmotionPublish("test", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("gotResult:%#v\n", gotResult)
}

func TestManager_UploadImage(t *testing.T) {
	m := NewManager(cookie)
	path := `D:\1.png`
	srcByte, err := os.ReadFile(path)
	if err != nil {
		return
	}
	picBase64 := base64.StdEncoding.EncodeToString(srcByte)
	if err != nil {
		t.Fatal(err)
	}
	gotResult, err := m.UploadImage(picBase64)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("gotResult:%#v\n", gotResult)
}

func TestManager_Msglist(t *testing.T) {
	m := NewManager(cookie)
	gotResult, err := m.EmotionMsglist("1", "1")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("gotResult:%#v\n", gotResult)
}
