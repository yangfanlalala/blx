package network

import (
	"bytes"
	"net/http"
	"time"
)

const (
	HttpMethodPost   = "POST"
	HttpMethodGet    = "GET"
	HttpMethodHead   = "HEAD"
	HttpMethodPut    = "PUT"
	HttpMethodPatch  = "PATCH"
	HttpMethodDelete = "DELETE"
	HttpMethodOption = "OPTIONS"
)

const (
	ContentTypeJSON       = "application/json"
	ContentTypeXml        = "application/xml"
	ContentTypeForm       = "application/x-www-form-urlencoded"
	ContentTypeFormData   = "application/x-www-form-urlencoded"
	ContentTypeURLEncoded = "application/x-www-form-urlencoded"
	ContentTypeHTML       = "text/html"
	ContentTypeText       = "text/plain"
	ContentTypeMultipart  = "multipart/form-data"
	ContentTypeJavascript = "application/javascript"
	ContentTypePdf        = "application/pdf"
	ContentTypeImageJpeg  = "image/jpeg"
	ContentTypeImageGif   = "image/gif"
	ContentTypeImagePng   = "image/png"
	ContentTypeStream     = "application/octet-stream"
)

const (
	MaxIdelConnsPerHost = 50
	MaxIdleConns        = 100
	ResponseHeadTimeout = 60 * time.Second
	DialTimeout         = 30 * time.Second
	SmallInterval       = 30 * time.Second
	LargeInterval       = 1200 * time.Second
)

const (
	UserAgent = "blx-go-sdk"
	GMTFormat = ""
)

func NewHttpRequest(method, url string, body []byte) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", UserAgent)
	return
}

type httpClient struct {
	http.Client
}

func NewHttpClient() *httpClient {
	transport := &http.Transport{
		//MaxIdleConnsPerHost:   MaxIdelConnsPerHost,
		//MaxIdleConns:          MaxIdleConns,
		//ResponseHeaderTimeout: ResponseHeadTimeout,
		//IdleConnTimeout:       DialTimeout,
	}
	client := http.Client{Transport: transport}
	return &httpClient{
		client,
	}
}

func (cli *httpClient) SetTimeout(dur time.Duration) {
	cli.Timeout = dur
}

func (cli *httpClient) SetTransport(tran http.RoundTripper) {
	cli.Transport = tran
}
