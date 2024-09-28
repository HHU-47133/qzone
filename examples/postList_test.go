package examples

import (
	"fmt"
	"github.com/HHU-47133/qzone"
	"testing"
	"time"
)

var (
	// cookie 登录成功后的 cookie
	cookie = "qrsig=cfc0b8a7b126f764823f2d5880f56136369697ada9f9292253bb2fd97b853f2b2e7f8aeb750a2ebf71c7941fa2eaecfb495cbefea63bd3f5;uin=o1778046356;skey=@BLfHfFf0Z;pt2gguin=o1778046356;p_uin=o1778046356;pt4_token=2cTfl1WCSbYrjuVAmRsBtShdB7zT6DVMAXgLkQ1RkCE_;p_skey=*Gtk6DIts5*ISDq91t4LLWqwF-Ob6T4bF1aHRz17BRE_;"
)

func TestGetPostList(t *testing.T) {
	m := qzone.NewManager(cookie)
	list, err := m.EmotionMsglist("20", "1")
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range list.Msglist {
		fmt.Println(v.Name, v.Conlist, v.Tid, v.Pic, v.Cmtnum)
		for _, com := range v.Commentlist {
			fmt.Println("  ·", com.Content, com.Name, com.CreateTime, com.Uin)
		}
		if err := m.DoLike(v.Tid); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}
