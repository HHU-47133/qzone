package qzone

import (
	"fmt"
	"net/http"
	"net/url"
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
)

// Ptqrshow 获得登录二维码
func Ptqrshow() (data []byte, qrsig string, ptqrtoken string, err error) {
	data, err = DialRequest(NewRequest(
		WithUrl(ptqrshowURL),
		WithClient(&http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}}),
		WithRespFunc(func(response *http.Response) {
			for _, v := range response.Cookies() {
				if v.Name == "qrsig" {
					qrsig = v.Value
					break
				}
			}
		})))
	if qrsig == "" {
		return
	}
	ptqrtoken = genderGTK(qrsig, 0)
	return
}

// Ptqrlogin 登录回调
func Ptqrlogin(qrsig string, qrtoken string) (data []byte, cookie string, err error) {
	urls := fmt.Sprintf(ptqrloginURL, qrtoken)
	data, err = DialRequest(NewRequest(
		WithUrl(urls),
		WithHeader(map[string]string{
			"cookie": "qrsig=" + qrsig,
		}),
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
	return
}

// LoginRedirect 登录成功回调
func LoginRedirect(redirectURL string) (cookie string, err error) {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return
	}
	values, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return
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
	return
}

// Manager qq空间信息管理
type Manager struct {
	Cookie string
	QQ     string
	Gtk    string
	Gtk2   string
	PSkey  string
	Skey   string
	Uin    string
}

// NewManager 初始化信息
func NewManager(cookie string) (m Manager) {
	cookie = strings.ReplaceAll(cookie, " ", "")
	for _, v := range strings.Split(cookie, ";") {
		name, val, f := strings.Cut(v, "=")
		if f {
			switch name {
			case "uin":
				m.Uin = val
			case "skey":
				m.Skey = val
			case "p_skey":
				m.PSkey = val
			}
		}
	}
	m.Gtk = genderGTK(m.Skey, 5381)
	m.Gtk2 = genderGTK(m.PSkey, 5381)
	m.QQ = strings.TrimPrefix(m.Uin, "o")
	m.Cookie = cookie
	return
}
