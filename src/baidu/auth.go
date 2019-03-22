package baidu

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/GZShi/text2speech/src/utils"
)

type Auth struct {
	APIKey    string `url:"client_id"`
	SecretKey string `url:"client_secret"`
	GrantType string `url:"grant_type"`
}
type AccessToken struct {
	AccessToken   string `json:"access_token"`
	ExpiresIn     uint   `json:"expires_in"`
	RefreshToken  string `json:"refresh_token"`
	Scope         string `json:"scope"`
	SessionKey    string `json:"session_key"`
	SessionSecret string `json:"session_secret"`
}

type AuthClient struct {
	cuid        string
	auth        *Auth
	accessToken *AccessToken
	updateMut   sync.RWMutex
	expire      time.Time
	stopChan    chan int
}

func NewAuthClient(APIKey, SecretKey, cuid string) *AuthClient {
	return &AuthClient{
		cuid:        cuid,
		auth:        &Auth{APIKey, SecretKey, "client_credentials"},
		accessToken: nil,
		expire:      time.Now().Add(time.Second * time.Duration(-10)),
		stopChan:    make(chan int),
	}
}

// GetToken 获取最新的Token
func (c *AuthClient) GetToken() (string, string, error) {
	c.updateMut.RLock()
	defer c.updateMut.RUnlock()
	if time.Now().Sub(c.expire) >= 0 {
		return "", "", errors.New("AccessToken expired")
	}
	if c.accessToken == nil {
		return "", "", errors.New("accesstoken not inited")
	}
	return c.accessToken.AccessToken, c.cuid, nil
}

// KeepFresh 开始刷新
func (c *AuthClient) KeepFresh() {
	c.update()
	go func() {
		interval := time.Second * time.Duration(30)
		for {
			select {
			case sig := <-c.stopChan:
				log.Print("stop fresh with signal ", sig)
				return
			case <-time.After(interval):
				err := c.update()
				if err == nil {
					t := c.accessToken
					log.Printf("token updated token=%v expire=%v", t.AccessToken, t.ExpiresIn)
					if t.ExpiresIn != 0 {
						interval = time.Second * time.Duration(c.accessToken.ExpiresIn-60)
					}
				}
			}
		}
	}()
}

// StopFresh 停止刷新
func (c *AuthClient) StopFresh() {
	c.stopChan <- 0
}

func (c *AuthClient) update() error {
	_, data, err := utils.Get(
		"https://openapi.baidu.com/oauth/2.0/token",
		c.auth,
		nil,
	)

	accessToken := &AccessToken{}
	err = json.Unmarshal(data, accessToken)
	if err != nil {
		return err
	}

	c.updateMut.Lock()
	c.accessToken = accessToken
	c.expire = time.Now().Add(time.Second * time.Duration(accessToken.ExpiresIn))
	c.updateMut.Unlock()

	return nil
}
