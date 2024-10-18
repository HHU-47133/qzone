package examples

import (
	"encoding/json"
	"fmt"
	"github.com/HHU-47133/qzone"
	"os"
	"strconv"
	"testing"
)

var resentShuoShuoData string

// 获取最新说说
func TestGetLatestShuoShuo(t *testing.T) {
	m := qzone.NewManager(Cfg.Cookie)
	ss, err := m.GetLatestShuoShuo(m.QQ)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[获取最新说说成功]", ss.Name, ss.Content, ss.Tid, ss.CreatedTime)
	Cfg.Tid = ss.Tid
}

// 获取说说总数
func TestGetShuoShuoCount(t *testing.T) {
	m := qzone.NewManager(Cfg.Cookie)
	cnt, err := m.GetShuoShuoCount(m.QQ)
	if err != nil {
		t.Fatal(err)
	}
	resentShuoShuoData = resentShuoShuoData + "【我正在测试自动化投稿】\n我说说总数是:" + strconv.FormatInt(cnt, 10)
	t.Log("[说说总数获取成功]" + strconv.FormatInt(cnt, 10))
}

// 获取指定说说一级评论总数
func TestGetLevel1CommentCount(t *testing.T) {
	m := qzone.NewManager(Cfg.Cookie)
	cnt, err := m.GetLevel1CommentCount(Cfg.Tid)
	if err != nil {
		t.Fatal(err)
	}
	resentShuoShuoData = resentShuoShuoData + "\n我上一条说说一级评论总数是:" + strconv.FormatInt(cnt, 10)
	t.Log("[成功获取说说一级评论总数]" + strconv.FormatInt(cnt, 10))
}

// 获取指定说说所有的一级评论
func TestShuoShuoCommentList(t *testing.T) {
	m := qzone.NewManager(Cfg.Cookie)
	cnt := int64(3)
	comments, _ := m.ShuoShuoCommentList(Cfg.Tid, cnt, 1000)
	resentShuoShuoData = resentShuoShuoData + "\n上条评论人是:"
	for i, comment := range comments {
		resentShuoShuoData = resentShuoShuoData + comment.OwnerName + " "
		t.Logf("[获取到说说评论][%d/%d]:%s %d %s %s", i, cnt, comment.OwnerName, comment.OwnerUin, comment.Content, comment.CreateTime)
	}
}

// 点赞说说
func TestDoLike(t *testing.T) {
	m := qzone.NewManager(Cfg.Cookie)
	dl, err := m.DoLike(Cfg.Tid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[点赞返回]" + dl.Msg)
}

//// 发布文字说说
//func TestPublishShuoShuoText(t *testing.T) {
//	m := qzone.NewManager(Cfg.Cookie)
//	pr, err := m.PublishShuoShuo(resentShuoShuoData, nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Logf("[发布文字说说返回]" + pr.Tid)
//}
//
//// 发布带图说说
//func TestPublishShuoShuoImg(t *testing.T) {
//	m := qzone.NewManager(Cfg.Cookie)
//	// 读取本地图片
//	srcByte, err := os.ReadFile(Cfg.ImgPath[0])
//	if err != nil {
//		t.Log("[测试图片1读取错误]", err)
//	}
//	// base64编码
//	pic1Base64 := base64.StdEncoding.EncodeToString(srcByte)
//
//	// 读取上传第二张图片
//	srcByte, err = os.ReadFile(Cfg.ImgPath[1])
//	if err != nil {
//		t.Log("[测试图片2读取错误]", err)
//	}
//	pic2Base64 := base64.StdEncoding.EncodeToString(srcByte)
//	// 发说说
//	pr, err := m.PublishShuoShuo("我正在测试QQ空间自动发说说功能", []string{pic1Base64, pic2Base64})
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Logf("[发布图片说说返回]" + pr.Tid)
//}

// 保存配置信息
func TestCfg(t *testing.T) {
	output, err := json.MarshalIndent(Cfg, "", "  ")
	if err != nil {
		t.Fatal("json配置修改解析失败:", err)
	}
	err = os.WriteFile("./config.json", output, 0644)
	if err != nil {
		fmt.Println("json配置文件写入失败:", err)
	}
}
