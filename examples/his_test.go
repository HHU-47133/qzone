package examples

import (
	"encoding/base64"
	"errors"
	"github.com/HHU-47133/qzone"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	qm *qzone.QZone
)

func getParams(offset, count int) string {
	// 构造 params 参数
	params := map[string]string{
		"uin":                qm.Info.QQ,
		"begin_time":         "",
		"end_time":           "",
		"getappnotification": "1",
		"getnotifi":          "1",
		"has_get_key":        "0",
		"offset":             strconv.Itoa(offset), // todo start
		"set":                "0",
		"count":              strconv.Itoa(count), // todo count
		"useutf8":            "1",
		"outputhtmlfeed":     "1",
		"scope":              "1",
		"format":             "json",
		"g_tk":               qm.Info.Gtk,
	}

	u := url.Values{}
	for k, v := range params {
		u.Add(k, v)
	}
	return u.Encode()
}

// decodeHtml 解码原始的data string
func decodeHtml(dataStr string) string {
	// 1. 正则匹配 "\xHH" 的 16 进制编码部分
	re := regexp.MustCompile(`\\x[0-9a-fA-F]{2}`)

	// 替换每个匹配项
	decoded := re.ReplaceAllStringFunc(dataStr, func(hex string) string {
		// 去掉 "\x" 前缀，并解析为整数
		hexValue, err := strconv.ParseInt(hex[2:], 16, 32)
		if err != nil {
			// 如果解析失败，保留原字符串
			return hex
		}
		// 转换为字符
		return string(rune(hexValue))
	})

	// 2. 去除反斜杠定义
	re2 := regexp.MustCompile(`\\+`)
	decoded = re2.ReplaceAllStringFunc(decoded, func(match string) string {
		if match == `\/` { // \/ -> /
			return `/`
		}
		return `` // 否则，去除反斜杠
	})

	return decoded
}

// extractHtml 返回匹配html:'(.*?)'的html代码切片
func extractHtml(parsed string) []string {
	//匹配 html:'(.*?)'
	//re := regexp.MustCompile(`html:'(.*?)'`)
	re := regexp.MustCompile(`html:'(.*?)',opuin`)
	matches := re.FindAllStringSubmatch(parsed, -1)
	htmls := make([]string, len(matches))
	for idx, match := range matches {
		htmls[idx] = match[1]
	}

	return htmls
}

func TestGetHistoryData(t *testing.T) {
	qm.WithCookie(cookie)
	headers := map[string]string{
		"cookie":                    qm.Info.Cookie,
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"authority":                 "user.qzone.qq.com",
		"pragma":                    "no-cache",
		"cache-control":             "no-cache",
		"accept-language":           "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"sec-ch-ua":                 "\"Not A(Brand\";v=\"99\", \"Microsoft Edge\";v=\"121\", \"Chromium\";v=\"121\"",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "\"Windows\"",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"Content-Type":              "application/json; charset=utf-8",
	}

	urls := "https://user.qzone.qq.com/proxy/domain/ic2.qzone.qq.com/cgi-bin/feeds/feeds2_html_pav_all?" +
		getParams(0, 10)
	data, err := qzone.DialRequest(qzone.NewRequest(qzone.WithUrl(urls), qzone.WithHeader(headers)))
	if err != nil {
		t.Log(err)
	}

	htmlStr := string(data)
	t.Log(htmlStr)

}

// TODO 完成TestGetTotal，获取全部历史数据
func TestGetTotal(t *testing.T) {
	//login(t)
	//t.Log("cookie=", qm.Info.Cookie)
	qm = new(qzone.QZone)
	qm.WithCookie("pt2gguin=o0168880679;uin=o0168880679;skey=@wnAhLSn8f;superuin=o0168880679;supertoken=707819007;superkey=sppGD*9rWyWwfYssBl-LVe*s3EKovjLOSfbGrYtWJl0_;pt_recent_uins=220b2340fb277d07f7e4050c0a403e45f22dfb49f0db81ad62a0693ad04603d3e2a736682ccbe429660864dbc73c0ae4ec2d9c6617828c01;RK=SJlcxow2OC;ptnick_168880679=4a4c;ptcz=c232bd9097ded791b9f452f86ee426e1b9d6d8589ab31a1c7b99611f4ffbc068;uin=o0168880679;skey=@wnAhLSn8f;pt2gguin=o0168880679;p_uin=o0168880679;pt4_token=-WIkUxsuA7zF5wZnGUqULzOEWukmt1-f-9JYakerKSY_;p_skey=DujVvqO4JnMwHhDEGEBGcdltGxutrACfDCC2JLd14XI_;")
	baseUrl := "https://user.qzone.qq.com/proxy/domain/ic2.qzone.qq.com/cgi-bin/feeds/feeds2_html_pav_all?"
	headers := map[string]string{
		"cookie":                    qm.Info.Cookie,
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"authority":                 "user.qzone.qq.com",
		"pragma":                    "no-cache",
		"cache-control":             "no-cache",
		"accept-language":           "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"sec-ch-ua":                 "\"Not A(Brand\";v=\"99\", \"Microsoft Edge\";v=\"121\", \"Chromium\";v=\"121\"",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "\"Windows\"",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"Content-Type":              "application/json; charset=utf-8",
	}

	low, high := 0, 20000
	total, count := 0, 100

	for low <= high {
		mid := (low + high) >> 1
		// 1. 请求数据
		urls := baseUrl + getParams(mid*count, count)
		data, err := qzone.DialRequest(qzone.NewRequest(qzone.WithUrl(urls), qzone.WithHeader(headers)))
		t.Logf("[low=%v,high=%v,offset=%v]", low, high, mid*count)
		t.Log("data=", string(data))
		if err != nil {
			t.Logf("[low=%v,high=%v,offset=%v] request url failed, err:%v\n", low, high, mid*count, err)
			return
		}
		// 2. 解析数据
		ans := matchWithRegexp(string(data), `total_number:(.*?),`, true)
		if len(ans) == 0 {
			t.Logf("[low=%v,high=%v,offset=%v] parse data failed, err:%v\n", low, high, mid*count, err)
			return
		}
		num, _ := strconv.Atoi(ans[0])
		if num <= 0 {
			high = mid - 1
		} else { // num > 0
			low = mid + 1
			total = mid*count + num
			t.Logf("[total=%v]", total)
		}

		time.Sleep(2 * time.Second)
	}
	t.Log("final total=", total)
}

type qzoneHistoryItem struct {
	SenderQQ        string    // 发送方QQ
	ActionType      string    // 互动类型
	ShuoshuoID      string    // 说说ID
	ShuoshuoContent string    // 说说内容
	Content         string    // 互动内容
	CreateTime      time.Time // 发送的时间  abstime: data-abstime="1732591182"
	ImgUrls         []string  // 互动内容的图片
	ShuoshuoImgUrls []string  // 说说内容
	// QZoneImages	string // TODO: 可考虑加入表情
}

func extractHistoryMsg(html string) (*qzoneHistoryItem, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, errors.New("parse history msg failed")
	}
	var item *qzoneHistoryItem
	doc.Find("li.f-s-s").Each(func(i int, s *goquery.Selection) {
		// sender qq
		senderQQ, _ := s.Find(".user-avatar").Attr("link")
		senderQQ = strings.TrimPrefix(senderQQ, "nameCard_")
		// 说说ID
		shuoshuoID, _ := s.Find("i[name='feed_data']").Attr("data-tid")
		// 说说消息中的图片
		var shuoshuoImgUrls []string
		s.Find(".f-ct-txtimg .img-box img").Each(func(j int, imgS *goquery.Selection) {
			attr, exists := imgS.Attr("src")
			if exists {
				shuoshuoImgUrls = append(shuoshuoImgUrls, attr)
			}
		})
		// 说说内容
		shuoshuoContent := s.Find(".f-ct-txtimg .txt-box .txt-box-title").Contents().FilterFunction(func(j int, txtS *goquery.Selection) bool {
			// 过滤掉 <a> 和 <span> 等子标签，只保留纯文本节点
			return goquery.NodeName(txtS) == "#text"
		}).Text()
		shuoshuoContent = strings.TrimSpace(shuoshuoContent)

		// createTime
		createTimeStr, _ := s.Find("i[name='feed_data']").Attr("data-abstime")
		createTime, _ := strconv.ParseInt(createTimeStr, 10, 64)
		// 互动类型
		actionType := s.Find(".f-nick .state").Text()

		// 互动内容
		comments := s.Find(".comments-content").Text()
		suffix := s.Find(".comments-content .comments-op").Text()
		if len(comments) > 0 {
			comments = strings.SplitN(comments, ": ", 2)[1]
			comments = strings.TrimSuffix(comments, suffix)
		}

		// 互动消息中的图片
		var imgUrls []string
		s.Find(".comments-content .comments-thumbnails img").Each(func(j int, imgS *goquery.Selection) {
			attrOnLoad, exists := imgS.Attr("onload")
			if exists {
				link := matchWithRegexp(attrOnLoad, `trueSrc:'(.*?)'`, true)
				if len(link) > 0 {
					imgUrls = append(imgUrls, link[0])
				}
			}
		})

		item = &qzoneHistoryItem{
			SenderQQ:        senderQQ,
			ActionType:      actionType,
			ShuoshuoID:      shuoshuoID,
			Content:         comments,
			CreateTime:      time.Unix(createTime, 0),
			ImgUrls:         imgUrls,
			ShuoshuoContent: shuoshuoContent,
			ShuoshuoImgUrls: shuoshuoImgUrls,
		}
	})
	return item, nil
}

// matchWithRegexp 返回data中所有匹配pattern的字符串，extract为true时，仅返回匹配到内容
func matchWithRegexp(data, pattern string, extract bool) []string {
	re := regexp.MustCompile(pattern)
	matched := re.FindAllStringSubmatch(data, -1)
	if matched == nil {
		return nil
	}

	res := make([]string, len(matched))
	for i, match := range matched {
		if extract {
			res[i] = match[1]
		} else {
			res[i] = match[0]
		}
	}

	return res
}

func TestMatch(t *testing.T) {
	ans := matchWithRegexp(`{
	"code":-3000,
	"subcode":-4001,
	"message":"need login",
	"notice":0,
	"time":1732810983,
	"tips":"0103-87"
}
`, `"code":(.*?),`, true)
	t.Logf("%#v", ans)
}

func login(t *testing.T) {
	// 创建QZone对象, 使用扫码登录
	qm = qzone.NewQZone()
	b64s, err := qm.GenerateQRCode()
	if err != nil {
		t.Fatal("扫码登录获取二维码失败:", err)
	}

	ddd, err := base64.StdEncoding.DecodeString(b64s)
	if err != nil {
		t.Fatal("扫码登录base64解码失败:", err)
	}

	err = os.WriteFile("./qrcode.png", ddd, 0666)
	if err != nil {
		t.Fatal("扫码登录写入二维码到文件失败:", err)
	}

	for {
		//0成功 1未扫描 2未确认 3已过期  -1系统错误
		status, err := qm.CheckQRCodeStatus()
		if err != nil {
			t.Fatal("扫码登录检测二维码状态失败:", err)
		}
		if status == 0 {
			break
		}
		t.Log("登录状态:", status)
		time.Sleep(2 * time.Second)
	}
}

func TestGetQZoneHistoryList(t *testing.T) {
	login(t)
	t.Log("cookie=", qm.Info.Cookie)
	//qm = qzone.NewQZone()
	//qm.WithCookie("")
	list, err := qm.GetQZoneHistoryList()
	if err != nil {
		t.Log("test GetQZoneHistoryList failed, err:", err)
		return
	}
	t.Logf("len=%v", len(list))
	for _, item := range list {
		t.Logf("%#v", item)
	}
}
