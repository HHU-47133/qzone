package qzone

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.SetPrefix("[qzone]")
}

// Manager qq空间信息管理
type Manager struct {
	Cookie string
	QQ     int64
	Gtk    string
	Gtk2   string
	PSkey  string
	Skey   string
	Uin    string
}

// QzoneLogin 扫码登录
func QzoneLogin(qrCodeOutputPath string, qrCodeInBytes chan []byte, retryNum int64) (m Manager, err error) {
	index := int64(1)
Outer:
	for ; index <= retryNum; index++ {
		// 1. 获取二维码信息（data），取出cookie重要参数（qrsig、ptqrtoken）
		data, qrsig, qrtoken, err := qrShow()
		if err != nil {
			log.Printf("空间登录失败:%s[%d/%d]", err.Error(), index, retryNum)
			return Manager{}, err
		}
		// 2. 保存二维码
		if qrCodeOutputPath != "" {
			err = os.WriteFile(qrCodeOutputPath, data, 0666)
			if err != nil {
				log.Printf("空间登录失败:二维码保存错误:%s[%d/%d]", err.Error(), index, retryNum)
				return Manager{}, err
			}
		}
		// 3. 向通道发送二维码数据
		if qrCodeInBytes != nil {
			<-qrCodeInBytes
			qrCodeInBytes <- data
		}
		log.Printf("空间登录尝试中[%d/%d]\n", index, retryNum)
	Inner:
		for {
			//查询扫码结果
			data, qrloginCookie, err := qrLogin(qrsig, qrtoken)
			if err != nil {
				log.Printf("空间登录失败:%s[%d/%d]", err.Error(), index, retryNum)
				return Manager{}, err
			}
			text := string(data)
			switch {
			case strings.Contains(text, "二维码未失效"):
				log.Printf("空间登录二维码已生成，请尽快扫码登录[%d/%d]", index, retryNum)
			case strings.Contains(text, "二维码认证中"):
				log.Printf("空间登录二维码已扫描，请点击确认[%d/%d]", index, retryNum)
			case strings.Contains(text, "二维码已失效") || strings.Contains(text, "本次登录已被拒绝"):
				if index <= retryNum-1 {
					log.Printf("空间登录二维码已失效,即将重新生成[%d/%d]", index, retryNum)
				}
				break Inner
			case strings.Contains(text, "登录成功"):
				_ = os.Remove(qrCodeOutputPath)
				dealedCheckText := strings.ReplaceAll(text, "'", "")
				redirectURL := strings.Split(dealedCheckText, ",")[2]
				// 4. 成功登录后，获取登录重定向URL
				redirectCookie, err := loginRedirect(redirectURL)
				if err != nil {
					log.Printf("空间登录失败:%s[%d/%d]", err.Error(), index, retryNum)
					return Manager{}, err
				}
				cookie := qrloginCookie + redirectCookie
				// 创建信息管理结构，携带登录回调cookie和重定向页面cookie
				m = NewManager(cookie)
				log.Printf("空间登录成功[%d/%d]", index, retryNum)
				break Outer
			}
			time.Sleep(2 * time.Second)
		}

	}
	if index > retryNum {
		err := errors.New("已尝试最大次数")
		log.Printf("空间登录失败:%s[%d/%d]", err.Error(), retryNum, retryNum)
		return Manager{}, err
	}
	return m, nil
}

// qrShow 获得登录二维码
func qrShow() (data []byte, qrsig string, qrtoken string, err error) {
	cookiesString := ""
	data, err = DialRequest(NewRequest(
		WithUrl(ptqrshowURL),
		WithClient(&http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}}),
		WithRespFunc(func(response *http.Response) {
			for _, v := range response.Cookies() {
				cookiesString = cookiesString + v.String()
				if v.Name == "qrsig" {
					qrsig = v.Value
					break
				}
			}
		})))
	if err != nil {
		er := errors.New("空间登录二维码显示错误:" + string(data))
		return nil, "", "", er
	}
	if qrsig == "" {
		er := errors.New("空间登录二维码cookie获取错误:" + cookiesString)
		return nil, "", "", er
	}
	qrtoken = genderGTK(qrsig, 0)
	return
}

// qrLogin 登录状态检测
func qrLogin(qrsig string, qrtoken string) (data []byte, cookie string, err error) {
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
	if err != nil {
		er := errors.New("空间登录状态检测错误:" + err.Error())
		return nil, "", er
	}
	return
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
	t, err := strconv.ParseInt(strings.TrimPrefix(m.Uin, "o"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	m.QQ = t
	m.Cookie = cookie
	return
}
