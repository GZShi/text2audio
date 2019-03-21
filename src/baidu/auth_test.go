package baidu

import (
	"fmt"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	key := "pZXEHHmetABCd1x2TmLdFrXL"
	secret := "CbKrOOhrR05ohCjIzyK1xaqpb9ET69Q3"
	client := NewAuthClient(key, secret, "test")

	go client.KeepFresh()
	<-time.After(time.Second * 1)

	token, cuid, err := client.GetToken()
	if err != nil {
		t.Error(err)
		return
	}
	if cuid != "test" {
		t.Error("not equal")
		return
	}

	fmt.Println("AccessToken", token)
	if token == "" {
		t.Error("token is empty")
		return
	}
	client.StopFresh()
}
