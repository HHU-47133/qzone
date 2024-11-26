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
	"time"
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
	// 获取QQ空间历史消息
	getQZoneHistory = "https://user.qzone.qq.com/proxy/domain/ic2.qzone.qq.com/cgi-bin/feeds/feeds2_html_pav_all?"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.SetPrefix("[qzone]")
}

// QZone QQ空间对象
// TODO:api调用之前是否可以便捷进行登录状态的判断?
type QZone struct {
	// 登录状态
	status int8 // 0 成功；1 未登录；2 已过期；
	// 暴露出去的字段
	Info info
	// 扫码登录流程使用
	qrLogin
	// API调用使用
	action
}

// info 暴露出去的字段
type info struct {
	QQ          string // QQ空间的账号
	Cookie      string // 登录成功的Cookie，保存以便下次使用
	Gtk         string // 临时添加，用于获取历史久远的空间数据拉取使用
	Gtk2        string // 临时添加，用于获取历史久远的空间数据拉取使用
	ExpiredTime time.Time
}

// qrLogin 扫码登录需要的字段
type qrLogin struct {
	qrsig   string // 二维码接口获取到的参数
	qrtoken string // 由qrsig计算而成
	cookie  string
}

// action 调用API需要使用的参数
type action struct {
	qq    int64 // QQ号
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

func (q *QZone) WithCookie(cookie string) *QZone {
	q.cookie = cookie
	q.unpack()
	// 设置为成功登录 TODO:可以做失效判断
	q.status = 0
	return q
}

// unpack 初始化信息,将成功扫码登录获取到的cookie解析
func (q *QZone) unpack() {
	cookie := strings.ReplaceAll(q.cookie, " ", "")
	for _, v := range strings.Split(cookie, ";") {
		name, val, f := strings.Cut(v, "=")
		if f {
			switch name {
			case "uin":
				q.uin = val
			case "skey":
				q.skey = val
			case "p_skey":
				q.pskey = val
			}
		}
	}
	q.gtk = genderGTK(q.skey, 5381)
	q.gtk2 = genderGTK(q.pskey, 5381)
	t, err := strconv.ParseInt(strings.TrimPrefix(q.uin, "o"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	q.qq = t
	q.cookie = cookie

	q.Info.Cookie = cookie
	q.Info.QQ = strings.TrimPrefix(q.uin, "o")
	q.Info.Gtk = q.gtk
	q.Info.Gtk2 = q.gtk2
	return
}

// GenerateQRCode 生成二维码，返回base64 二维码ID 用于查询扫码情况
func (q *QZone) GenerateQRCode() (string, error) {
	cookiesString := ""
	q.qrsig = ""
	data, err := DialRequest(NewRequest(
		WithUrl(ptqrshowURL),
		WithClient(&http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}}),
		WithRespFunc(func(response *http.Response) {
			for _, v := range response.Cookies() {
				cookiesString = cookiesString + v.String()
				if v.Name == "qrsig" {
					q.qrsig = v.Value
					break
				}
			}
		})))
	if err != nil {
		er := errors.New("空间登录二维码显示错误:" + string(data))
		return "", er
	}

	if q.qrsig == "" {
		er := errors.New("空间登录二维码cookie获取错误:" + cookiesString)
		return "", er
	}
	base64 := base64.StdEncoding.EncodeToString(data)
	q.qrtoken = genderGTK(q.qrsig, 0)
	return base64, nil
}

// CheckQRCodeStatus 检查二维码状态 //0成功 1未扫描 2未确认 3已过期  -1系统错误
func (q *QZone) CheckQRCodeStatus() (int8, error) {
	if q.status == 0 {
		return 0, nil
	}
	qrtoken := q.qrtoken
	qrsign := q.qrsig
	qcookie := q.cookie
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
		q.cookie = qcookie
		// 创建信息管理结构，携带登录回调cookie和重定向页面cookie
		q.unpack()
		q.status = 0
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
