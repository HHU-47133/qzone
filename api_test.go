package qzone

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/FloatTech/floatbox/binary"
)

var (
	// cookie 填写登录成功后的cookie
	cookie = "qrsig=69aa8b3f28b9cfe17ccdbb027ab9be6889247678d392a4afa3f5d872e08bac3a0cd649e876892cdb424ba6933a11de687feddb068d059147;uin=o1778046356;skey=@BLfHfFf0Z;pt2gguin=o1778046356;p_uin=o1778046356;pt4_token=8pYL28kcg6smOZN-NtkYVovq5VhtqUaJBh8reeTULis_;p_skey=6jlXB0HmdM9hJ83AkOirKy94IDC6w8h4hrRke9cNYP4_;"
)

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
		text := binary.BytesToString(data)
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
	err = os.WriteFile("cookie.txt", binary.StringToBytes(ptqrloginCookie+redirectCookie), 0666)
	if err != nil {
		t.Fatal(err)
	}
	gotResult, err := m.EmotionPublish("真好", []string{picBase64})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("gotResult:%#v\n", gotResult)
}
