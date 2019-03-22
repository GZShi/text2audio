package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/GZShi/text2audio/src/baidu"
	"github.com/GZShi/text2audio/src/xfyun"
	"github.com/kataras/iris"
)

func response(ctx iris.Context, err error, data interface{}) {
	if err != nil {
		ctx.JSON(struct {
			ErrorCode int         `json:"errCode"`
			ErrorInfo string      `json:"errInfo"`
			Data      interface{} `json:"data"`
		}{-1, err.Error(), nil})
		return
	}
	ctx.JSON(struct {
		ErrorCode int         `json:"errCode"`
		ErrorInfo string      `json:"errInfo"`
		Data      interface{} `json:"data"`
	}{0, "", data})
}

func parseAuthStr(str string) ([]string, error) {
	strs := strings.Split(str, ":")
	if len(strs) != 2 {
		return nil, errors.New("bad auth string")
	}
	return strs, nil
}

func main() {
	var listenAddr string
	var baiduAuthRaw, xfyunAuthRaw string
	flag.StringVar(&listenAddr, "l", ":8081", "[host]:[port]")
	flag.StringVar(&baiduAuthRaw, "baidu", "", "[api key]:[secret key]")
	flag.StringVar(&xfyunAuthRaw, "xfyun", "", "[app id]:[app key]")
	flag.Parse()

	baiduAuth, err := parseAuthStr(baiduAuthRaw)
	if err != nil {
		panic(err)
	}

	xfyunAuth, err := parseAuthStr(xfyunAuthRaw)
	if err != nil {
		panic(err)
	}

	var filecacheMut sync.RWMutex
	filecache := make(map[string][]byte)

	// init baiduclient
	baiduclient := baidu.NewAuthClient(baiduAuth[0], baiduAuth[1], "serversideclient")
	baiduclient.KeepFresh()

	app := iris.New()

	// 创建语音文件，返回AudioFileToken
	app.Post("/tts/audio-file", func(ctx iris.Context) {
		var payload struct {
			API         string `json:"api"`
			ContentType string `json:"contentType"`
			Text        string `json:"text"`
			File        string `json:"file"`
			Xfyun       struct {
				VoiceName  string `json:"voiceName"`
				Speed      int64  `json:"speed"`
				Volumn     int64  `json:"volumn"`
				Pitch      int64  `json:"pitch"`
				EngineType string `json:"engineType"`
			} `json:"xfyun"`
			Baidu struct {
				Person int `json:"person"`
				Speed  int `json:"speed"`
				Volumn int `json:"volumn"`
				Pitch  int `json:"pitch"`
			} `json:"baidu"`
		}

		rawBody, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			response(ctx, err, nil)
			return
		}

		err = json.Unmarshal(rawBody, &payload)
		if err != nil {
			response(ctx, err, nil)
			return
		}

		// 把借口数据进行HASH，结果作为文件名缓存
		h := md5.New()
		h.Write(rawBody)
		tag := hex.EncodeToString(h.Sum(nil))

		filecacheMut.RLock()
		_, has := filecache[tag]
		filecacheMut.RUnlock()
		if has {
			response(ctx, nil, struct {
				Tag string `json:"tag"`
			}{tag})
			return
		}

		var mp3 []byte

		switch payload.API {
		case "xfyun":
			api := xfyun.NewText2Audio(xfyunAuth[0], xfyunAuth[1])
			api.VoiceName = payload.Xfyun.VoiceName
			api.Speed = strconv.FormatInt(payload.Xfyun.Speed, 10)
			api.Pitch = strconv.FormatInt(payload.Xfyun.Pitch, 10)
			api.Volume = strconv.FormatInt(payload.Xfyun.Volumn, 10)
			api.EngineType = payload.Xfyun.EngineType
			mp3, err = api.Get(payload.Text)
		case "baidu":
			api := baidu.NewText2Audio(baiduclient)
			api.Person = payload.Baidu.Person
			api.Speed = payload.Baidu.Speed
			api.Pitch = payload.Baidu.Pitch
			api.Volume = payload.Baidu.Volumn
			mp3, err = api.Get(payload.Text)
		default:
			response(ctx, fmt.Errorf("unknown api '%s'", payload.API), nil)
			return
		}

		if err != nil {
			response(ctx, err, nil)
			return
		}

		filecacheMut.Lock()
		filecache[tag] = mp3
		filecacheMut.Unlock()

		response(ctx, nil, struct {
			Tag string `json:"tag"`
		}{tag})
	})

	// 获取语音文件，可以下载和在线播放
	app.Get("/tts/audio-file/{tag:string}", func(ctx iris.Context) {
		tag := ctx.Params().Get("tag")

		filecacheMut.RLock()
		mp3, has := filecache[tag]
		filecacheMut.RUnlock()

		if !has {
			ctx.StatusCode(404)
			response(ctx, fmt.Errorf("'%s' not found", tag), nil)
			return
		}

		ctx.ContentType("audio/mp3")
		ctx.Header("Content-Disposition", `filename="tts.mp3"`)
		ctx.Write(mp3)
	})

	app.Run(iris.Addr(listenAddr))
}
