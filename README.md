# qzone 

> 提供qq空间基础功能接口


- 导入项目

```go
go get -u github.com/HHU-47133/qzone
```


## 功能接口

### 登录

- 详见 `examples/login_test.go`

```go
// 1. 获取二维码信息（data），取出cookie重要参数（qrsig、ptqrtoken）
data, qrsig, ptqrtoken, err = Ptqrshow()
// 2. 保存二维码
err = os.WriteFile("ptqrcode.png", data, 0666)
// 3. 查询登录回调，检测登录状态
for {
    data, ptqrloginCookie, err = Ptqrlogin(qrsig, ptqrtoken)
	...
	// 4. 成功登录后，获取登录重定向URL
    redirectCookie, err = LoginRedirect(redirectURL)
}
// 5. 创建信息管理结构，携带登录回调cookie和重定向页面cookie
m := NewManager(ptqrloginCookie + redirectCookie)
// 6. 执行其它接口操作
...
```

### 发布说说

```go
// PublishShuoShuo(content string, base64imgList []string)
// content：说说内容
// base64imgList：base64编码图片列表
ShuoShuoPublishResp, err = m.PublishShuoShuo("content", []string{picBase64})
```

### 获取用户所有说说

```go
// ShuoShuoList(uin string, num int, replynum int)
// uin：QQ号
// num：说说数量
// replynum：评论数量
ShuoShuoResp, err = m.ShuoShuoList(m.QQ, 20, 5)
```


### 获取说说所有一级评论

```go
// GetShuoShuoComments(tid string)
// tid：说说ID
comments, err := m.GetShuoShuoComments("4844244d9011f866f3d90500")
```

### 单个说说URL

```go
"https://user.qzone.qq.com/"+QQ号+"/mood/"+说说tid
```


### 上传图片

```go
// 1. 读取本地图片
srcByte, err = os.ReadFile(path)
// 2. base64编码
picBase64 = base64.StdEncoding.EncodeToString(srcByte)
// 3. 上传图片 
UploadImageResp, err = m.UploadImage(picBase64)
```


### 获取好友信息

- 接口限制，只能获取亲密度排序前200的好友
- 简略信息

```go
FriendInfoEasyResp, err = m.FriendList()
```


### 获取好友详细信息

```go
// FriendInfoDetail(uin int64)
// uin：好友QQ号
FriendInfoDetailResp, err = m.FriendInfoDetail(friend.uin)
```

### 说说点赞

```go
// DoLike(tid string)
// tid：说说ID
LikeResp, err = m.DoLike(shuoshuo.tid)
```


## model 

- 请求响应结构，简洁信息参考 `model.go` 文件，详细信息参考 `types.go` 文件