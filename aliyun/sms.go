package aliyun

import (
	"github.com/yangfanlalala/blx/common"
	"github.com/yangfanlalala/blx/crypto"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	smsFormat           = "json"
	smsRegionID         = "cn-beijing"
	smsSignatureMethod  = "HMAC-SHA1"
	smsSignatureVersion = "1.0"
	smsVersion          = "2017-05-25"
	smsSendApiURL       = "dysmsapi.aliyuncs.com"
)

var (
	smsLock = &sync.Once{}
	smsCli  *smsClient
)

type smsSendResponse struct {
	RequestID string `json:"RequestId"`
	BizID     string `json:"BizId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
}

type smsClient struct {
	AccessKeyID  string
	AccessSecret string
	httpClient   *http.Client
}

type smsBaseParams struct {
	Signature        string `name:"Signature"`
	AccessKeyID      string `name:"AccessKeyId"`
	Action           string `name:"Action"`
	Format           string `name:"Format"`
	RegionId         string `name:"RegionId"`
	SignatureMethod  string `name:"SignatureMethod"`
	SignatureNonce   string `name:"SignatureNonce"`
	SignatureVersion string `name:"SignatureVersion"`
	Timestamp        string `name:"Timestamp"`
	Version          string `name:"Version"`
}

type smsSendParams struct {
	smsBaseParams
	PhoneNumbers  string `name:"PhoneNumbers"`
	SignName      string `name:"SignName"`
	TemplateCode  string `name:"TemplateCode"`
	OutID         string `name:"OutId"`
	TemplateParam string `name:"TemplateParam"`
}

func NewSmsClient(ak, as, endpoint string, param map[string]string) *smsClient {
	smsLock.Do(func() {
		httpClient := &http.Client{
			Timeout: 30,
		}
		smsCli = &smsClient{
			httpClient: httpClient,
		}
	})
	return smsCli
}

func (cli *smsClient) Send(tplCode, signName, requestID string, phoneNumbers []string, param map[string]interface{}) (string, string, error) {
	tplParam, _ := json.Marshal(param)
	params := map[string]string{
		"AccessKeyId":      cli.AccessKeyID,
		"Action":           "SendSms",
		"Format":           smsFormat,
		"RegionId":         smsRegionID,
		"SignatureMethod":  smsSignatureMethod,
		"SignatureVersion": smsSignatureVersion,
		"SignatureNonce":   time.Now().String(),
		"Timestamp":        time.Now().UTC().Format(common.ISO8601Format),
		"Version":          smsVersion,
		"PhoneNumbers":     strings.Join(phoneNumbers, ","),
		"SignName":         signName,
		"TemplateCode":     tplCode,
		"OutId":            requestID,
		"TemplateParam":    string(tplParam),
	}
	queryString, stringToSign := cli.buildStringToSign(params)
	fmt.Println(stringToSign)
	sign := common.AliSpecialEncode(crypto.HmacSha1(stringToSign, cli.AccessSecret+"&"))
	fmt.Println("Debug Sign", sign)
	request, err := http.NewRequest("GET", "https://"+smsSendApiURL+"?"+queryString+"&Signature="+sign, bytes.NewReader([]byte("")))
	if err != nil {
		return "", "", err
	}
	rsp, err := cli.httpClient.Do(request)
	if err != nil {
		return "", "", err
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(body))
	if err != nil {
		return "", "", err
	}
	if rsp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("requests aliyun sms service error, status[%d] msg[%s]", rsp.StatusCode, body)
	}
	result := &smsSendResponse{}
	json.Unmarshal(body, result)
	if result.Code != "OK" {
		return "", "", fmt.Errorf("requests aliyun sms service error, code[%s] msg[%s]", result.Code, result.Message)
	}
	return result.RequestID, result.BizID, nil
}

func (cli *smsClient) BatchSend() {

}

func (cli *smsClient) Fetch() {

}

func (cli *smsClient) buildStringToSign(params map[string]string) (queryString, stringToSign string) {
	val := &url.Values{}
	for k, v := range params {
		val.Add(k, v)
	}
	queryString = val.Encode()
	stringToSign = strings.Replace(queryString, "+", "%20", -1)
	stringToSign = strings.Replace(stringToSign, "~", "%7E", -1)
	stringToSign = strings.Replace(stringToSign, "*", "%2A", -1)
	stringToSign = "GET&%2F&" + common.AliSpecialEncode(stringToSign)
	return
}
