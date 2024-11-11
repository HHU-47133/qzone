package examples

import (
	"github.com/HHU-47133/qzone"
	"testing"
)

// 测试好友相关
func TestFriendList(t *testing.T) {
	// 从cookie创建空间对象
	qm := qzone.NewQZone().WithCookie(cookie)
	friends, _ := qm.FriendList()
	for i := 0; i < 10; i++ {
		t.Log("[好友简略信息]", friends[i].Name, friends[i].Uin, friends[i].Online, friends[i].Image, friends[i].GroupName)
		//fid, _ := m.FriendInfoDetail(friends[i].Uin) //TODO:详细信息获取有时候会莫名报错可能需要代理IP
		//t.Log("[好友详细信息]", fid.Uin, fid.Age, fid.Nickname, fid.Sex, fid.Birthyear, fid.Birthday, fid.Country, fid.Province, fid.City, fid.Mailname, fid.Mailcellphone, fid.Avatar, fid.Signature)
	}
	// 保存第一个好友的QQ, 方便后续操作
	friendID = friends[0].Uin
}

// 测试QQ群列表
func TestQQGroupList(t *testing.T) {
	groups, _ := qzone.NewQZone().WithCookie(cookie).QQGroupList()
	for _, v := range groups {
		t.Log("[QQ群信息]", v.GroupCode, v.GroupName, v.TotalMember, v.NotFriends)
	}
	groupID = groups[0].GroupCode
}

// 测试QQ群友列表
func TestGroupMemberList(t *testing.T) {
	groupMember, _ := qzone.NewQZone().WithCookie(cookie).QQGroupMemberList(groupID)
	for _, v := range groupMember {
		t.Log("[QQ群非好友信息]", v.Uin, v.NickName, v.AvatarURL, v.GroupCode)
	}
}
