package qzone

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HHU-47133/qzone/models"
	"github.com/tidwall/gjson" // TODO: 疑问？需要处理 body 中的 \n
	"math"
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

// PublishShuoShuo 发布说说，content文本内容，base64imgList图片数组
func (m *Manager) PublishShuoShuo(content string, base64imgList []string) (*models.ShuoShuoPublishResp, error) {
	var (
		uir         *models.UploadImageResp
		err         error
		picBo       string
		richval     string
		richvalList = make([]string, 0, 9)
		picBoList   = make([]string, 0, 9)
	)

	for _, base64img := range base64imgList {
		uir, err = m.UploadImage(base64img)
		if err != nil {
			return nil, err
		}
		picBo, richval, err = getPicBoAndRichval(uir)
		if err != nil {
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
		Hostuin:        m.QQ,
		CodeVersion:    "1",
		Format:         "json",
		Qzreferrer:     userQzoneURL + "/" + m.QQ,
	}
	if len(base64imgList) > 0 {
		epr.Richtype = "1"
		epr.Richval = strings.Join(richvalList, "\t")
		epr.Subrichtype = "1"
		epr.PicBo = strings.Join(picBoList, ",")
	}
	return m.publishShuoShuo(epr)
}

// publishShuoShuo 发送说说
func (m *Manager) publishShuoShuo(epr models.EmotionPublishRequest) (*models.ShuoShuoPublishResp, error) {
	url := fmt.Sprintf(emotionPublishURL, m.Gtk2)
	payload := strings.NewReader(structToStr(epr))
	data, err := DialRequest(NewRequest(WithMethod("POST"), WithUrl(url), WithBody(payload),
		WithHeader(map[string]string{
			"referer": userQzoneURL,
			"origin":  userQzoneURL,
			"cookie":  m.Cookie,
		})))
	if err != nil {
		return nil, err
	}

	jsonStr := string(data)
	ssp := &models.ShuoShuoPublishResp{
		Code: gjson.Get(jsonStr, "code").Int(),
		Tid:  gjson.Get(jsonStr, "tid").String(),
		Now:  gjson.Get(jsonStr, "now").Int(),
		//Feedinfo: gjson.Get(jsonStr, "feedinfo").String(),
		Message: gjson.Get(jsonStr, "message").String(),
	}

	return ssp, nil
}

// ShuoShuoList 获取用户QQ号为uin且最多num个说说列表，每个说说获取上限replynum个评论数量
func (m *Manager) ShuoShuoList(uin string, num int, replynum int) ([]*models.ShuoShuoResp, error) {
	mlr := models.MsgListRequest{
		Uin:                uin,
		Ftype:              "0",
		Sort:               "0",
		Num:                strconv.Itoa(num),
		Replynum:           strconv.Itoa(replynum),
		GTk:                m.Gtk2,
		Callback:           "_preloadCallback",
		CodeVersion:        "1",
		Format:             "json",
		NeedPrivateComment: "1",
	}

	url := msglistURL + "?" + structToStr(mlr)

	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"referer": userQzoneURL,
		"origin":  userQzoneURL,
		"cookie":  m.Cookie,
	})))
	if err != nil {
		return nil, err
	}
	jsonStr := string(data)

	resLen := gjson.Get(jsonStr, "msglist.#").Int()
	results := make([]*models.ShuoShuoResp, min(resLen, int64(num)))
	index := 0

	lists := gjson.Get(jsonStr, "msglist").Array()
	for _, shuoshuo := range lists {
		ss := &models.ShuoShuoResp{
			Uin:         shuoshuo.Get("uin").Int(),
			Name:        shuoshuo.Get("name").String(),
			Tid:         shuoshuo.Get("tid").String(),
			Content:     shuoshuo.Get("content").String(),
			CreateTime:  shuoshuo.Get("createTime").String(),
			CreatedTime: shuoshuo.Get("created_time").Int(),
			Pictotal:    shuoshuo.Get("pictotal").Int(),
			Cmtnum:      shuoshuo.Get("cmtnum").Int(),
			Secret:      shuoshuo.Get("secret").Int(),
		}

		pics := shuoshuo.Get("pic").Array()
		for _, pic := range pics {
			ss.Pic = append(ss.Pic, models.PicResp{
				PicId:      pic.Get("pic_id").String(),
				Url1:       pic.Get("url1").String(),
				Url2:       pic.Get("url2").String(),
				Url3:       pic.Get("url3").String(),
				Smallurl:   pic.Get("smallurl").String(),
				Curlikekey: pic.Get("curlikekey").String(),
			})
		}

		results[index] = ss
		index++
	}

	return results, nil
}

// FriendList 获取亲密度前200的好友，第一位好友是小Q（系统）
func (m *Manager) FriendList() ([]*models.FriendInfoEasyResp, error) {
	url := fmt.Sprintf(friendURL, m.Gtk2) + "&uin=" + m.QQ
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"referer": userQzoneURL,
		"origin":  userQzoneURL,
		"cookie":  m.Cookie,
	})))
	if err != nil {
		return nil, err
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		return nil, errors.New("好友正则解析错误")
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

// FriendInfoDetail 获取好友详细信息
func (m *Manager) FriendInfoDetail(uin int64) (*models.FriendInfoDetailResp, error) {
	url := fmt.Sprintf(detailFriendURL, m.Gtk2) + "&uin=" + strconv.FormatUint(uint64(uin), 10)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"referer": userQzoneURL,
		"origin":  userQzoneURL,
		"cookie":  m.Cookie,
	})))
	if err != nil {
		return nil, err
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		return nil, errors.New("好友正则解析错误")
	}
	jsonStr := r[1]

	fid := &models.FriendInfoDetailResp{}
	if err := json.Unmarshal([]byte(jsonStr), fid); err != nil {
		return nil, err
	}
	return fid, nil
}

// UploadImage 上传图片
func (m *Manager) UploadImage(base64img string) (*models.UploadImageResp, error) {
	uir := models.UploadImageRequest{
		Filename:      "filename",
		Uin:           m.QQ,
		Skey:          m.Skey,
		Zzpaneluin:    m.QQ,
		PUin:          m.QQ,
		PSkey:         m.PSkey,
		Uploadtype:    "1",
		Albumtype:     "7",
		Exttype:       "0",
		Refer:         "shuoshuo",
		OutputType:    "json",
		Charset:       "utf-8",
		OutputCharset: "utf-8",
		UploadHd:      "1",
		HdWidth:       "2048",
		HdHeight:      "10000",
		HdQuality:     "96",
		BackUrls:      "http://upbak.photo.qzone.qq.com/cgi-bin/upload/cgi_upload_image,http://119.147.64.75/cgi-bin/upload/cgi_upload_image",
		URL:           fmt.Sprintf(uploadImageURL, m.Gtk2),
		Base64:        "1",
		Picfile:       base64img,
		Qzreferrer:    userQzoneURL + "/" + m.QQ,
	}

	url := fmt.Sprintf(uploadImageURL, m.Gtk2)
	payload := strings.NewReader(structToStr(uir))
	data, err := DialRequest(NewRequest(WithMethod("POST"), WithUrl(url), WithBody(payload),
		WithHeader(map[string]string{
			"referer": userQzoneURL,
			"origin":  userQzoneURL,
			"cookie":  m.Cookie,
		})))
	if err != nil {
		return nil, err
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		return nil, errors.New("上传失败")
	}
	jsonStr := r[1]
	uploadImageResp := &models.UploadImageResp{
		Pre:        gjson.Get(jsonStr, "data.pre").String(),
		URL:        gjson.Get(jsonStr, "data.url").String(),
		Width:      gjson.Get(jsonStr, "data.width").Int(),
		Height:     gjson.Get(jsonStr, "data.height").Int(),
		OriginURL:  gjson.Get(jsonStr, "data.origin_url").String(),
		Contentlen: gjson.Get(jsonStr, "data.contentlen").Int(),
		Ret:        gjson.Get(jsonStr, "ret").Int(),
		Albumid:    gjson.Get(jsonStr, "data.albumid").String(),
		Lloc:       gjson.Get(jsonStr, "data.lloc").String(),
		Sloc:       gjson.Get(jsonStr, "data.sloc").String(),
		Type:       gjson.Get(jsonStr, "data.type").Int(),
	}
	return uploadImageResp, nil
}

// DoLike 空间点赞
func (m *Manager) DoLike(tid string) (*models.LikeResp, error) {
	lr := models.LikeRequest{
		Qzreferrer: userQzoneURL + m.QQ,
		Opuin:      m.QQ,
		Unikey:     userQzoneURL + m.QQ + "/mood/" + tid,
		From:       "1",
		Fid:        tid,
		Typeid:     "0",
		Appid:      "311",
	}
	lr.Curkey = lr.Unikey
	return m.LikeRaw(lr)
}

// LikeRaw 空间点赞
func (m *Manager) LikeRaw(lr models.LikeRequest) (*models.LikeResp, error) {
	url := fmt.Sprintf(likeURL, m.Gtk2)
	payload := strings.NewReader(structToStr(lr))
	data, err := DialRequest(NewRequest(WithMethod("POST"), WithUrl(url),
		WithBody(payload), WithHeader(map[string]string{
			"referer": userQzoneURL,
			"origin":  userQzoneURL,
			"cookie":  m.Cookie,
		})))
	if err != nil {
		return nil, err
	}
	r := cReLike.FindStringSubmatch(string(data))
	if len(r) < 2 {
		return nil, errors.New("点赞正则解析失败")
	}
	likeResp := &models.LikeResp{
		Ret: gjson.Get(r[1], "ret").Int(),
		Msg: gjson.Get(r[1], "msg").String(),
	}
	if likeResp.Msg != "succ" {
		return nil, errors.New(fmt.Sprintf("点赞失败：[%s]", likeResp.Msg))
	}
	return likeResp, nil
}

// GetShuoShuoComments 根据说说ID获取所有评论，仅限本用户（m.QQ） // TODO: 待测试
func (m *Manager) GetShuoShuoComments(tid string) (comments []*models.Comment, err error) {
	url := fmt.Sprintf(getCommentsURL, m.QQ, 1, 1, tid, m.Gtk2)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"cookie": m.Cookie,
	})))
	if err != nil {
		return nil, err
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		return nil, errors.New("说说评论正则解析错误")
	}
	jsonRaw := r[1]

	// 说说的一级评论总数
	numOfComments := gjson.Get(jsonRaw, "cmtnum").Int()
	t := int(math.Ceil(float64(numOfComments) / 20.0))
	var i int
	for range t {
		commentsTemp, err := m.getShuoShuoCommentsRaw(20, i, tid)
		if err != nil {
			return nil, err
		}
		comments = append(comments, commentsTemp...)
		i = i + 20
	}
	return comments, nil
}

// getShuoShuoCommentsRaw 从第pos条评论开始获取num条评论，num最大为20
func (m *Manager) getShuoShuoCommentsRaw(num int, pos int, tid string) (comments []*models.Comment, err error) {
	url := fmt.Sprintf(getCommentsURL, m.QQ, pos, num, tid, m.Gtk2)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"cookie": m.Cookie,
	})))
	if err != nil {
		return nil, err
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		return nil, errors.New("说说评论正则解析错误")
	}
	jsonRaw := r[1]
	// fmt.Println("🧡🧡🧡取说说评论测试🧡🧡🧡：", jsonRaw)

	// 取出评论数据
	commentJsonList := gjson.Get(jsonRaw, "commentlist").Array()
	for _, com := range commentJsonList {
		comment := &models.Comment{
			ShuoShuoID: tid,
			OwnerName:  com.Get("owner.name").String(),
			OwnerUin:   com.Get("owner.uin").Int(),
			Content:    com.Get("content").String(),
			PicContent: make([]string, 0),
			CreateTime: time.Unix(com.Get("create_time").Int(), 0),
		}
		// 添加图片评论的图片到结构体
		for _, pic := range com.Get("rich_info").Array() {
			comment.PicContent = append(comment.PicContent, pic.Get("burl").String())
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// getPicBoAndRichval 上传图片资源
func getPicBoAndRichval(data *models.UploadImageResp) (picBo, richval string, err error) {
	var flag bool
	if data.Ret != 0 {
		err = errors.New("上传失败")
		return
	}
	_, picBo, flag = strings.Cut(data.URL, "&bo=")
	if !flag {
		err = errors.New("上传图片返回的地址错误")
		return
	}
	richval = fmt.Sprintf(",%s,%s,%s,%d,%d,%d,,%d,%d", data.Albumid, data.Lloc, data.Sloc, data.Type, data.Height, data.Width, data.Height, data.Width)
	return
}
