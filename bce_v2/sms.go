package bce

import (
	"github.com/yangfanlalala/blx/crypto"
	"github.com/yangfanlalala/blx/network"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	SmsEndpoint   = "sms.bj.baidubce.com"
	SmsApiSend    = "/bce/v2/message"
	SmsExpiration = 3600
)

type SmsContentVar map[string]string

type SmsPayload struct {
	InvokeId     string        `json:"invokeId"`
	PhoneNumber  string        `json:"phoneNumber"`
	TemplateCode string        `json:"templateCode"`
	ContentVar   SmsContentVar `json:"contentVar"`
}

type smsSendResponse struct {
	RequestId string
	Code      string
	Message   string
}

type smsClient struct {
	Credentials *Credentials
}

func NewSmsClient(credentials *Credentials) *smsClient {
	return &smsClient{Credentials: credentials}
}

func (cli *smsClient) SendMessage(payload *SmsPayload) (err error) {
	reqs, _ := json.Marshal(payload)
	req, err := network.NewHttpRequest(network.HttpMethodPost, "https://"+SmsEndpoint+SmsApiSend, reqs)
	if err != nil {
		return
	}
	fmt.Println(string(reqs))
	req.Header.Set("Content-Type", network.ContentTypeJSON)
	req.Header.Set("x-bce-content-sha256", crypto.Sha256(string(reqs)))
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(reqs)))
	signature := new(AuthSignerV1).Sign(req, cli.Credentials, &SignOptionsV1{Timestamp: time.Now().UTC(), ExpireSeconds: SmsExpiration})
	req.Header.Set("Authorization", signature)
	fmt.Println(signature)
	client := network.NewHttpClient()
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("http requests is not ok, status code[%d] message[%s]", res.StatusCode, body)
		return
	}
	sms := &smsSendResponse{}
	if err = json.Unmarshal(body, sms); err != nil {
		return
	}
	if sms.Code != "100" {
		err = fmt.Errorf("send sms is not ok, code[%s] message[%s] request_id[%s]", sms.Code, sms.Message, sms.RequestId)
		return
	}
	return
}
