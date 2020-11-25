package bce

import (
	"blxee.com/utils/common"
	"blxee.com/utils/crypto"
	blxHttp "blxee.com/utils/network"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	AuthVersionV2    = "bce-auth-v2"
	SDKVersion       = "blx-go-sdk/2.0"
	SignSeparator    = "\n"
	RequestSeparator = "\n"
	FieldSeparator   = "/"
	HeaderSeparator  = ";"
)

type Signer interface {
	Sign(*http.Request, *Credentials, *signOptions)
}

type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

type signOptions struct {
	Date    time.Time
	Region  string
	Service string
}

func NewSignOptions(date time.Time, region, service string) *signOptions {
	return &signOptions{date, region, service}
}

type AuthSigner struct {
}

func (*AuthSigner) Sign(req *http.Request, cred *Credentials, opt *signOptions) string {
	req.Header.Set("x-bce-date", time.Now().UTC().Format(common.ISO8601Format))
	req.Header.Set("Host", req.Host)
	date := opt.Date.Format("20060102")
	prefix := AuthVersionV2 + FieldSeparator + cred.AccessKeyId + FieldSeparator + date + FieldSeparator + opt.Region + FieldSeparator + opt.Service
	signKey := crypto.HmacSha256Hex(cred.SecretAccessKey, prefix)
	canonicalRequest, headers := getCanonicalRequest(req)
	signature := crypto.HmacSha256Hex(signKey, canonicalRequest)
	return prefix + FieldSeparator + headers + FieldSeparator + signature
}

//构建CanonicalRequest
func getCanonicalRequest(req *http.Request) (string, string) {
	canonicalURI := getCanonicalURI(req.URL.Path)
	canonicalQueryString := getCanonicalQueryString(req.URL.Query())
	canonicalHeaders, headers := getCanonicalHeaders(req.Header)
	return req.Method + RequestSeparator + canonicalURI + RequestSeparator + canonicalQueryString + RequestSeparator + canonicalHeaders, headers
}

//构建CanonicalURI
func getCanonicalURI(uri string) string {
	if len(uri) == 0 {
		return "/"
	}
	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}
	return blxHttp.UriEncodeExceptSlash(uri)
}

//构建CanonicalQueryString
func getCanonicalQueryString(pairs url.Values) string {
	length := len(pairs)
	if length == 0 {
		return ""
	}
	result := make([]string, 0, length)
	for key, value := range pairs {
		if key == "" {
			continue
		}
		item := ""
		if len(value) == 0 {
			item += blxHttp.UriEncode(key) + "="
		} else {
			item += blxHttp.UriEncode(key) + "=" + blxHttp.UriEncode(value[0])
		}
		result = append(result, item)
	}
	sort.Strings(result)
	return strings.Join(result, "&")
}

//构建CanonicalHeaders
func getCanonicalHeaders(headers http.Header) (string, string) {
	container := make([]string, 0, len(headers))
	fields := make([]string, 0, len(headers))
	for key, value := range headers {
		if key == "" {
			continue
		}
		lowerKey := strings.ToLower(key)
		headValue := strings.TrimSpace(value[0])
		container = append(container, blxHttp.UriEncode(lowerKey)+":"+blxHttp.UriEncode(headValue))
		fields = append(fields, lowerKey)
	}
	sort.Strings(container)
	sort.Strings(fields)
	return strings.Join(container, SignSeparator), strings.Join(fields, HeaderSeparator)
}
