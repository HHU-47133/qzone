package examples

import (
	"encoding/base64"
	"fmt"
	"github.com/HHU-47133/qzone"
	"os"
	"testing"
)

var (
	// cookie 登录成功后的 cookie
	cookie = "pt2gguin=o1778046356;uin=o1778046356;skey=@NUftqZ3Sz;superuin=o1778046356;supertoken=1507384333;superkey=t7CIc38A*taBvketZpLctoUQYasJRLWR4XRP7M*4Gb4_;pt_recent_uins=e75c7ec9ec417d7ea917bd19da77e9a119369850a1ca2538d04eace4a4207714aca2f5199dfc5c74203d8a58865663cfc3e41f53e56d3350;RK=SqfoqwciGJ;ptnick_1778046356=e69e97c2b7e4b883e5a49c;ptcz=96c9e1cfde41dcc3ff599fa29a0d8ac47a01553e44f9635a86b00df6afe26456;uin=o1778046356;skey=@NUftqZ3Sz;pt2gguin=o1778046356;p_uin=o1778046356;pt4_token=bm1UAyvj9t1GJL7trkGCXPpZKzJl4ILTCk9DnpANpWE_;p_skey=pvhxGuxysp-fP3MLJNhAOHpFczuylP0jGL1y0JkPDZM_;"

	// ImageTestPath1 测试图片1路径
	ImageTestPath1 = "D:\\1.png"
	// ImageTestPath2 测试图片2路径
	ImageTestPath2 = "D:\\2.jpg"
)

// 获取所有的说说
func TestGetPostList(t *testing.T) {
	m := qzone.NewManager(cookie)
	ssl, err := m.ShuoShuoList(m.QQ, 20, 5)
	if err != nil {
		t.Fatal(err)
	}

	for _, shuoshuo := range ssl {
		fmt.Println(shuoshuo)
	}
}

// 获取说说所有的一级评论
func TestGetComments(t *testing.T) {
	m := qzone.NewManager(cookie)
	comments, err := m.GetShuoShuoComments("94d5fa69b8fcf766a1630b00")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("🧡🧡🧡评论结构体🧡🧡🧡：")
	for _, comment := range comments {
		fmt.Printf("%+v\n", comment)
	}
}

// 上传图片
func TestUploadImage(t *testing.T) {
	m := qzone.NewManager(cookie)
	// 读取本地图片
	srcByte, err := os.ReadFile(ImageTestPath1)
	if err != nil {
		t.Log("read image error", err)
	}
	// base64编码
	picBase64 := base64.StdEncoding.EncodeToString(srcByte)
	// 上传图片
	uploadResult, err := m.UploadImage(picBase64)
	if err != nil {
		t.Log("upload image error: ", err)
	}
	t.Log(uploadResult)
}

// 发布说说
func TestPublishShuoShuo(t *testing.T) {
	m := qzone.NewManager(cookie)
	// 1. 读取本地图片
	srcByte, err := os.ReadFile(ImageTestPath1)
	if err != nil {
		t.Log("read image error", err)
	}
	// 2. base64编码
	pic1Base64 := base64.StdEncoding.EncodeToString(srcByte)
	// 3. 上传图片
	uploadResult, err := m.UploadImage(pic1Base64)
	if err != nil {
		t.Log("upload image error: ", err)
	}
	t.Log(uploadResult)

	// 读取上传第二张图片
	srcByte, _ = os.ReadFile(ImageTestPath2)
	pic2Base64 := base64.StdEncoding.EncodeToString(srcByte)
	_, _ = m.UploadImage(pic2Base64)
	// 4. 发说说
	publishResult, err := m.PublishShuoShuo("content", []string{pic1Base64, pic2Base64})
	if err != nil {
		t.Log("publish post error: ", err)
	}
	t.Log(publishResult)
}
