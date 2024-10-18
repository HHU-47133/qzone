package qzone

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// genderGTK 生成GTK
func genderGTK(sKey string, hash int) string {
	for _, s := range sKey {
		us, _ := strconv.Atoi(fmt.Sprintf("%d", s))
		hash += (hash << 5) + us
	}
	return fmt.Sprintf("%d", hash&0x7fffffff)
}

func structToStr(in interface{}) (payload string) {
	keys := make([]string, 0, 16)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		get := field.Tag.Get("json")
		if get != "" {
			var t string
			if v.Field(i).Kind() == reflect.Int64 {
				t = strconv.FormatInt(v.Field(i).Int(), 10)
			} else {
				t = v.Field(i).Interface().(string)
			}

			keys = append(keys, get+"="+url.QueryEscape(t))
		}
	}
	payload = strings.Join(keys, "&")
	return
}

// 获取说说详情页面
func getShuoShuoUnikey(uin string, tid string) (unikey string) {
	return fmt.Sprintf("http://user.qzone.qq.com/%s/mood/%s", uin, tid)
}
