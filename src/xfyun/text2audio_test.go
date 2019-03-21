package xfyun

import (
	"io/ioutil"
	"testing"
)

func TestText2Audio(t *testing.T) {
	appID := "5c937cd5"
	appKey := "925f72ad4baa7a39bc52def6c1b77d64"

	client := NewText2Audio(appID, appKey)

	data, err := client.Get("如果炸弹人的技能伤害是物理的话，那么战士流派是最好的应对选择。")
	if err != nil {
		t.Error(err)
		return
	}

	err = ioutil.WriteFile("testdata/test.mp3", data, 0644)
	if err != nil {
		t.Error(err)
		return
	}
}
