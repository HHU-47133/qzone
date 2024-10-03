package examples

import (
	"encoding/base64"
	"github.com/HHU-47133/qzone"
	"os"
	"testing"
)

var (
	// cookie 登录成功后的 cookie
	cookie = "qrsig=ee0ab68371f81acebb61f376031ddbe5fe05788f442a8c1c1e5c6b45e8626a5de6cabe9100e9aec0046233f02b2c92ad18fd04ef76fd9651;uin=o1294222408;skey=@rh555XIz6;pt2gguin=o1294222408;p_uin=o1294222408;pt4_token=lYmEvOyIkgCUFQBGmFvmnVEUXaTbTmjJywe7nZm2Ofg_;p_skey=QUIv-nrrgdVK4BtkG1-uNNKBwtaDyrm3OPliG6dT4n0_;"
	//用于测试评论获取的说说tid
	tid = ""
	//用于测试的好友qq
	friendQQ = ""
	// ImageTestPath1 测试图片1路径
	ImageTestPath1 = "D:\\1.png"
	// ImageTestPath2 测试图片2路径
	ImageTestPath2 = "D:\\2.jpg"
)

// 调用低级别API获取指定数量说说
func TestGetPostListRaw(t *testing.T) {
	m := qzone.NewManager(cookie)
	ssl, err := m.ShuoShuoListRaw(friendQQ, 5, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	for i, shuoshuo := range ssl {
		t.Logf("got shuoshuo No.[%d]: %+v", i, shuoshuo)
		if i == 0 {
			tid = shuoshuo.Tid
		}
	}
}

// 调用高级别API获取全部说说
func TestGetPostList(t *testing.T) {
	m := qzone.NewManager(cookie)
	ssl, err := m.ShuoShuoList(friendQQ)
	if err != nil {
		t.Fatal(err)
	}
	for i, shuoshuo := range ssl {
		t.Logf("[全部说说获取成功] No.[%d]", i)
		t.Log(shuoshuo.Uin, shuoshuo.Name, shuoshuo.Content, shuoshuo.Pictotal)
	}
}

// 获取说说所有的一级评论
func TestGetComments(t *testing.T) {
	m := qzone.NewManager(cookie)
	comments, err := m.GetShuoShuoComments(tid)
	if err != nil {
		t.Log("get comments failed:", err)
	}
	for i, comment := range comments {
		t.Logf("got comment No.[%d]:%+v\n", i, comment)
	}
}

// 上传图片
func TestUploadImage(t *testing.T) {
	m := qzone.NewManager(cookie)
	// 读取本地图片
	srcByte, err := os.ReadFile(ImageTestPath1)
	if err != nil {
		t.Log("[读取本地图片失败]", err)
	}
	// base64编码
	picBase64 := base64.StdEncoding.EncodeToString(srcByte)
	// 上传图片
	uploadResult, err := m.UploadImage(picBase64)
	if err != nil {
		t.Log("[上传文件失败]", err)
	}
	t.Log("[上传图片成功]", uploadResult.URL)
}

// 获取说说总数
func TestShuoShuoCount(t *testing.T) {
	m := qzone.NewManager(cookie)
	cnt, err := m.GetShuoShuoCount(friendQQ)
	if err != nil {
		t.Fatal("[获取说说总数失败]", err)
	}
	t.Logf("[%s]获取说说总数成功:%d", friendQQ, cnt)
}

// 发布说说
func TestPublishShuoShuo(t *testing.T) {
	m := qzone.NewManager(cookie)
	// 1. 读取本地图片
	srcByte, err := os.ReadFile(ImageTestPath1)
	if err != nil {
		t.Log("[read image error]", err)
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
	srcByte, err = os.ReadFile(ImageTestPath2)
	if err != nil {
		t.Log("read image2 error", err)
	}
	pic2Base64 := base64.StdEncoding.EncodeToString(srcByte)
	_, _ = m.UploadImage(pic2Base64)
	// 4. 发说说
	publishResult, err := m.PublishShuoShuo("content", []string{pic1Base64, pic2Base64})
	if err != nil {
		t.Log("publish post error: ", err)
	}
	t.Log(publishResult)
}
