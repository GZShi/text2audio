package baidu

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/GZShi/text2speech/src/utils"
	"github.com/google/go-querystring/query"
)

// Text2Audio 结构
type Text2Audio struct {
	authClient *AuthClient
	// 必填
	ClientType int    `url:"ctp"`
	Language   string `url:"lan"`

	// Speed 选填 语速，取值0-15，默认为5中语速
	Speed int `url:"spd"`
	// Pitch 选填 音调，取值0-15，默认为5中语速
	Pitch int `url:"pit"`
	// Volume 选填 音调，取值0-15，默认为5中语速
	Volume int `url:"vol"`
	// Person 选填 发声人，0普通女声，1普通男声，3情感合成-度逍遥，4情感合成-度丫丫
	Person int `url:"per"`
	// AudioUE 选填 格式，3-mp3，4-pcm16k，5-pcm-8k，6-wav
	AudioUE int `url:"aue"`
}

// NewText2Audio 创建
func NewText2Audio(client *AuthClient) *Text2Audio {
	return &Text2Audio{
		authClient: client,
		ClientType: 1,
		Language:   "zh",
		Speed:      4,
		Pitch:      5,
		Volume:     5,
		Person:     1,
		AudioUE:    3,
	}
}

// Get 将段文本转化为语音文件
func (t *Text2Audio) Get(text string) (audioData []byte, err error) {
	token, cuid, err := t.authClient.GetToken()
	if err != nil {
		return nil, err
	}

	encodedText := url.QueryEscape(text)
	encodedToken := url.QueryEscape(token)
	encodedCUID := url.QueryEscape(cuid)

	q, err := query.Values(t)
	if err != nil {
		return nil, err
	}
	body := fmt.Sprintf(
		"tex=%s&tok=%s&cuid=%s&%s",
		encodedText,
		encodedToken,
		encodedCUID,
		q.Encode(),
	)

	res, data, err := utils.Post(
		"https://tsn.baidu.com/text2audio",
		[]byte(body),
		&map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	)
	if err != nil {
		return nil, err
	}

	contentType := res.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "audio") {
		return data, nil
	}

	return nil, errors.New(string(data))
}
