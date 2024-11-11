package examples

import (
	"encoding/base64"
	"github.com/HHU-47133/qzone"
	"os"
	"testing"
	"time"
)

var (
	qm       = qzone.NewQZone()
	qrID     string
	b64s     string
	groupID  int64
	friendID int64
	tid      string
	imgPath  = [2]string{"./1.png", "./2.png"}
)

// 登录测试
func TestLogin(t *testing.T) {
	// 读取测试配置文件
	// 给一个userID用于获取二维码，成功返回base64数据和二维码id
	b64s, _ = qm.GenerateQRCode()
	ddd, _ := base64.StdEncoding.DecodeString(b64s) //成图片文件并把文件写入到buffer
	_ = os.WriteFile("./qrcode.png", ddd, 0666)
	for {
		status, err := qm.CheckQRCodeStatus()
		if err != nil {
			return
		}
		if status == 0 {
			break
		}
		t.Log("登录状态:", status)
		time.Sleep(2 * time.Second)
	}
}
