package examples

import (
	"fmt"
	"github.com/HHU-47133/qzone"
	"strconv"
	"testing"
)

func TestGetFriendLists(t *testing.T) {
	m := qzone.NewManager(cookie)
	friends, err := m.FriendList()
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, v := range friends {
		fmt.Println("好友简略信息：", v.Name, v.Uin, v.Online, v.Image, v.GroupName)
		//fid, _ := m.FriendInfoDetail(v.Uin) TODO:详细信息获取有时候会莫名报错
		//fmt.Println("好友详细信息：", fid.Uin, fid.Age, fid.Nickname, fid.Sex, fid.Birthyear, fid.Birthday, fid.Country, fid.Province, fid.City, fid.Mailname, fid.Mailcellphone, fid.Avatar, fid.Signature)
	}
	friendQQ = strconv.FormatInt(friends[0].Uin, 10)
}

func TestQQGroupList(t *testing.T) {
	m := qzone.NewManager(cookie)
	groups, err := m.GetQQGroup()
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, v := range groups {
		fmt.Println("[QQ群信息]", v.GroupCode, v.GroupName, v.TotalMember, v.NotFriends)
	}
}
