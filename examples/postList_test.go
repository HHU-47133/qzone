package examples

import (
	"fmt"
	"github.com/HHU-47133/qzone"
	"testing"
	"time"
)

var (
	// cookie 登录成功后的 cookie
	cookie = "qrsig=cfe103dff2139455380c239f3df6117b1c68a1212960259d1fd18aa3a8fd34d42d43a65631ac2a5339c543883ad39a8fa649698db59caeff;uin=o1294222408;skey=@zHCefuSGJ;pt2gguin=o1294222408;p_uin=o1294222408;pt4_token=Wn6iBid5eZ42zxS44m9273hvwlkE*A19zIJ*3K7MmrM_;p_skey=SJZj0twlEHOfmhNBVAywpWWWbp6WLaAEcl*mw3N1K2w_;"
)

// 获取所有的说说
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

// 获取说说所有的一级评论
func TestGetComments(t *testing.T) {
	m := qzone.NewManager(cookie)
	comments, err := m.GetShuoShuoComments("4844244d9011f866f3d90500")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("🧡🧡🧡评论结构体🧡🧡🧡：")
	for _, comment := range comments {
		fmt.Printf("%+v\n", comment)

	}

}
