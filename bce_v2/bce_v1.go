package bce

import (
	"github.com/yangfanlalala/blx/common"
	"github.com/yangfanlalala/blx/crypto"
	blxHttp "github.com/yangfanlalala/blx/network"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var (
	BceAuthVersionV1 = "bce-auth-v1"
	SignJoiner       = "\n"
	SignHeaderJoiner = ";"
	HeaderToSign     = map[string]struct{}{
		"host":           {},
		"content-type":   {},
		"content-length": {},
		"content-md5":    {},
	}
)

const (
	SDK_VERSION = "blx-go-sdk/1.0"
)

type CredentialsV1 struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

func (b *CredentialsV1) String() string {
	str := "ak: " + b.AccessKeyId + ", sk: " + b.SecretAccessKey
	if len(b.SessionToken) > 0 {
		return str + ", sessionToken: " + b.SessionToken
	}
	return str
}

type SignOptionsV1 struct {
	Timestamp     time.Time
	ExpireSeconds int
}

type AuthSignerV1 struct{}

func (b *AuthSignerV1) Sign(req *http.Request, cred *Credentials, opt *SignOptionsV1) string {
	timestamp := opt.Timestamp.Format(common.ISO8601Format)
	req.Header.Set("Host", req.URL.Host)
	req.Header.Set("x-bce-date", timestamp)
	if cred.SessionToken != "" {
		req.Header.Set("x-bce-security-token", cred.SessionToken)
	}
	authStringPrefix := fmt.Sprintf("%s/%s/%s/%d", BceAuthVersionV1, cred.AccessKeyId, timestamp, opt.ExpireSeconds)
	signingKey := crypto.HmacSha256Hex(cred.SecretAccessKey, authStringPrefix)
	canonicalHeader, signHeader := V1CanonicalHeader(req.Header)
	canonicalRequestSet := []string{req.Method, V1CanonicalUri(req.URL.Path), V1CanonicalQuery(req.URL.Query()), canonicalHeader}
	canonicalRequest := strings.Join(canonicalRequestSet, SignJoiner)
	fmt.Println(canonicalRequest)
	signature := crypto.HmacSha256Hex(signingKey, canonicalRequest)
	return authStringPrefix + "/" + signHeader + "/" + signature
}

func V1CanonicalUri(path string) string {
	if len(path) == 0 {
		return "/"
	}
	canonicalPath := path
	if strings.HasPrefix(path, "/") {
		canonicalPath = path[1:]
	}
	return "/" + blxHttp.UriEncodeExceptSlash(canonicalPath)
}

func V1CanonicalQuery(params url.Values) string {
	if len(params) == 0 {
		return ""
	}
	result := make([]string, 0, len(params))
	for key, value := range params {
		if strings.ToLower(key) == strings.ToLower("") {
			continue
		}
		var item string
		if len(value) == 0 {
			item = blxHttp.UriEncode(key) + "="
		} else {
			item = blxHttp.UriEncode(key) + "=" + blxHttp.UriEncode(value[0])
		}
		result = append(result, item)
	}
	sort.Strings(result)
	return strings.Join(result, "&")
}

func V1CanonicalHeader(headers http.Header) (string, string) {
	var canonicalHeaders []string
	var signHeaders []string
	for key, val := range headers {
		headKey := strings.ToLower(key)
		if headKey == strings.ToLower("") {
			continue
		}
		_, headExists := HeaderToSign[headKey]
		if headExists || (strings.HasPrefix(headKey, "x-bce-")) {
			headValue := strings.TrimSpace(val[0])
			encodedValue := blxHttp.UriEncode(headKey) + ":" + blxHttp.UriEncode(headValue)
			canonicalHeaders = append(canonicalHeaders, encodedValue)
			signHeaders = append(signHeaders, headKey)
		}
	}
	sort.Strings(canonicalHeaders)
	sort.Strings(signHeaders)
	return strings.Join(canonicalHeaders, SignJoiner), strings.Join(signHeaders, SignHeaderJoiner)
}
