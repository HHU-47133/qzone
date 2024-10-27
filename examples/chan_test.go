package examples

import (
	"github.com/HHU-47133/qzone"
	"testing"
)

// 登录测试
func TestChan(t *testing.T) {
	img := make(chan []byte, 1)

	// 调用登录接口
	go func() {
		_, err := qzone.QzoneLogin("qr.png", img, 2)
		if err != nil {
			t.Fatal(err)
		}
	}()
	for a := range img {
		t.Logf("a")
		t.Log(a)
	}
}
