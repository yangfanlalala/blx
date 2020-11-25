package bce

import (
	"github.com/yangfanlalala/blx/crypto"
	blxHttp "github.com/yangfanlalala/blx/network"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	SMS_HOST      = "sms.bj.baidubce.com"
	SMS_URI_SEND  = "/bce/v2/message"
	SMS_URI_FETCH = "/v1/message/{messageId}"
	SMS_TPL_1     = ""
	SMS_TPL_2     = ""
	SMS_EXPIRE    = 500000
)

type SMSBase struct {
}

type SMSClient struct {
	TplID string
	Cred  Credentials
}

func NewSmsClient(payload string, tplID string) *SMSClient {
	return &SMSClient{}
}

func (c *SMSClient) BuildRequest(method, url, payload string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", blxHttp.ContentTypeJSON)
	req.Header.Add("x-bce-content-sha256", crypto.Sha256(payload))
	now := time.Now()
	req.Header.Add("x-bce-date", now.Format("2006-01-02"))
	req.Header.Add("Host", req.Host)
	signOpt := &SignOptions{
		Timestamp:     now.UTC().Unix(),
		ExpireSeconds: SMS_EXPIRE,
	}
	signer := new(AuthSigner)
	signature := signer.Sign(req, &c.Cred, signOpt)
	fmt.Println(signature)
	req.Header.Add("Authorization", signature)
	return req, nil
}

func (c *SMSClient) SendSMS(payload map[string]interface{}) (*http.Response, error) {
	httpClient := &blxHttp.Client{}
	body, _ := json.Marshal(payload)
	req, _ := c.BuildRequest(blxHttp.METHOD_POST, SMS_HOST+SMS_URI_SEND, string(body))
	return httpClient.Do(req)
}
