package qzone

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	userQzoneURL      = "https://user.qzone.qq.com"
	ua                = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
	contentType       = "application/x-www-form-urlencoded"
	params            = "g_tk=%v"
	inpcqqURL         = "https://h5.qzone.qq.com/feeds/inpcqq?uin=%v&qqver=5749&timestamp=%v"
	emotionPublishURL = userQzoneURL + "/proxy/domain/taotao.qzone.qq.com/cgi-bin/emotion_cgi_publish_v6?" + params
	uploadImageURL    = "https://up.qzone.qq.com/cgi-bin/upload/cgi_upload_image?" + params
	msglistURL        = userQzoneURL + "/proxy/domain/taotao.qq.com/cgi-bin/emotion_cgi_msglist_v6"
	likeURL           = userQzoneURL + "/proxy/domain/w.qzone.qq.com/cgi-bin/likes/internal_dolike_app?" + params
	ptqrshowURL       = "https://ssl.ptlogin2.qq.com/ptqrshow?appid=549000912&e=2&l=M&s=3&d=72&v=4&t=0.31232733520361844&daid=5&pt_3rd_aid=0"
	ptqrloginURL      = "https://xui.ptlogin2.qq.com/ssl/ptqrlogin?u1=https://qzs.qq.com/qzone/v5/loginsucc.html?para=izone&ptqrtoken=%v&ptredirect=0&h=1&t=1&g=1&from_ui=1&ptlang=2052&action=0-0-1656992258324&js_ver=22070111&js_type=1&login_sig=&pt_uistyle=40&aid=549000912&daid=5&has_onekey=1&&o1vId=1e61428d61cb5015701ad73d5fb59f73"
	checkSigURL       = "https://ptlogin2.qzone.qq.com/check_sig?pttype=1&uin=%v&service=ptqrlogin&nodirect=1&ptsigx=%v&s_url=https://qzs.qq.com/qzone/v5/loginsucc.html?para=izone&f_url=&ptlang=2052&ptredirect=100&aid=549000912&daid=5&j_later=0&low_login_hour=0&regmaster=0&pt_login_type=3&pt_aid=0&pt_aaid=16&pt_light=0&pt_3rd_aid=0"
	friendURL         = "https://h5.qzone.qq.com/proxy/domain/r.qzone.qq.com/cgi-bin/tfriend/friend_show_qqfriends.cgi?" + params
	detailFriendURL   = "https://h5.qzone.qq.com/proxy/domain/base.qzone.qq.com/cgi-bin/user/cgi_userinfo_get_all?" + params
	getCommentsURL    = "https://h5.qzone.qq.com/proxy/domain/taotao.qq.com/cgi-bin/emotion_cgi_msgdetail_v6?uin=%s&pos=%d&num=%d&tid=%s&format=jsonp&g_tk=%s"
	// 获取点赞列表的URL
	getLikeListURL = "https://h5.qzone.qq.com/proxy/domain/users.qzone.qq.com/cgi-bin/likes/get_like_list_app?"
	// 获取QQ群URL
	getQQGroupURL = "https://user.qzone.qq.com/proxy/domain/r.qzone.qq.com/cgi-bin/tfriend/qqgroupfriend_extend.cgi?"
	// 获取QQ群成员非好友URL
	getQQGroupMemberURL = "https://user.qzone.qq.com/proxy/domain/r.qzone.qq.com/cgi-bin/tfriend/qqgroupfriend_groupinfo.cgi?"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.SetPrefix("[qzone]")
}

type QZone struct {
	qrsig   string // 二维码接口获取到的参数
	qrtoken string // 由Qrsig计算而成
	cookie  string
	status  int8 // 0 成功；1 未登录；2 已过期；

	qq    int64
	gtk   string
	gtk2  string
	pskey string
	skey  string
	uin   string
}

// NewQZone 创建管理类
func NewQZone() *QZone {
	return &QZone{
		status: 1,
	}
}

func (qz *QZone) WithCookie(cookie string) *QZone {
	qz.cookie = cookie
	qz.unpack()
	qz.status = 0
	return qz
}

// unpack 初始化信息
func (qz *QZone) unpack() {
	cookie := strings.ReplaceAll(qz.cookie, " ", "")
	for _, v := range strings.Split(cookie, ";") {
		name, val, f := strings.Cut(v, "=")
		if f {
			switch name {
			case "uin":
				qz.uin = val
			case "skey":
				qz.skey = val
			case "p_skey":
				qz.pskey = val
			}
		}
	}
	qz.gtk = genderGTK(qz.skey, 5381)
	qz.gtk2 = genderGTK(qz.pskey, 5381)
	t, err := strconv.ParseInt(strings.TrimPrefix(qz.uin, "o"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	qz.qq = t
	qz.cookie = cookie
	return
}

// GenerateQRCode 生成二维码，返回base64 二维码ID 用于查询扫码情况
func (qz *QZone) GenerateQRCode() (string, error) {
	cookiesString := ""
	qz.qrsig = ""
	data, err := DialRequest(NewRequest(
		WithUrl(ptqrshowURL),
		WithClient(&http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}}),
		WithRespFunc(func(response *http.Response) {
			for _, v := range response.Cookies() {
				cookiesString = cookiesString + v.String()
				if v.Name == "qrsig" {
					qz.qrsig = v.Value
					break
				}
			}
		})))
	if err != nil {
		er := errors.New("空间登录二维码显示错误:" + string(data))
		return "", er
	}

	if qz.qrsig == "" {
		er := errors.New("空间登录二维码cookie获取错误:" + cookiesString)
		return "", er
	}
	base64 := base64.StdEncoding.EncodeToString(data)
	qz.qrtoken = genderGTK(qz.qrsig, 0)
	return base64, nil
}

// CheckQRCodeStatus 检查二维码状态 //0成功 1未扫描 2未确认 3已过期   -1系统错误
func (qz *QZone) CheckQRCodeStatus() (int8, error) {
	if qz.status == 0 {
		return 0, nil
	}

	qrtoken := qz.qrtoken
	qrsign := qz.qrsig
	qcookie := qz.cookie
	urls := fmt.Sprintf(ptqrloginURL, qrtoken)
	data, err := DialRequest(NewRequest(
		WithUrl(urls),
		WithHeader(map[string]string{
			"cookie": "qrsig=" + qrsign,
		}),
		WithClient(&http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}}),
		WithRespFunc(func(response *http.Response) {
			for _, v := range response.Cookies() {
				if v.Value != "" {
					qcookie += v.Name + "=" + v.Value + ";"
				}
			}
		})))
	if err != nil {
		er := errors.New("空间登录状态检测错误:" + err.Error())
		return -1, er
	}
	text := string(data)
	switch {
	case strings.Contains(text, "二维码未失效"):
		return 1, nil
	case strings.Contains(text, "二维码认证中"):
		return 2, nil
	case strings.Contains(text, "二维码已失效") || strings.Contains(text, "本次登录已被拒绝"):
		return 3, nil
	case strings.Contains(text, "登录成功"):
		dealedCheckText := strings.ReplaceAll(text, "'", "")
		redirectURL := strings.Split(dealedCheckText, ",")[2]
		redirectCookie, err := loginRedirect(redirectURL)
		if err != nil {
			er := errors.New("空间登录重定向失败:" + err.Error())
			return -1, er
		}
		qcookie += redirectCookie
		qz.cookie = qcookie
		// 创建信息管理结构，携带登录回调cookie和重定向页面cookie
		qz.unpack()
		qz.status = 0
		return 0, nil
	}
	return 0, nil
}

// loginRedirect 登录成功回调
func loginRedirect(redirectURL string) (cookie string, err error) {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return "", errors.New("空间登录重定向链接解析错误:" + err.Error())
	}
	values, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", errors.New("空间登录重定向链接查询参数解析错误:" + err.Error())
	}

	urls := fmt.Sprintf(checkSigURL, values["uin"][0], values["ptsigx"][0])
	_, err = DialRequest(NewRequest(
		WithUrl(urls),
		WithClient(&http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}}),
		WithRespFunc(func(response *http.Response) {
			for _, v := range response.Cookies() {
				if v.Value != "" {
					cookie += v.Name + "=" + v.Value + ";"
				}
			}
		})))
	if err != nil {
		return "", errors.New("空间登录重定向链接请求错误:" + err.Error())
	}
	return
}
