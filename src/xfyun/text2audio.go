package xfyun

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/GZShi/text2speech/src/utils"
)

type xparam struct {
	AUE        string `json:"aue"`
	AUF        string `json:"auf"`
	VoiceName  string `json:"voice_name"`
	Speed      string `json:"speed"`
	Volume     string `json:"volume"`
	Pitch      string `json:"pitch"`
	EngineType string `json:"engine_type"`
	TextType   string `json:"text_type"`
}

type Text2Audio struct {
	xparam
	AppID  string
	AppKey string
}

func NewText2Audio(appID, appKey string) *Text2Audio {
	return &Text2Audio{
		xparam: xparam{
			AUE:        "lame",
			AUF:        "audio/L16;rate=16000",
			VoiceName:  "xiaoyan",
			Speed:      "50",
			Volume:     "50",
			Pitch:      "50",
			EngineType: "intp65",
			TextType:   "text",
		},
		AppID:  appID,
		AppKey: appKey,
	}
}

func (t *Text2Audio) Get(text string) (audioData []byte, err error) {
	// 构造请求Param
	param, err := json.Marshal(t.xparam)
	if err != nil {
		return nil, err
	}
	base64Param := base64.StdEncoding.EncodeToString(param)

	// 计算校验和
	currTime := strconv.FormatInt(time.Now().Unix(), 10)
	w := md5.New()
	io.WriteString(w, t.AppKey+currTime+base64Param)
	checksum := hex.EncodeToString(w.Sum(nil))

	// 构造请求formdata
	body := "text=" + url.QueryEscape(text)

	res, data, err := utils.Post(
		"http://api.xfyun.cn/v1/service/v1/tts",
		[]byte(body),
		&map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"X-CurTime":    currTime,
			"X-Appid":      t.AppID,
			"X-Param":      base64Param,
			"X-CheckSum":   checksum,
		},
	)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(res.Header.Get("Content-Type"), "audio") {
		return data, nil
	}
	return nil, errors.New(string(data))
}
