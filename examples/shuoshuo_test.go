package examples

import (
	"encoding/base64"
	"github.com/HHU-47133/qzone"
	"os"
	"testing"
)

var (
	// cookie 登录成功后的 cookie
	cookie = "pt2gguin=o1294222408;uin=o1294222408;skey=@mfSWWpb16;superuin=o1294222408;supertoken=3435612033;superkey=0FogUHzNVgEpO2-y-oSliot1JlqWf2srnh6Sd9luQRE_;pt_recent_uins=a2bacb455090de383ebad516abb5cfde11b4b0a9b967ab2c247652cdccbdc0116f98e3ebe9f37908de8bc7397e7e9108a20e6a200c1d29a7;RK=SuFZBxEiGc;ptnick_1294222408=52;ptcz=617337e299936ec833f69206ac432e555a2634372b4def485b543a7fa9d52860;uin=o1294222408;skey=@mfSWWpb16;pt2gguin=o1294222408;p_uin=o1294222408;pt4_token=LxNb*jT4yInDpXQh901ONUN2xACkAcgmTj05mQKnZxc_;p_skey=5YWlqL0jheDJc8xnbJYrkqCldt-jr1zw6Y7vWprbDUI_;"
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
	ssl, err := m.ShuoShuoListRaw(m.QQ, 1, 0, 0)
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

// 获取最新说说
func TestLatestShuoShuo(t *testing.T) {
	m := qzone.NewManager(cookie)
	ss, err := m.GetLatestShuoShuo(m.QQ)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[获取最新说说成功]", ss.Name, ss.Content, ss.Tid, ss.CreatedTime)
}

// 调用高级别API获取全部说说
func TestGetPostList(t *testing.T) {
	m := qzone.NewManager(cookie)
	ssl, err := m.ShuoShuoList(m.QQ)
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
