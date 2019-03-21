package utils

import (
	"strings"
	"testing"

	"github.com/google/go-querystring/query"
)

func TestGetHTTP(t *testing.T) {
	_, data, err := Get("http://baidu.com", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}

	body := string(data)
	if !strings.HasPrefix(body, "<html>") {
		t.Error("is not html response")
		return
	}
}
func TestGetHTTPS(t *testing.T) {
	_, data, err := Get("https://txy.eval.im", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}

	body := string(data)
	if !strings.HasPrefix(body, "<!DOCTYPE html>") {
		t.Error("is not html response")
		return
	}
}

func TestPostHTTPS(t *testing.T) {
	q, err := query.Values(struct {
		FuncID  string `url:"funcid"`
		Body    string `url:"bodystr"`
		Timeout int    `url:"timeout"`
	}{
		"MysqlPro:GetPublishList",
		"IX,SPEC=2834,STRUCT=MysqlPro\r\n" +
			"1,i_group_type|2,i_group_id|\r\n" +
			"4|1|\r\n",
		20000,
	})
	if err != nil {
		t.Error(err)
		return
	}
	_, data, err := Post(
		"https://wjwj.newone.com.cn:9092/wj/Exec",
		[]byte(q.Encode()),
		&map[string]string{
			"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		},
	)

	body := string(data)
	if !strings.HasPrefix(body, "success") {
		t.Error("response is not success")
		return
	}
}
