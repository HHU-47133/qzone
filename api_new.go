package qzone

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"math"
	"net/http"
	"time"
)

const (
	getCommentsURL = "https://h5.qzone.qq.com/proxy/domain/taotao.qq.com/cgi-bin/emotion_cgi_msgdetail_v6?uin=%s&pos=%d&num=%d&tid=%s&format=jsonp&g_tk=%s"
)

// 暂时没用到
type CommentRaw struct {
	IsPasswordLuckyMoneyCmtRight string    `json:"IsPasswordLuckyMoneyCmtRight"`
	Abledel                      int       `json:"abledel"`
	Content                      string    `json:"content"` //评论内容，为空则为图片评论，Pic
	CreateTime                   string    `json:"createTime"`
	CreateTime2                  string    `json:"createTime2"`
	CreateTime0                  time.Time `json:"create_time"` //发送时间戳
	List3                        []struct {
		Abledel     int    `json:"abledel"`
		Content     string `json:"content"`
		CreateTime  string `json:"createTime"`
		CreateTime2 string `json:"createTime2"`
		CreateTime0 int    `json:"create_time"`
		Name        string `json:"name"`
		Owner       struct {
			Name string `json:"name"`
			Uin  int    `json:"uin"`
		} `json:"owner"`
		SourceName string `json:"source_name"`
		SourceURL  string `json:"source_url"`
		T3Source   int    `json:"t3_source"`
		Tid        int    `json:"tid"`
		Uin        int    `json:"uin"`
	} `json:"list_3"` //二级评论
	Name  string `json:"name"`
	Owner struct {
		Name string `json:"name"` //评论发送人的昵称
		Uin  int64  `json:"uin"`  //评论发送人的QQ
	} `json:"owner"`
	Pic []struct {
		BHeight  int    `json:"b_height"`
		BURL     string `json:"b_url"`
		BWidth   int    `json:"b_width"`
		HdHeight int    `json:"hd_height"`
		HdURL    string `json:"hd_url"`
		HdWidth  int    `json:"hd_width"`
		OURL     string `json:"o_url"`
		SHeight  int    `json:"s_height"`
		SURL     string `json:"s_url"`
		SWidth   int    `json:"s_width"`
		Who      int    `json:"who"`
	} `json:"pic"` //图片评论细节信息，一般不使用
	Pictotal int `json:"pictotal"` //图片总数
	Private  int `json:"private"`
	ReplyNum int `json:"replyNum"` //二级评论数
	RichInfo []struct {
		Burl string `json:"burl"` //评论图片链接地址
		Type int    `json:"type"`
		Who  int    `json:"who"`
	} `json:"rich_info"` //评论图片
	SourceName string      `json:"source_name"`
	SourceURL  string      `json:"source_url"`
	T2Source   int         `json:"t2_source"`
	T2Subtype  int         `json:"t2_subtype"`
	T2Termtype int         `json:"t2_termtype"`
	T2WcNick   interface{} `json:"t2_wc_nick"`
	T3Subtype  int         `json:"t3_subtype"`
	T3Termtype int         `json:"t3_termtype"`
	Tid        int         `json:"tid"`
	Uin        int         `json:"uin"`
}

// Comment 评论简单结构体，目前支持一级评论
type Comment struct {
	ShuoShuoID string    //当前评论所属的说说ID
	OwnerName  string    //当前评论人的昵称
	OwnerUin   int64     //当前评论人的QQ
	Content    string    //评论内容，为空则是图片评论
	PicContent []string  //图片评论链接
	CreateTime time.Time //发布评论的时间戳
}

// 根据说说ID获取所有评论
func (m *Manager) GetShuoShuoComments(tid string) (comments []Comment, err error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf(getCommentsURL, m.QQ, 1, 1, tid, m.Gtk2), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("cookie", m.Cookie)
	request.Header.Add("user-agent", ua)
	request.Header.Add("content-type", contentType)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	jsonRaw := cRe.FindStringSubmatch(string(data))[1]
	// 说说的一级评论总数
	numOfComments := gjson.Get(jsonRaw, "cmtnum").Int()
	t := int(math.Ceil(float64(numOfComments) / 20.0))
	fmt.Println("我是t", t, "  ", numOfComments)
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

// 从第pos条评论开始获取num条评论，num最大为20
func (m *Manager) getShuoShuoCommentsRaw(num int, pos int, tid string) (comments []Comment, err error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf(getCommentsURL, m.QQ, pos, num, tid, m.Gtk2), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("cookie", m.Cookie)
	request.Header.Add("user-agent", ua)
	request.Header.Add("content-type", contentType)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	jsonRaw := cRe.FindStringSubmatch(string(data))[1]
	//fmt.Println("🧡🧡🧡取说说评论测试🧡🧡🧡：", jsonRaw)

	// 取出评论数据
	commentJsonList := gjson.Get(jsonRaw, "commentlist").Array()
	var comment Comment
	for _, com := range commentJsonList {
		comment = Comment{
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
