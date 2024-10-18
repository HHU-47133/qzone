package qzone

import (
	"errors"
	"fmt"
	"github.com/HHU-47133/qzone/models"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
	"strings"
	"time"
)

// GetShuoShuoCommentsRaw 从第pos条评论开始获取num条评论，num最大为20
func (m *Manager) shuoShuoCommentsRaw(num int, pos int, tid string) (comments []*models.Comment, err error) {
	url := fmt.Sprintf(getCommentsURL, strconv.FormatInt(m.QQ, 10), pos, num, tid, m.Gtk2)
	data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
		"cookie": m.Cookie,
	})))
	if err != nil {
		er := errors.New("说说评论列表请求错误:" + err.Error())
		log.Println("说说评论列表获取失败:", er.Error())
		return nil, err
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		er := errors.New("说说评论正则解析错误:" + err.Error())
		log.Println("说说评论列表获取失败:", er.Error())
		return nil, er
	}
	jsonRaw := r[1]

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

// UploadImage 上传图片
func (m *Manager) uploadImage(base64img string) (*models.UploadImageResp, error) {
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
		Qzreferrer:    userQzoneURL + "/" + strconv.FormatInt(m.QQ, 10),
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
		return nil, errors.New("图片上传请求错误:" + err.Error())
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		return nil, errors.New("图片上传响应解析错误:" + string(data))
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

// getPicBoAndRichval 获取已上传图片重要信息
func (m *Manager) getPicBoAndRichval(data *models.UploadImageResp) (picBo, richval string, err error) {
	var flag bool
	if data.Ret != 0 {
		err = errors.New("已上传图片信息错误:fuck")
		return
	}
	_, picBo, flag = strings.Cut(data.URL, "&bo=")
	if !flag {
		err = errors.New("已上传图片URL错误:" + data.URL)
		return
	}
	richval = fmt.Sprintf(",%s,%s,%s,%d,%d,%d,,%d,%d", data.Albumid, data.Lloc, data.Sloc, data.Type, data.Height, data.Width, data.Height, data.Width)
	return
}

// ShuoShuoListRaw 获取用户QQ号为uin且最多num个说说列表，每个说说获取上限replynum个评论数量 TODO:replynum无效果，说说评论展示条数和num绑定
func (m *Manager) shuoShuoListRaw(uin int64, num int, pos int, replynum int) ([]*models.ShuoShuoResp, error) {
	mlr := models.MsgListRequest{
		Uin:                uin,
		Ftype:              "0",
		Sort:               "0",
		Pos:                strconv.Itoa(pos),
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
		er := errors.New("说说列表请求错误:" + err.Error())
		log.Println("说说列表获取失败:", er.Error())
		return nil, er
	}
	jsonStr := string(data)
	// 判断是否有访问权限
	forbid := gjson.Get(jsonStr, "message").String()
	if forbid != "" {
		er := errors.New("说说列表解析错误:" + forbid)
		log.Println("说说列表获取失败:", er.Error())
		return nil, er
	}

	var resLen int64
	if !gjson.Get(jsonStr, "msglist.#").Exists() {
		er := errors.New("说说列表解析错误:" + jsonStr)
		log.Println("说说列表获取失败:", er.Error())
		return nil, er
	}
	resLen = gjson.Get(jsonStr, "msglist.#").Int()
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
			PicTotal:    shuoshuo.Get("pictotal").Int(),
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
