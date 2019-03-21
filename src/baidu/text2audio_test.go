package baidu

import (
	"io/ioutil"
	"testing"
)

func TestText2Audio(t *testing.T) {
	key := "pZXEHHmetABCd1x2TmLdFrXL"
	secret := "CbKrOOhrR05ohCjIzyK1xaqpb9ET69Q3"
	client := NewAuthClient(key, secret, "test")
	client.KeepFresh()

	t2a := NewText2Audio(client)
	mp3Data, err := t2a.Get("如果炸弹人的技能伤害是物理的话，那么战士流派是最好的应对选择。")
	if err != nil {
		t.Error(err)
		return
	}

	err = ioutil.WriteFile("testdata/test.mp3", mp3Data, 0644)
	if err != nil {
		t.Error(err)
		return
	}

	client.StopFresh()
}
