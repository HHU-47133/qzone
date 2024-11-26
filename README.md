# qzone 

> 提供qq空间基础功能接口

**！！本项目尚未开发完毕,改动较大！！**

开发进度
- [x] 基础接口封装
- [x] 扫码登录
- [ ] 规范接口返回字段
- [ ] 接口的统一分页设计
- [ ] 对低级别接口进一步封装，实现便捷功能
- [ ] 探索“与我相关”推送接口，解析历史数据
- [ ] 探索账号密码、便捷登录


- 导入项目
```go
go get -u github.com/HHU-47133/qzone
```
## 功能接口
- 具体实现请参看 `examples/*_test.go`
- 管理类实现 `manager.go`; 接口实现 `api.go`
### 登录流程
```go
// 1、创建对象
qm := qzone.NewQZone()
```
```go
// 2、获取二维码
// 成功返回"base64编码的二维码数据"
b64s, err := qm.GenerateQRCode()
```
```go
// 3、检测二维码扫码状态
// 0成功 1未扫描 2未确认 3已过期 -1系统错误
status, err := qm.CheckQRCodeStatus()

// 成功登录后qm对象会暴露公开字段qm.Info
type info struct {
    QQ          string // QQ空间的账号
    Cookie      string // 登录成功的Cookie，保存以便下次使用
    ExpiredTime time.Time
}

// 保存cookie方便下次创建对象
cookie := qm.Info.Cookie
```
- 从cookie创建
```go
// 你可以直接通过cookie创建一个空间操作对象
// cookie可以从扫码登录成功后qm.Info.Cookie获取
qm := qzone.NewQZone().WithCookie(cookie)
```
### 好友、群相关
- 群列表获取
```go
func (q *QZone) QQGroupList() ([]*models.QQGroupResp, error)
```
- 好友获取
```go
func (q *QZone) FriendList() ([]*models.FriendInfoEasyResp, error)
```
- 群友(非好友)获取
```go
func (q *QZone) QQGroupMemberList(gid int64) ([]*models.QQGroupMemberResp, error)
```
- 好友详细信息获取
```go
// uin:本人QQ
func (q *QZone) FriendInfoDetail(uin int64) (*models.FriendInfoDetailResp, error)
```
### 说说相关
- 说说发布
```go
// content:文本内容
// base64imgList:图片数组,为nil则只发文字
func (q *QZone) PublishShuoShuo(content string, base64imgList []string) (*models.ShuoShuoPublishResp, error)
```
- 说说获取
```go
// uin:有访问权限的QQ
// num:获取说说个数
// ms:延迟访问毫秒
func (q *QZone) ShuoShuoList(uin int64, num int64, ms int64) (ShuoShuo []*models.ShuoShuoResp, err error)
```
- 说说总数获取
```go
// uin:有访问权限的QQ
// 实际能访问的说说数量<=说说总数(封存动态)
func (q *QZone) GetShuoShuoCount(uin int64) (cnt int64, err error)
```
- 说说一级评论总数
```go
// tid:说说id（限制本人）
func (q *QZone) GetLevel1CommentCount(tid string) (cnt int64, err error)
```
- 说说评论内容获取
```go
// tid:说说id（限制本人）
// num:评论上限
// ms:延迟访问毫秒
func (q *QZone) ShuoShuoCommentList(tid string, num int64, ms int64) 
```
- 最新说说获取
```go
// uin:有访问权限的QQ
func (q *QZone) GetLatestShuoShuo(uin int64) (*models.ShuoShuoResp, error)
```

- 历史消息数据获取
```go
// GetQZoneHistory 获取QQ空间历史消息（限制本人）
func (q *QZone) GetQZoneHistory() ([]*models.QZoneHistoryItem, error)
````

### 其他
- 单个说说地址
```go
"https://user.qzone.qq.com/"+QQ号+"/mood/"+说说tid
```


### model 

- 请求响应结构，简洁信息参考 `model.go` 文件，详细信息参考 `types.go` 文件