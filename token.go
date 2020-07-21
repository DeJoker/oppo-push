package oppopush

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	DefaultTimeToLive = 3600 * 1000 // one hour
)

type OppoToken struct {
	AccessToken string `json:"access_token"`
	CreateTime  int64  `json:"create_time"`
}

//GetToken 获取AccessToken值
func (c *OppoPush) GetToken(appKey, masterSecret string) (*OppoToken, error) {
	nowMilliSecond := time.Now().UnixNano() / 1e6
	if (nowMilliSecond-c.TokenIns.CreateTime) < DefaultTimeToLive && c.TokenIns.AccessToken != "" {
		return c.TokenIns, nil
	}
	timestamp := strconv.FormatInt(nowMilliSecond, 10)
	shaByte := sha256.Sum256([]byte(appKey + timestamp + masterSecret))
	sign := fmt.Sprintf("%x", shaByte)
	params := url.Values{}
	params.Add("app_key", appKey)
	params.Add("sign", sign)
	params.Add("timestamp", timestamp)

	res, err := doPost(PushHost+AuthURL, params)
	if err != nil {
		return nil, err
	}

	var result AuthSendResult
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, errors.New(result.Message)
	}
	c.TokenIns.AccessToken = result.Data.AuthToken
	c.TokenIns.CreateTime = result.Data.CreateTime
	return c.TokenIns, nil
}
