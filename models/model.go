package models

import "time"

// ShuoShuoPublishResp 发布说说响应结构体
type ShuoShuoPublishResp struct {
	Tid      string // 说说Id
	Code     int64  // 响应状态码，0成功
	Now      int64  // 发布时间戳
	Feedinfo string // 说说页面html元素
	Message  string // ？错误后返回的消息
}

// QQGroupReq 获取QQ群请求结构体
type QQGroupReq struct {
	Uin     string `json:"uin"`
	Do      string `json:"do"`
	Rd      string `json:"rd"`
	Fupdate string `json:"fupdate"`
	Clean   string `json:"clean"`
	GTk     string `json:"g_tk"`
}

// QQGroupResp 获取QQ群响应结构体
type QQGroupResp struct {
	GroupCode   int64  `json:"groupcode"`    //群号
	GroupName   string `json:"groupname"`    //群名
	TotalMember int64  `json:"total_member"` //群人数
	NotFriends  int64  `json:"notfriends"`   //群里非好友人数
}

// ShuoShuoResp 说说响应结构体
type ShuoShuoResp struct {
	Uin         int64  // 用户QQ
	Name        string // 用户昵称
	Tid         string // 说说Id
	Content     string // 说说内容
	CreateTime  string // 说说创建时间
	CreatedTime int64  // 说说创建时间戳
	Pictotal    int64  // 图片总数
	Cmtnum      int64  // 评论数量
	Secret      int64  // 是否为私密动态
	Pic         []PicResp
}

// PicResp 说说响应结构体中的图片数据
type PicResp struct {
	PicId      string // 图片Id
	Url1       string // 原图更小
	Url2       string // 原图大小
	Url3       string // 原图指定hw
	Smallurl   string // 缩略图
	Curlikekey string // 链接
	Unilikekey string // 链接
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

// LikeResp 点赞响应结构体
type LikeResp struct {
	Ret int64
	Msg string
}

// UploadImageResp 上传图片响应结构体
type UploadImageResp struct {
	Pre        string // 低分辨率url
	URL        string // 完整url
	Width      int64  // 宽
	Height     int64  // 高
	OriginURL  string // 图片的原始url
	Contentlen int64  // 图片大小（字节）
	Albumid    string
	Lloc       string
	Sloc       string
	Type       int64
	Ret        int64
}

// FriendInfoEasyResp 好友简略信息响应结构体
type FriendInfoEasyResp struct {
	Uin       int64  // QQ号
	Groupid   int64  // 分组ID
	GroupName string // 分组名称
	Name      string // 名称
	Remark    string // 备注
	Image     string // 头像
	Online    int64  // 在线状态
}

// FriendInfoDetailResp 好友详细信息响应结构体
type FriendInfoDetailResp struct {
	Uin           int64  `json:"uin"`           // QQ号
	Nickname      string `json:"nickname"`      // 昵称
	Signature     string `json:"signature"`     // 签名
	Avatar        string `json:"avatar"`        // 上古头像
	Sex           int64  `json:"sex"`           // 性别，1男
	Age           int64  `json:"age"`           // 年龄
	Birthyear     int64  `json:"birthyear"`     // 生日年份
	Birthday      string `json:"birthday"`      // 生日月-天
	Country       string `json:"country"`       // 国家
	Province      string `json:"province"`      // 省份
	City          string `json:"city"`          // 城市
	Career        string `json:"career"`        // 职业
	Company       string `json:"company"`       // 公司
	Mailname      string `json:"mailname"`      // 邮件名称
	Mailcellphone string `json:"mailcellphone"` // 邮件绑定手机号
	Mailaddr      string `json:"mailaddr"`      // 邮件地址
}
