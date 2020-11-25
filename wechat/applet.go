package wechat

import (
	blxHttp "blxee.com/utils/network"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const APPLET_HOST = "https://api.weixin.qq.com"
const APPLET_URI_MINI_CODE_TO_SESSION = "/sns/jscode2session"

type Applet struct {
	AppID  string
	Secret string
}

type AppletWatermark struct {
	AppID     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

type AppletError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type AppletSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	AppletError
}

type AppletPhoneNumber struct {
	PhoneNumber     string          `json:"phoneNumber"`
	PurePhoneNumber string          `json:"purePhoneNumber"`
	CountryCode     string          `json:"countryCode"`
	Watermark       AppletWatermark `json:"watermark"`
}

type AppletPhoneNumberRsp struct {
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
	AppletError
}

type AppletUserInfo struct {
	OpenID    string          `json:"openId"`
	NickName  string          `json:"nickName"`
	Gender    string          `json:"gender"`
	City      string          `json:"city"`
	Province  string          `json:"province"`
	Country   string          `json:"country"`
	AvatarURL string          `json:"avatarUrl"`
	UnionId   string          `json:"unionId"`
	Watermark AppletWatermark `json:"watermark"`
}

type AppletUserInfoRsp struct {
	UserInfo      AppletUserInfo `json:"userInfo"`
	RawData       string         `json:"rawData"`
	Signature     string         `json:"signature"`
	EncryptedData string         `json:"encryptedData"`
	Iv            string         `json:"iv"`
	AppletError
}

func NewApplet(appid, secret string) *Applet {
	return &Applet{appid, secret}
}

func (wx *Applet) GetSession(code string) (*AppletSession, error) {
	httpClient := blxHttp.NewHttpClient()
	httpURL := APPLET_HOST + APPLET_URI_MINI_CODE_TO_SESSION + "?appid=" + wx.AppID + "&secret=" + wx.Secret + "&js_code=" + code + "&grant_type=authorization_code"

	rsp, err := httpClient.Get(httpURL)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	body, _ := ioutil.ReadAll(rsp.Body)
	session := new(AppletSession)
	if err := json.Unmarshal(body, session); err != nil {
		return nil, err
	}
	if session.ErrCode != 0 {
		return nil, fmt.Errorf("requests wx srv error: code[%d] message[%s]", session.ErrCode, session.ErrMsg)
	}
	return session, nil
}

func AppletGetPhone(cipher, iv, session string) (*AppletPhoneNumber, error) {
	phone := new(AppletPhoneNumber)
	err := Decrypt(cipher, iv, session, phone)
	if err != nil {
		return nil, err
	}
	return phone, nil
}

func (wx *Applet) AppletGetUser() (*AppletUserInfo, error) {
	return &AppletUserInfo{}, nil
}
