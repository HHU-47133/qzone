package examples

import (
	"encoding/base64"
	"github.com/HHU-47133/qzone"
	"os"
	"testing"
	"time"
)

var (
	cookie   = ""
	groupID  int64
	friendID int64
	tid      string
	imgPath  = [2]string{"./1.png", "./2.png"}
)

// 登录测试
func TestLogin(t *testing.T) {
	// 创建QZone对象, 使用扫码登录
	qm := qzone.NewQZone()
	b64s, err := qm.GenerateQRCode()
	if err != nil {
		t.Fatal("扫码登录获取二维码失败:", err)
	}

	ddd, err := base64.StdEncoding.DecodeString(b64s)
	if err != nil {
		t.Fatal("扫码登录base64解码失败:", err)
	}

	err = os.WriteFile("./qrcode.png", ddd, 0666)
	if err != nil {
		t.Fatal("扫码登录写入二维码到文件失败:", err)
	}

	for {
		//0成功 1未扫描 2未确认 3已过期  -1系统错误
		status, err := qm.CheckQRCodeStatus()
		if err != nil {
			t.Fatal("扫码登录检测二维码状态失败:", err)
		}
		if status == 0 {
			break
		}
		t.Log("登录状态:", status)
		time.Sleep(2 * time.Second)
	}
	// 保存cookie以便下次使用
	cookie = qm.Info.Cookie
}
