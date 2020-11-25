package bce

import (
	"blxee.com/utils/common"
	"blxee.com/utils/crypto"
	"blxee.com/utils/network"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	BosPermissionGetObject     = "GetObject"
	BosPermissionPutObject     = "PutObject"
	BosPermissionDeleteObject  = "DeleteObject"
	BosPermissionListObjects   = "ListObjects"
	BosPermissionGetObjectMeta = "GetObjectMeta"
)

type BosPolicy struct {
	Expiration common.ISO8601Datetime `json:"expiration"`
	Conditions []interface{}          `json:"conditions"`
}

type BosPost struct {
	Signature string `json:"signature"`
	Policy    string `json:"policy"`
}

type bosClient struct {
	Credentials *Credentials
	Bucket      string
	Endpoint    string
}

func NewBosClient(cred *Credentials, bucket, endpoint string) *bosClient {
	return &bosClient{cred, bucket, endpoint}
}

func (bos *bosClient) GetObject(source string) (err error) {
	req, err := network.NewHttpRequest(network.HttpMethodGet, "https://"+bos.Bucket+"."+bos.Endpoint+"/"+strings.TrimLeft(source, "/"), []byte(""))
	if err != nil {
		return
	}
	//req.Header.Set("Date", time.Now().UTC().Format(time.RFC3339Nano))
	//signature := new(AuthSigner).Sign(req, bos.Credentials, &signOptions{time.Now().UTC(), "bj", "bos"})
	signature := new(AuthSignerV1).Sign(req, bos.Credentials, &SignOptionsV1{Timestamp: time.Now().UTC(), ExpireSeconds: 3600})
	fmt.Println(signature)
	req.Header.Set("Authorization", signature)
	cli := network.NewHttpClient()
	res, err := cli.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("http requests is not ok, status code[%d] message[%s]", res.StatusCode, body)
		return
	}
	fmt.Println(string(body))
	return nil
}

func (bos *bosClient) PutObject(source, path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	size := buf.Size()
	body := make([]byte, 0, size)
	for {
		byt, err := buf.ReadByte()
		if err != nil {
			break
		}
		body = append(body, byt)
	}
	req, err := network.NewHttpRequest(network.HttpMethodPut, "http://"+bos.Bucket+"."+bos.Endpoint+"/"+strings.TrimLeft(source, "/"), body)
	//req.Header.Set("Content-Length", strconv.Itoa(size))
	signature := new(AuthSignerV1).Sign(req, bos.Credentials, &SignOptionsV1{Timestamp: time.Now().UTC(), ExpireSeconds: 3600})
	req.Header.Set("Authorization", signature)
	cli := network.NewHttpClient()
	res, err := cli.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("http requests is not ok, status[%d] message[%s]", res.StatusCode, body)
		return
	}
	return nil
}

func (bos *bosClient) DeleteObject() {

}

func (bos *bosClient) PostObject(policy *BosPolicy) (policyRes *BosPost) {
	jason, _ := json.Marshal(policy)
	base64policy := base64.StdEncoding.EncodeToString(jason)
	signature := crypto.HmacSha256Hex(bos.Credentials.SecretAccessKey, base64policy)
	policyRes = &BosPost{Signature: signature, Policy: base64policy}
	return
}
