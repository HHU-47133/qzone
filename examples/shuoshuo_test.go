package examples

import (
	"encoding/base64"
	"github.com/HHU-47133/qzone"
	"os"
	"strconv"
	"testing"
)

var resentShuoShuoData string

// 获取最新说说
func TestGetLatestShuoShuo(t *testing.T) {
	qm := qzone.NewQZone().WithCookie(cookie)
	// 取出空间对象暴露的QQ值
	QQ := qm.Info.QQ
	uin, _ := strconv.ParseInt(QQ, 10, 64)
	ss, err := qm.GetLatestShuoShuo(uin)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[获取最新说说成功]", ss.Name, ss.Content, ss.Tid, ss.CreatedTime)
	// 保存说说ID备用
	tid = ss.Tid
}

// 获取说说总数
func TestGetShuoShuoCount(t *testing.T) {
	qm := qzone.NewQZone().WithCookie(cookie)
	// 取出空间对象暴露的QQ值
	QQ := qm.Info.QQ
	uin, _ := strconv.ParseInt(QQ, 10, 64)
	cnt, err := qm.GetShuoShuoCount(uin)
	if err != nil {
		t.Fatal(err)
	}
	resentShuoShuoData = resentShuoShuoData + "【我正在测试自动化投稿】\n我说说总数是:" + strconv.FormatInt(cnt, 10)
	t.Log("[说说总数获取成功]" + strconv.FormatInt(cnt, 10))
}

// 获取指定个数的说说
func TestShuoShuoList(t *testing.T) {
	qm := qzone.NewQZone().WithCookie(cookie)
	// 取出空间对象暴露的QQ值
	QQ := qm.Info.QQ
	cnt := int64(3)
	uin, _ := strconv.ParseInt(QQ, 10, 64)
	shuoshuos, _ := qm.ShuoShuoList(uin, cnt, 1000)
	for i, shuo := range shuoshuos {
		t.Logf("[获取到说说][%d/%d]:%s", i, cnt, shuo.Content)
	}
}

// 获取指定说说一级评论总数
// @{uin:2546229294,nick:爱莉希雅的,who:1}
func TestGetLevel1CommentCount(t *testing.T) {
	qm := qzone.NewQZone().WithCookie(cookie)
	cnt, err := qm.GetLevel1CommentCount(tid)
	if err != nil {
		t.Fatal(err)
	}
	resentShuoShuoData = resentShuoShuoData + "\n我上一条说说一级评论总数是:" + strconv.FormatInt(cnt, 10)
	t.Log("[成功获取说说一级评论总数]" + strconv.FormatInt(cnt, 10))
}

// 获取指定说说所有的一级评论
func TestShuoShuoCommentList(t *testing.T) {
	qm := qzone.NewQZone().WithCookie(cookie)
	cnt := int64(90)
	comments, _ := qm.ShuoShuoCommentList(tid, cnt, 1000)
	resentShuoShuoData = resentShuoShuoData + "\n上条评论人是:"
	for i, comment := range comments {
		resentShuoShuoData = resentShuoShuoData + comment.OwnerName + " "
		t.Logf("[获取到说说评论][%d/%d]:%s %d %s %s", i, cnt, comment.OwnerName, comment.OwnerUin, comment.Content, comment.PicContent)
	}
}

// 点赞说说
func TestDoLike(t *testing.T) {
	qm := qzone.NewQZone().WithCookie(cookie)
	dl, err := qm.DoLike(tid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[点赞返回]" + dl.Msg)
}

// 发布文字说说
func TestPublishShuoShuoText(t *testing.T) {
	qm := qzone.NewQZone().WithCookie(cookie)
	pr, err := qm.PublishShuoShuo(resentShuoShuoData, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("[发布文字说说返回]" + pr.Tid)
}

// 发布带图说说
func TestPublishShuoShuoImg(t *testing.T) {
	qm := qzone.NewQZone().WithCookie(cookie)
	// 读取本地图片
	srcByte, err := os.ReadFile(imgPath[0])
	if err != nil {
		t.Log("[测试图片1读取错误]", err)
	}
	// base64编码
	pic1Base64 := base64.StdEncoding.EncodeToString(srcByte)

	// 读取上传第二张图片
	srcByte, err = os.ReadFile(imgPath[1])
	if err != nil {
		t.Log("[测试图片2读取错误]", err)
	}
	pic2Base64 := base64.StdEncoding.EncodeToString(srcByte)
	// 发说说
	pr, err := qm.PublishShuoShuo("我正在测试QQ空间自动发说说功能", []string{pic1Base64, pic2Base64})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("[发布图片说说返回]" + pr.Tid)
}
