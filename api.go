package qzone

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HHU-47133/qzone/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"log"
	"math"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// cReLike 点赞响应正则，frameElement.callback();
	cReLike = regexp.MustCompile(`(?s)frameElement.callback\((.*)\)`)
	// cRe 正则，_Callback();
	cRe = regexp.MustCompile(`(?s)_Callback\((.*)\)`)
)

// QQGroupList 群列表获取
func (q *QZone) QQGroupList() ([]*models.QQGroupResp, error) {
	gr := &models.QQGroupReq{
		Uin:     q.qq,
		Do:      "1",
		Rd:      fmt.Sprintf("%010.8f", rand.Float64()),
		Fupdate: "1",
		Clean:   "1",
		GTk:     q.gtk2,
	}
	url := getQQGroupURL + structToStr(gr)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"user-agent": ua,
		"cookie":     q.cookie,
	})))
	if err != nil {
		er := errors.New("QQ群请求错误:" + err.Error())
		log.Println("QQ群获取失败:", er.Error())
		return nil, er
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		er := errors.New("QQ群响应正则解析错误:" + string(data))
		log.Println("QQ群获取失败:", er.Error())
		return nil, er
	}
	jsonStr := r[1]
	resLen := gjson.Get(jsonStr, "data.group.#").Int()
	results := make([]*models.QQGroupResp, resLen)
	index := 0
	groups := gjson.Get(jsonStr, "data.group").Array()
	for _, group := range groups {
		gro := &models.QQGroupResp{
			GroupCode:   group.Get("groupcode").Int(),
			GroupName:   group.Get("groupname").String(),
			TotalMember: group.Get("total_member").Int(),
			NotFriends:  group.Get("notfriends").Int(),
		}
		results[index] = gro
		index++
	}
	return results, nil
}

// QQGroupMemberList 群友(非好友)列表获取
func (q *QZone) QQGroupMemberList(gid int64) ([]*models.QQGroupMemberResp, error) {
	gmr := &models.QQGroupMemberReq{
		Uin:     q.qq,
		Gid:     gid,
		Fupdate: "1",
		Type:    "1",
		GTk:     q.gtk2,
	}
	url := getQQGroupMemberURL + structToStr(gmr)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"user-agent": ua,
		"cookie":     q.cookie,
	})))
	if err != nil {
		er := errors.New("QQ群非好友请求错误:" + err.Error())
		log.Println("QQ群非好友获取失败:", er.Error())
		return nil, er
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		er := errors.New("QQ群非好友正则解析错误:" + string(data))
		log.Println("QQ群非好友获取失败:", er.Error())
		return nil, er
	}
	jsonStr := r[1]
	resLen := gjson.Get(jsonStr, "data.notfriends").Int()
	results := make([]*models.QQGroupMemberResp, resLen)
	index := 0
	groupMembers := gjson.Get(jsonStr, "data.friends").Array()
	for _, groupMember := range groupMembers {
		gro := &models.QQGroupMemberResp{
			Uin:       groupMember.Get("fuin").Int(),
			NickName:  groupMember.Get("name").String(),
			AvatarURL: groupMember.Get("img").String(),
		}
		gro.GroupCode = gjson.Get(jsonStr, "data.groupcode").Int()
		results[index] = gro
		index++
	}
	return results, nil
}

// FriendList 好友列表获取 TODO:有时候显示亲密度前200好友
func (q *QZone) FriendList() ([]*models.FriendInfoEasyResp, error) {
	url := fmt.Sprintf(friendURL, q.gtk2) + "&uin=" + strconv.FormatInt(q.qq, 10)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"referer": userQzoneURL,
		"origin":  userQzoneURL,
		"cookie":  q.cookie,
	})))
	if err != nil {
		er := errors.New("好友列表请求错误:" + err.Error())
		log.Println("好友列表获取失败:", er.Error())
		return nil, er
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		er := errors.New("好友列表正则解析错误:" + string(data))
		log.Println("好友列表获取失败:", er.Error())
		return nil, er
	}
	jsonStr := r[1]
	resLen := gjson.Get(jsonStr, "items.#").Int()
	results := make([]*models.FriendInfoEasyResp, resLen)
	index := 0

	friends := gjson.Get(jsonStr, "items").Array()
	for _, friend := range friends {
		fie := &models.FriendInfoEasyResp{
			Uin:     friend.Get("uin").Int(),
			Groupid: friend.Get("groupid").Int(),
			Name:    friend.Get("name").String(),
			Remark:  friend.Get("remark").String(),
			Image:   friend.Get("image").String(),
			Online:  friend.Get("online").Int(),
		}
		results[index] = fie
		index++
	}

	groupName := gjson.Get(jsonStr, "gpnames.#.gpname").Array()
	for i := 0; i < index; i++ {
		results[i].GroupName = groupName[results[i].Groupid].String()
	}
	return results, nil
}

// FriendInfoDetail 好友详细信息获取
func (q *QZone) FriendInfoDetail(uin int64) (*models.FriendInfoDetailResp, error) {
	url := fmt.Sprintf(detailFriendURL, q.gtk2) + "&uin=" + strconv.FormatUint(uint64(uin), 10)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"referer": userQzoneURL,
		"origin":  userQzoneURL,
		"cookie":  q.cookie,
	})))
	if err != nil {
		er := errors.New("好友详细信息请求错误:" + err.Error())
		log.Println("好友详细信息获取失败:", er.Error())
		return nil, er
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		er := errors.New("好友详细信息正则解析错误:" + string(data))
		log.Println("好友详细信息获取失败:", er.Error())
		return nil, er
	}
	jsonStr := r[1]

	fid := &models.FriendInfoDetailResp{}
	if err := json.Unmarshal([]byte(jsonStr), fid); err != nil {
		er := errors.New("好友详细信息JSON绑定错误:" + err.Error())
		log.Println("好友详细信息获取失败:", er.Error())
		return nil, er
	}
	return fid, nil
}

// PublishShuoShuo 发布说说，content文本内容，base64imgList图片数组
func (q *QZone) PublishShuoShuo(content string, base64imgList []string) (*models.ShuoShuoPublishResp, error) {
	var (
		uir         *models.UploadImageResp
		err         error
		picBo       string
		richval     string
		richvalList = make([]string, 0, 9)
		picBoList   = make([]string, 0, 9)
	)

	for _, base64img := range base64imgList {
		uir, err = q.uploadImage(base64img)
		if err != nil {
			log.Println("说说发布失败:", err.Error())
			return nil, err
		}
		picBo, richval, err = q.getPicBoAndRichval(uir)
		if err != nil {
			log.Println("说说发布失败:", err.Error())
			return nil, err
		}
		richvalList = append(richvalList, richval)
		picBoList = append(picBoList, picBo)
	}

	epr := models.EmotionPublishRequest{
		SynTweetVerson: "1",
		Paramstr:       "1",
		Who:            "1",
		Con:            content,
		Feedversion:    "1",
		Ver:            "1",
		UgcRight:       "1",
		ToSign:         "0",
		Hostuin:        q.qq,
		CodeVersion:    "1",
		Format:         "json",
		Qzreferrer:     userQzoneURL + "/" + strconv.FormatInt(q.qq, 10),
	}
	if len(base64imgList) > 0 {
		epr.Richtype = "1"
		epr.Richval = strings.Join(richvalList, "\t")
		epr.Subrichtype = "1"
		epr.PicBo = strings.Join(picBoList, ",")
	}
	url := fmt.Sprintf(emotionPublishURL, q.gtk2)
	payload := strings.NewReader(structToStr(epr))
	data, err := DialRequest(NewRequest(WithMethod("POST"), WithUrl(url), WithBody(payload),
		WithHeader(map[string]string{
			"referer": userQzoneURL,
			"origin":  userQzoneURL,
			"cookie":  q.cookie,
		})))
	if err != nil {
		er := errors.New("说说发布请求错误:" + err.Error())
		log.Println("说说发布失败:", er.Error())
		return nil, er
	}

	jsonStr := string(data)
	ssp := &models.ShuoShuoPublishResp{
		Code:    gjson.Get(jsonStr, "code").Int(),
		Tid:     gjson.Get(jsonStr, "tid").String(),
		Now:     gjson.Get(jsonStr, "now").Int(),
		Message: gjson.Get(jsonStr, "message").String(),
	}
	if ssp.Message != "" {
		er := errors.New("说说发布错误:" + ssp.Message)
		log.Println("说说发布失败:", er.Error())
		return nil, er
	}
	return ssp, nil
}

// ShuoShuoList 获取所有说说 实际能访问的说说个数 <= 说说总数(空间仅展示近半年等情况) (有空间访问权限即可)
func (q *QZone) ShuoShuoList(uin int64, num int64, ms int64) (ShuoShuo []*models.ShuoShuoResp, err error) {
	cnt := num
	t := int(math.Ceil(float64(cnt) / 20.0))
	var i int
	//获取最大数量，控制i的取值
	maxCnt, err := q.GetShuoShuoCount(uin)
	if err != nil {
		log.Println("说说获取失败:", err.Error())
		return nil, err
	}
	for range t {
		if i >= int(maxCnt) {
			break
		}
		ShuoShuoTemp, err := q.shuoShuoListRaw(uin, 20, i, 0)
		if err != nil {
			log.Println("所有说说获取失败:", err.Error())
			return nil, err
		}
		if len(ShuoShuoTemp) == 0 {
			break
		}
		if len(ShuoShuo) < int(cnt) {
			ShuoShuo = append(ShuoShuo, ShuoShuoTemp[0:min(len(ShuoShuoTemp), int(cnt)-len(ShuoShuo))]...)
			i = i + 20
			time.Sleep(time.Millisecond * time.Duration(ms))
		}
	}
	return ShuoShuo, nil
}

// GetShuoShuoCount 获取用户QQ号为uin的说说总数（有空间访问权限即可）
func (q *QZone) GetShuoShuoCount(uin int64) (cnt int64, err error) {
	mlr := models.MsgListRequest{
		Uin:                uin,
		Ftype:              "0",
		Sort:               "0",
		Pos:                "0",
		Num:                "1",
		Replynum:           "0",
		GTk:                q.gtk2,
		Callback:           "_preloadCallback",
		CodeVersion:        "1",
		Format:             "json",
		NeedPrivateComment: "1",
	}
	url := msglistURL + "?" + structToStr(mlr)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"referer": userQzoneURL,
		"origin":  userQzoneURL,
		"cookie":  q.cookie,
	})))
	if err != nil {
		er := errors.New("说说总数请求错误:" + err.Error())
		log.Println("说说总数获取失败:", er.Error())
		return -1, er
	}
	jsonStr := string(data)
	// 判断是否有访问权限
	forbid := gjson.Get(jsonStr, "message").String()
	if forbid != "" {
		er := errors.New("说说总数响应错误:" + forbid)
		log.Println("说说总数响应失败:", er.Error())
		return -1, er
	}
	cnt = gjson.Get(jsonStr, "total").Int()
	return cnt, nil
}

// GetLevel1CommentCount 获取一级评论总数(限制本人)
func (q *QZone) GetLevel1CommentCount(tid string) (cnt int64, err error) {
	url := fmt.Sprintf(getCommentsURL, strconv.FormatInt(q.qq, 10), 0, 1, tid, q.gtk2)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"cookie": q.cookie,
	})))
	if err != nil {
		er := errors.New("说说评论请求错误:" + err.Error())
		log.Println("说说评论请求失败:", er.Error())
		return -1, er
	}
	r := cRe.FindStringSubmatch(string(data))
	//log.Println("空指针异常测试：" + string(data))
	if len(r) < 2 {
		er := errors.New("说说评论正则解析错误:" + err.Error())
		log.Println("说说评论请求失败:", er.Error())
		return -1, er
	}
	jsonRaw := r[1]

	// 说说的一级评论总数
	numOfComments := gjson.Get(jsonRaw, "cmtnum").Int()
	return numOfComments, nil
}

// ShuoShuoCommentList 根据说说ID获取评论（限制本人）
func (q *QZone) ShuoShuoCommentList(tid string, num int64, ms int64) (comments []*models.Comment, err error) {
	numOfComments := num
	t := int(math.Ceil(float64(numOfComments) / 20.0))
	//获取最大数量，控制i的取值
	maxCnt, err := q.GetLevel1CommentCount(tid)
	if err != nil {
		log.Println("说说评论获取失败:", err.Error())
		return nil, err
	}
	var i int
	for range t {
		if i >= int(maxCnt) {
			break
		}
		commentsTemp, err := q.shuoShuoCommentsRaw(20, i, tid)
		if err != nil {
			log.Println("说说评论获取失败:", err.Error())
			return nil, err
		}
		if len(commentsTemp) == 0 {
			break
		}
		if len(comments) < int(num) {
			comments = append(comments, commentsTemp[0:min(len(commentsTemp), int(num)-len(comments))]...)
			i = i + 20
			time.Sleep(time.Millisecond * time.Duration(ms))
		}

	}
	return comments, nil
}

// GetLatestShuoShuo 获取用户QQ号为uin的最新说说（有空间访问权限即可）
func (q *QZone) GetLatestShuoShuo(uin int64) (*models.ShuoShuoResp, error) {
	ss, err := q.shuoShuoListRaw(uin, 1, 0, 0)
	fmt.Println("ss!!!!!!:", ss)
	if err != nil {
		er := errors.New("最新说说获取错误:" + err.Error())
		log.Println("最新说说获取失败:", er.Error())
		return nil, er
	}
	return ss[0], nil
}

// DoLike 说说空间点赞 TODO:疑似无效
func (q *QZone) DoLike(tid string) (*models.LikeResp, error) {
	lr := models.LikeRequest{
		Qzreferrer: userQzoneURL + strconv.FormatInt(q.qq, 10),
		Opuin:      q.qq,
		Unikey:     userQzoneURL + strconv.FormatInt(q.qq, 10) + "/mood/" + tid,
		From:       "1",
		Fid:        tid,
		Typeid:     "0",
		Appid:      "311",
	}
	lr.Curkey = lr.Unikey
	url := fmt.Sprintf(likeURL, q.gtk2)
	payload := strings.NewReader(structToStr(lr))
	data, err := DialRequest(NewRequest(WithMethod("POST"), WithUrl(url),
		WithBody(payload), WithHeader(map[string]string{
			"referer": userQzoneURL,
			"origin":  userQzoneURL,
			"cookie":  q.cookie,
		})))
	if err != nil {
		er := errors.New("点赞请求错误:" + err.Error())
		log.Println("空间点赞失败:", er.Error())
		return nil, er
	}
	r := cReLike.FindStringSubmatch(string(data))
	if len(r) < 2 {
		er := errors.New("点赞响应解析错误:" + string(data))
		log.Println("空间点赞失败:", er.Error())
		return nil, er
	}
	likeResp := &models.LikeResp{
		Ret: gjson.Get(r[1], "ret").Int(),
		Msg: gjson.Get(r[1], "msg").String(),
	}
	if likeResp.Msg != "succ" {
		er := errors.New("点赞未生效" + likeResp.Msg)
		log.Println("空间点赞失败:", er.Error())
		return nil, er
	}
	return likeResp, nil
}

// getQZoneHistoryList 获取QQ空间历史消息（限制本人），offset和count分别表示每次请求的偏移量和数目
func (q *QZone) getQZoneHistoryList(offset, count int64) ([]*models.QZoneHistoryItem, error) {
	// 匿名函数列表，完成子操作
	// decodeHtml 解码其中的html字符（例如\x3C）
	decodeHtml := func(dataStr string) string {
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
	// extractHtml 提取其中的html部分
	extractHtml := func(parsed string) []string {
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
	// extractHistoryMsg 解析html代码，提取一条消息的数据
	extractHistoryMsg := func(html string) (*models.QZoneHistoryItem, error) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return nil, errors.New("parse history msg failed")
		}
		var item *models.QZoneHistoryItem
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

			item = &models.QZoneHistoryItem{
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

	// 请求历史消息数据
	data, err := q.queryQZoneHistoryList(offset, count)
	if err != nil {
		er := errors.New("QQ空间历史数据请求错误:" + err.Error())
		log.Println("QQ空间历史数据请求失败:" + er.Error())
		return nil, er
	}
	// 解码并提取HTML数据
	htmlSlice := extractHtml(decodeHtml(string(data)))
	items := make([]*models.QZoneHistoryItem, len(htmlSlice))
	// 分别对html切片中的每一条数据进行处理
	for idx, html := range htmlSlice {
		item, err := extractHistoryMsg(html)
		if err != nil {
			er := errors.New("QQ空间历史数据解析错误:" + err.Error())
			log.Println("QQ空间历史数据解析失败:" + er.Error())
			return nil, er
		}
		items[idx] = item
	}
	return items, nil
}

func (q *QZone) queryQZoneHistoryList(offset, count int64) ([]byte, error) {
	qzhr := models.QZoneHistoryReq{
		Uin:                q.qq,
		Offset:             offset,
		Count:              count,
		BeginTime:          "",
		EndTime:            "",
		Getappnotification: "1",
		Getnotifi:          "1",
		HasGetKey:          "0",
		Useutf8:            "1",
		Outputhtmlfeed:     "1",
		Scope:              "1",
		Set:                "0",
		Format:             "json",
		Gtk:                q.gtk,
	}
	headers := map[string]string{
		"cookie":                    q.cookie,
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
	url_ := getQZoneHistory + structToStr(qzhr)

	data, err := DialRequest(NewRequest(WithUrl(url_), WithHeader(headers)))
	if err != nil {
		return nil, err
	}
	// cookie过期或者发生了其他错误
	ans := matchWithRegexp(string(data), `"code":(.*?),`, true)
	if ans != nil {
		code, _ := strconv.Atoi(ans[0])
		if code == -3000 {
			er := errors.New("请求历史消息数据错误: cookie失效或其他错误")
			log.Println("请求历史消息数据失败:", er.Error())
			return nil, er
		}
	}

	return data, nil
}

// GetQZoneHistoryList 获取本人QQ空间的所有历史消息
func (q *QZone) GetQZoneHistoryList() ([]*models.QZoneHistoryItem, error) {
	// 0. 函数中使用到的匿名函数
	// getTotal 获取历史消息总数
	getTotal := func() (int64, error) {
		var (
			low, high int64 = 0, 20000
			total     int64 = 0
			count     int64 = 100
		)

		for low <= high {
			mid := (low + high) >> 1
			// 1. 请求数据
			data, err := q.queryQZoneHistoryList(mid*count, count)
			if err != nil {
				er := errors.New("QQ空间历史消息获取错误:" + err.Error())
				log.Println("QQ空间历史消息获取失败:", er.Error())
				return total, er
			}
			// 2. 解析数据
			ans := matchWithRegexp(string(data), `total_number:(.*?),`, true)
			if ans == nil {
				er := errors.New("QQ空间历史消息解析错误")
				log.Println("QQ空间历史消息解析失败:", er.Error())
				return total, er
			}
			num, _ := strconv.ParseInt(ans[0], 10, 64)
			if num <= 0 {
				high = mid - 1
			} else { // num > 0
				low = mid + 1
				total = mid*count + num
			}

			time.Sleep(2 * time.Second)
		}

		return total, nil
	}

	// 1. getTotal
	total, err := getTotal()
	if err != nil {
		return nil, err
	}
	// 2. 每次请求10条，并拼接结果
	totalItems := make([]*models.QZoneHistoryItem, total)
	for i := 0; i <= int(total)/10; i++ {
		offset := i * 10
		items, err := q.getQZoneHistoryList(int64(offset), 10)
		if err != nil {
			return nil, err
		}
		for idx, item := range items {
			totalItems[offset+idx] = item
		}
	}

	return totalItems, nil
}
