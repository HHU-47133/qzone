package examples

import (
	"errors"
	"fmt"
	"github.com/HHU-47133/qzone"
	"github.com/PuerkitoBio/goquery"
	"net/url"
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

// TODO 完成getTotal，获取全部历史数据
func getTotal() int {
	l, r := 0, 100000
	var m int
	baseUrl := "https://user.qzone.qq.com/proxy/domain/ic2.qzone.qq.com/cgi-bin/feeds/feeds2_html_pav_all?"
	header := map[string]string{
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
	for l <= r {
		m = (l + r) >> 1
		fmt.Printf("l=%v,r=%v\n", l, r)
		url_ := baseUrl + getParams(m, 100)
		data, err := qzone.DialRequest(qzone.NewRequest(qzone.WithUrl(url_), qzone.WithHeader(header)))
		if err != nil {
			fmt.Println("get response failed, err:", err)
			return -1
		}
		if strings.Contains(string(data), "li") {
			l = m + 1
		} else {
			r = m - 1
		}
	}

	return m
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
