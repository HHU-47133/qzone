package examples

import (
	"fmt"
	"github.com/HHU-47133/qzone"
	"testing"
)

func TestGetFriendLists(t *testing.T) {
	m := qzone.NewManager(cookie)
	flv, err := m.FriendList()
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, v := range flv.FriendInfosEasy {
		fmt.Println(v)
	}
}
