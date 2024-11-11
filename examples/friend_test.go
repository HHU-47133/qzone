package examples

import (
	"github.com/HHU-47133/qzone"
	"testing"
)

// 测试好友相关
func TestFriendList(t *testing.T) {
	friends, _ := qm.Qpack.FriendList()
	for i := 0; i < 10; i++ {
		t.Log("[好友简略信息]", friends[i].Name, friends[i].Uin, friends[i].Online, friends[i].Image, friends[i].GroupName)
		//fid, _ := m.FriendInfoDetail(friends[i].Uin) //TODO:详细信息获取有时候会莫名报错可能需要代理IP
		//t.Log("[好友详细信息]", fid.Uin, fid.Age, fid.Nickname, fid.Sex, fid.Birthyear, fid.Birthday, fid.Country, fid.Province, fid.City, fid.Mailname, fid.Mailcellphone, fid.Avatar, fid.Signature)
	}
	friendID = friends[0].Uin
}

// 测试QQ群列表
func TestQQGroupList(t *testing.T) {
	coo := "pt2gguin=o1294222408;uin=o1294222408;skey=@QXH2rd096;superuin=o1294222408;supertoken=3234127316;superkey=ieLifDO91hLNj3dHg2bUPm-gdNb1V5DffORaE9XM*-g_;pt_recent_uins=3d237f4a66eb24b516a34b203f0e956a9dc2bec812610943a0b6318d38fb5bbd405146b68274510809b3d702f3298c762f53283bf135e612;RK=SvFZBxECGe;ptnick_1294222408=52;ptcz=648a792fbf49676db46e3733151abe6a14ec23270984a118257e57d3e782bc04;uin=o1294222408;skey=@QXH2rd096;pt2gguin=o1294222408;p_uin=o1294222408;pt4_token=T*yo73*0ldXmkXgedKEPoJ82YNCIDku0XHvlEOpkYtM_;p_skey=MxjTxW9Im5N9ek2g*pOSIdY*Adlnw7gtJvUGqGLzp1Q_;"
	groups3, _ := qzone.NewQpack(coo).QQGroupList()
	for _, v := range groups3 {
		t.Log("[QQ群信息]", v.GroupCode, v.GroupName, v.TotalMember, v.NotFriends)
	}
	// 将groupQQ设置为第一个群组
	t.Log("第二次获取群")
	groups1, _ := qm.Qpack.QQGroupList()
	for _, v := range groups1 {
		t.Log("[QQ群信息]", v.GroupCode, v.GroupName, v.TotalMember, v.NotFriends)
	}
}

// 测试QQ群友列表
func TestGroupMemberList(t *testing.T) {
	groupMember, _ := qm.Qpack.QQGroupMemberList(groupID)
	for _, v := range groupMember {
		t.Log("[QQ群非好友信息]", v.Uin, v.NickName, v.AvatarURL, v.GroupCode)
	}
}
