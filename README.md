# qzone 

> 提供qq空间基础功能接口

**！！本项目尚未开发完毕,改动较大！！**

- 导入项目
```go
go get -u github.com/HHU-47133/qzone
```
## 功能接口
- 具体实现请参看 `examples/*_test.go`
- 管理类实现 `manager.go`; 接口实现 `api.go`
### 登录流程
- 结构介绍
```go
// QManager 用于管理多个QQ空间登录
type QManager struct {
    Mu    sync.RWMutex
    Store map[string]*qsession
}
```
```go
// qsession 用于管理QQ空间单次扫码登录
type qsession struct {
    UserID     string // 标识本次登录权限所有人
    QrCodeID   string // 标识本次二维码ID
    Qrsig      string // 扫码登录需要使用的参数
    Qrtoken    string // 由Qrsig计算而成
    Cookie     string // 扫码登录过程使用
    ExpiryTime time.Time // 本次登录二维码过期时间
    Qpack      *Qpack    // 成功登录后创建的空间对象
}
```
```go
// Qpack 空间操作对象，api都绑定在这里
type Qpack struct {
    Cookie string
    QQ     int64
    Gtk    string
    Gtk2   string
    PSkey  string
    Skey   string
    Uin    string
}
```
- 登录过程
```go
// 1、创建管理对象
qm := qzone.NewQManager()
```
```go
// 2、获取二维码
// 传入一个string类型的参数用于标识本次登录授权人
// 成功返回"base64编码的二维码数据","用于查询登录状态的二维码ID"
b64s, qrID, _ = qm.GenerateQRCode("test-uid")
```
```go
// 3、检测二维码扫码状态
// 0成功 1未扫描 2未确认 3已过期 -1系统错误
// 扫码成功后会自动创建Qpack对象:qm.Store[qrID].Qpack
status, err := qm.CheckQRCodeStatus(qrID, "test-uid")
```
- 其他
```go
// 你可以直接通过cookie创建一个空间操作对象
// cookie可以从qm.Store[qrID].Cookie获取
qp := qzone.NewQpack(cookie)
```
### 好友、群相关
- 群列表获取
```go
func (m *Qpack) QQGroupList() ([]*models.QQGroupResp, error)
```
- 好友获取
```go
func (m *Qpack) FriendList() ([]*models.FriendInfoEasyResp, error)
```
- 群友(非好友)获取
```go
func (m *Qpack) QQGroupMemberList(gid int64) ([]*models.QQGroupMemberResp, error)
```
- 好友详细信息获取
```go
// uin:本人QQ
func (m *Qpack) FriendInfoDetail(uin int64) (*models.FriendInfoDetailResp, error)
```
### 说说相关
- 说说发布
```go
// content:文本内容
// base64imgList:图片数组,为nil则只发文字
func (m *Qpack) PublishShuoShuo(content string, base64imgList []string) (*models.ShuoShuoPublishResp, error)
```
- 说说获取
```go
// uin:有访问权限的QQ
// num:获取说说个数
// ms:延迟访问毫秒
func (m *Qpack) ShuoShuoList(uin int64, num int64, ms int64) (ShuoShuo []*models.ShuoShuoResp, err error)
```
- 说说总数获取
```go
// uin:有访问权限的QQ
// 实际能访问的说说数量<=说说总数(封存动态)
func (m *Qpack) GetShuoShuoCount(uin int64) (cnt int64, err error)
```
- 说说一级评论总数
```go
// tid:说说id（限制本人）
func (m *Qpack) GetLevel1CommentCount(tid string) (cnt int64, err error)
```
- 说说评论内容获取
```go
// tid:说说id（限制本人）
// num:评论上限
// ms:延迟访问毫秒
func (m *Qpack) ShuoShuoCommentList(tid string, num int64, ms int64) 
```
- 最新说说获取
```go
// uin:有访问权限的QQ
func (m *Qpack) GetLatestShuoShuo(uin int64) (*models.ShuoShuoResp, error)
```
### 其他
- 单个说说地址
```go
"https://user.qzone.qq.com/"+QQ号+"/mood/"+说说tid
```


### model 

- 请求响应结构，简洁信息参考 `model.go` 文件，详细信息参考 `types.go` 文件