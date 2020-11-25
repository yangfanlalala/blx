package aliyun

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	ossLock = &sync.Once{}
	ossCli  *ossClient
)

const (
	ActionPutObject             = "oss:PutObject"
	ActionPostObject            = "oss:PostObject"
	ActionDeleteObject          = "oss:DeleteObject"
	ActionDeleteMultipleObjects = "oss:DeleteMultipleObjects"
	ActionGetObject             = "oss:GetObject"
)

type ossClient struct {
	AccessID     string
	AccessSecret string
	Bucket       string
	Endpoint     string
	httpClient   *http.Client
}

func NewOssClient(ak, as, bucket, endpoint string) *ossClient {
	ossLock.Do(func() {
		httpCli := &http.Client{
			Timeout: 30 * time.Second,
		}
		ossCli = &ossClient{
			AccessID:     ak,
			AccessSecret: as,
			Bucket:       bucket,
			Endpoint:     endpoint,
			httpClient:   httpCli,
		}
	})
	return ossCli
}

func (cli *ossClient) signedString(req *http.Request, canoncalizedResource string) string {
	t := make(map[string]string)
	for k, v := range req.Header {
		if strings.HasPrefix(strings.ToLower(k), "x-oss-") {
			t[strings.ToLower(k)] = v[0]
		}
	}
	//字典排序
	hs := newHeaderSorter(t)
	hs.Sort()

	canonicalizedHeader := ""
	for i := range hs.Keys {
		canonicalizedHeader += hs.Keys[i] + "\n"
	}
	date := req.Header.Get("Date")
	contentType := req.Header.Get("Content-Type")
	contentMd5 := req.Header.Get("Content-Md5")
	signedString := req.Method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedHeader + canoncalizedResource

	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(cli.AccessSecret))
	io.WriteString(h, signedString)
	return "OSS " + cli.AccessID + ":" + base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (cli *ossClient) getSignedString(stringToSign string) string {
	hash := hmac.New(func() hash.Hash {
		return sha1.New()
	}, []byte(cli.AccessSecret))
	io.WriteString(hash, stringToSign)
	return "OSS " + cli.AccessID + ":" + base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (cli *ossClient) PutObject(content []byte, object string) error {
	url := "https://" + cli.Bucket + "." + cli.Endpoint + "/" + strings.TrimLeft(object, "/")
	req, err := http.NewRequest("PUT", url, bytes.NewReader(content))
	if err != nil {
		return err
	}
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	req.Header.Set("Authorization", cli.signedString(req, "/"+cli.Bucket+"/"+strings.TrimLeft(object, "/")))
	rsp, err := cli.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("put oss object not ok, status code[%d], body[%s]", rsp.StatusCode, body)
	}
	return nil
}

func (cli *ossClient) GetObject() ([]byte, error) {
	return []byte{}, nil
}

func (cli *ossClient) DeleteObject() error {
	return nil
}

type PostForm struct {
	Endpoint    string
	AccessKeyID string
	Policy      string
	Signature   string
	ExpiredAt   string
}

func (cli *ossClient) WebPostObject(path string, sizeBytes int, expired time.Time) *PostForm {
	expires := expired.Format("2006-01-02T15:04:05.000Z")
	tpl := `{"expiration": "%s", "conditions":[{"bucket": "%s"}, ["starts-with", "$key", "%s"], ["content-length-range", 1, %d]]}`
	policy := fmt.Sprintf(tpl, expires, cli.Bucket, path, sizeBytes)
	policy = base64.StdEncoding.EncodeToString([]byte(policy))
	hmc := hmac.New(func() hash.Hash {
		return sha1.New()
	}, []byte(cli.AccessSecret))
	io.WriteString(hmc, policy)
	signature := base64.StdEncoding.EncodeToString(hmc.Sum([]byte("")))
	form := &PostForm{
		Policy:      policy,
		Signature:   signature,
		ExpiredAt:   expires,
		Endpoint:    "https://" + cli.Bucket + "." + cli.Endpoint,
		AccessKeyID: cli.AccessID,
	}
	return form
}
