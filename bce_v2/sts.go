package bce

import (
	"github.com/yangfanlalala/blx/network"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ServiceOss = "bce:oss"

	EffectAllow = "Allow"
	EffectDeny  = "Deny"
)

type stsRequest struct {
	//Id string `json:"id"`
	AccessControlList []*accessControlList `json:"accessControlList"`
}

type accessControlList struct {
	Eid        string   `json:"eid"`
	Service    string   `json:"logic"`
	Region     string   `json:"region"`
	Effect     string   `json:"effect"`
	Resource   []string `json:"resource"`
	Permission []string `json:"permission"`
}

func NewAccessControlList(region, effect, service string, resource, permission []string) *accessControlList {
	return &accessControlList{
		Service:    service,
		Region:     region,
		Effect:     effect,
		Resource:   resource,
		Permission: permission,
	}
}

type SessionToke struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	CreateTime      string `json:"createTime"`
	Expiration      string `json:"expiration"`
	UserId          string `json:"userId"`
}

type stsClient struct {
	Credentials *Credentials
	Endpoint    string
}

func NewStsClient(cred *Credentials, endpoint string) *stsClient {
	return &stsClient{Credentials: cred, Endpoint: endpoint}
}

func (sts *stsClient) GetSessionToken(access *accessControlList, expire int64) (sessToken *SessionToke, err error) {
	reqs, err := json.Marshal(stsRequest{AccessControlList: []*accessControlList{access}})
	if err != nil {
		return
	}
	fmt.Println(string(reqs))
	req, err := network.NewHttpRequest(network.HttpMethodPost, fmt.Sprintf("https://sts.bj.baidubce.com/v1/sessionToken?durationSeconds=%d", expire), reqs)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", network.ContentTypeJSON)
	//signer := AuthSigner{}
	//signature := signer.Sign(req, sts.Credentials, &signOptions{time.Now().UTC(), access.Region, access.Service})
	signer := AuthSignerV1{}
	signature := signer.Sign(req, sts.Credentials, &SignOptionsV1{time.Now().UTC(), 900})
	fmt.Println("Signature: ", signature)
	req.Header.Set("Authorization", signature)
	cli := network.NewHttpClient()
	res, err := cli.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("http requests is not ok, status code [%d] message [%s]", res.StatusCode, string(body))
		return
	}
	sessToken = new(SessionToke)
	err = json.Unmarshal(body, sessToken)
	return
}
