package qzone

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HHU-47133/qzone/models"
	"github.com/tidwall/gjson" // TODO: 疑问？需要处理 body 中的 \n
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
func (m *Manager) QQGroupList() ([]*models.QQGroupResp, error) {
	gr := &models.QQGroupReq{
		Uin:     m.QQ,
		Do:      "1",
		Rd:      fmt.Sprintf("%010.8f", rand.Float64()),
		Fupdate: "1",
		Clean:   "1",
		GTk:     m.Gtk2,
	}
	url := getQQGroupURL + structToStr(gr)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"user-agent": ua,
		"cookie":     m.Cookie,
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
func (m *Manager) QQGroupMemberList(gid int64) ([]*models.QQGroupMemberResp, error) {
	gmr := &models.QQGroupMemberReq{
		Uin:     m.QQ,
		Gid:     gid,
		Fupdate: "1",
		Type:    "1",
		GTk:     m.Gtk2,
	}
	url := getQQGroupMemberURL + structToStr(gmr)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"user-agent": ua,
		"cookie":     m.Cookie,
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
func (m *Manager) FriendList() ([]*models.FriendInfoEasyResp, error) {
	url := fmt.Sprintf(friendURL, m.Gtk2) + "&uin=" + strconv.FormatInt(m.QQ, 10)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"referer": userQzoneURL,
		"origin":  userQzoneURL,
		"cookie":  m.Cookie,
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
func (m *Manager) FriendInfoDetail(uin int64) (*models.FriendInfoDetailResp, error) {
	url := fmt.Sprintf(detailFriendURL, m.Gtk2) + "&uin=" + strconv.FormatUint(uint64(uin), 10)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"referer": userQzoneURL,
		"origin":  userQzoneURL,
		"cookie":  m.Cookie,
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
		uir, err = m.uploadImage(base64img)
		if err != nil {
			log.Println("说说发布失败:", err.Error())
			return nil, err
		}
		picBo, richval, err = m.getPicBoAndRichval(uir)
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
		Hostuin:        m.QQ,
		CodeVersion:    "1",
		Format:         "json",
		Qzreferrer:     userQzoneURL + "/" + strconv.FormatInt(m.QQ, 10),
	}
	if len(base64imgList) > 0 {
		epr.Richtype = "1"
		epr.Richval = strings.Join(richvalList, "\t")
		epr.Subrichtype = "1"
		epr.PicBo = strings.Join(picBoList, ",")
	}
	url := fmt.Sprintf(emotionPublishURL, m.Gtk2)
	payload := strings.NewReader(structToStr(epr))
	data, err := DialRequest(NewRequest(WithMethod("POST"), WithUrl(url), WithBody(payload),
		WithHeader(map[string]string{
			"referer": userQzoneURL,
			"origin":  userQzoneURL,
			"cookie":  m.Cookie,
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

// ShuoShuoList 获取所有说说 实际能访问的说说个数 <= 说说总数(空间仅展示近半年等情况)
func (m *Manager) ShuoShuoList(uin int64, num int64, ms int64) (ShuoShuo []*models.ShuoShuoResp, err error) {
	cnt := num
	t := int(math.Ceil(float64(cnt) / 20.0))
	var i int
	for range t {
		ShuoShuoTemp, err := m.shuoShuoListRaw(uin, 20, i, 0)
		if err != nil {
			log.Println("所有说说获取失败:", err.Error())
			return nil, err
		}
		ShuoShuo = append(ShuoShuo, ShuoShuoTemp...)
		i = i + 20
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
	return ShuoShuo, nil
}

// GetShuoShuoCount 获取用户QQ号为uin的说说总数
func (m *Manager) GetShuoShuoCount(uin int64) (cnt int64, err error) {
	mlr := models.MsgListRequest{
		Uin:                uin,
		Ftype:              "0",
		Sort:               "0",
		Pos:                "0",
		Num:                "1",
		Replynum:           "0",
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

// GetLevel1CommentCount 获取一级评论总数
func (m *Manager) GetLevel1CommentCount(tid string) (cnt int64, err error) {
	url := fmt.Sprintf(getCommentsURL, strconv.FormatInt(m.QQ, 10), 1, 1, tid, m.Gtk2)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"cookie": m.Cookie,
	})))
	if err != nil {
		er := errors.New("说说评论请求错误:" + err.Error())
		log.Println("说说评论请求失败:", er.Error())
		return -1, er
	}
	r := cRe.FindStringSubmatch(string(data))
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

// ShuoShuoCommentList 根据说说ID获取评论，仅限本用户（m.QQ） // TODO: 待测试
func (m *Manager) ShuoShuoCommentList(tid string, num int64, ms int64) (comments []*models.Comment, err error) {
	numOfComments := num
	t := int(math.Ceil(float64(numOfComments) / 20.0))
	var i int
	for range t {
		commentsTemp, err := m.shuoShuoCommentsRaw(20, i, tid)
		if err != nil {
			return nil, err
		}
		comments = append(comments, commentsTemp...)
		i = i + 20
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
	return comments, nil
}

// GetLatestShuoShuo 获取用户QQ号为uin的最新说说
func (m *Manager) GetLatestShuoShuo(uin int64) (*models.ShuoShuoResp, error) {
	ss, err := m.shuoShuoListRaw(uin, 1, 0, 0)
	fmt.Println("ss!!!!!!:", ss)
	if err != nil {
		er := errors.New("最新说说获取错误:" + err.Error())
		log.Println("最新说说获取失败:", er.Error())
		return nil, er
	}
	return ss[0], nil
}

// DoLike 说说空间点赞 TODO:疑似无效
func (m *Manager) DoLike(tid string) (*models.LikeResp, error) {
	lr := models.LikeRequest{
		Qzreferrer: userQzoneURL + strconv.FormatInt(m.QQ, 10),
		Opuin:      m.QQ,
		Unikey:     userQzoneURL + strconv.FormatInt(m.QQ, 10) + "/mood/" + tid,
		From:       "1",
		Fid:        tid,
		Typeid:     "0",
		Appid:      "311",
	}
	lr.Curkey = lr.Unikey
	url := fmt.Sprintf(likeURL, m.Gtk2)
	payload := strings.NewReader(structToStr(lr))
	data, err := DialRequest(NewRequest(WithMethod("POST"), WithUrl(url),
		WithBody(payload), WithHeader(map[string]string{
			"referer": userQzoneURL,
			"origin":  userQzoneURL,
			"cookie":  m.Cookie,
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
