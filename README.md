# qzone 

> 提供qq空间基础功能接口

- 导入项目
```go
go get -u github.com/HHU-47133/qzone
```
## 功能接口
- 具体实现请参看 `examples/*_test.go`
- 管理类实现 `manager.go`; 接口实现 `api.go`
### 登录
```go
// qrCodeOutputPath:二维码输出路径,例："./1.png"
// qrCodeInBytes:二维码字节流输出通道,向有缓冲区的通道发送最新字节流数据
// retryNum:尝试扫码登录的最大次数
func QzoneLogin(qrCodeOutputPath string, qrCodeInBytes chan []byte, retryNum int64) (m Manager, err error)
```
### 好友、群相关
- 群列表获取
```go
func (m *Manager) QQGroupList() ([]*models.QQGroupResp, error)
```
- 好友获取
```go
func (m *Manager) FriendList() ([]*models.FriendInfoEasyResp, error)
```
- 群友(非好友)获取
```go
func (m *Manager) QQGroupMemberList(gid int64) ([]*models.QQGroupMemberResp, error)
```
- 好友详细信息获取
```go
// uin:本人QQ
func (m *Manager) FriendInfoDetail(uin int64) (*models.FriendInfoDetailResp, error)
```
### 说说相关
- 说说发布
```go
// content:文本内容
// base64imgList:图片数组,为nil则只发文字
func (m *Manager) PublishShuoShuo(content string, base64imgList []string) (*models.ShuoShuoPublishResp, error)
```
- 说说获取
```go
// uin:有访问权限的QQ
// num:获取说说个数
// ms:延迟访问毫秒
func (m *Manager) ShuoShuoList(uin int64, num int64, ms int64) (ShuoShuo []*models.ShuoShuoResp, err error)
```
- 说说总数获取
```go
// uin:有访问权限的QQ
// 实际能访问的说说数量<=说说总数(封存动态)
func (m *Manager) GetShuoShuoCount(uin int64) (cnt int64, err error)
```
- 说说一级评论总数
```go
// tid:说说id（限制本人）
func (m *Manager) GetLevel1CommentCount(tid string) (cnt int64, err error)
```
- 说说评论内容获取
```go
// tid:说说id（限制本人）
// num:评论上限
// ms:延迟访问毫秒
func (m *Manager) ShuoShuoCommentList(tid string, num int64, ms int64) 
```
- 最新说说获取
```go
// uin:有访问权限的QQ
func (m *Manager) GetLatestShuoShuo(uin int64) (*models.ShuoShuoResp, error)
```
### 其他
- 单个说说地址
```go
"https://user.qzone.qq.com/"+QQ号+"/mood/"+说说tid
```


### model 

- 请求响应结构，简洁信息参考 `model.go` 文件，详细信息参考 `types.go` 文件