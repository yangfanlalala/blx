package bce

import (
	"blxee.com/utils/common"
	"blxee.com/utils/crypto"
	blxHttp "blxee.com/utils/network"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

var (
	BCE_AUTH_VERSION   = "bce-auth-v1"
	SIGN_JOINER        = "\n"
	SIGN_HEADER_JOINER = ";"
	HEADERS_TO_SIGN    = map[string]struct{}{
		"host":           {},
		"content-type":   {},
		"content-length": {},
		"content-md5":    {},
	}
)

const (
	SDK_VERSION = "blx-go-sdk/1.0"
)

type Signer interface {
	Sign(*http.Request, *Credentials, *SignOptions)
}

type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

func (b *Credentials) String() string {
	str := "ak: " + b.AccessKeyId + ", sk: " + b.SecretAccessKey
	if len(b.SessionToken) > 0 {
		return str + ", sessionToken: " + b.SessionToken
	}
	return str
}

type SignOptions struct {
	Timestamp     int64
	ExpireSeconds int
}

type AuthSigner struct{}

func (b *AuthSigner) Sign(req *http.Request, cred *Credentials, opt *SignOptions) string {
	timestamp := common.FormatISO8601Date(opt.Timestamp)
	authStringPrefix := fmt.Sprintf("%s/%s/%s/%d", BCE_AUTH_VERSION, cred.AccessKeyId, timestamp, opt.ExpireSeconds)
	signingKey := crypto.HmacSha256Hex(cred.SecretAccessKey, authStringPrefix)
	canonicalHeader, signHeader := getCanonicalHeaders(req.Header)
	fmt.Println(canonicalHeader)
	canonicalRequestSet := []string{req.Method, getCanonicalURIPath(req.URL.RequestURI()), getCanonicalQueryString(req.URL.Query()), canonicalHeader}
	canonicalRequest := strings.Join(canonicalRequestSet, SIGN_JOINER)
	signature := crypto.HmacSha256Hex(signingKey, canonicalRequest)
	return authStringPrefix + "/" + signHeader + "/" + signature
}

func (b *AuthSigner) getAuthStringPrefix() string {
	return ""
}

func getCanonicalURIPath(path string) string {
	if len(path) == 0 {
		return "/"
	}
	canonicalPath := path
	if strings.HasPrefix(path, "/") {
		canonicalPath = path[1:]
	}
	return "/" + blxHttp.UriEncodeExceptSlash(canonicalPath)
}

func getCanonicalQueryString(params url.Values) string {
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

func getCanonicalHeaders(headers http.Header) (string, string) {
	var canonicalHeaders []string
	var signHeaders []string
	for key, val := range headers {
		headKey := strings.ToLower(key)
		if headKey == strings.ToLower("") {
			continue
		}
		_, headExists := HEADERS_TO_SIGN[headKey]
		if headExists || (strings.HasPrefix(headKey, "x-bce-")) {
			headValue := strings.TrimSpace(val[0])
			encodedValue := blxHttp.UriEncode(headKey) + ":" + blxHttp.UriEncode(headValue)
			canonicalHeaders = append(canonicalHeaders, encodedValue)
			signHeaders = append(signHeaders, headKey)
		}
	}
	sort.Strings(canonicalHeaders)
	sort.Strings(signHeaders)
	return strings.Join(canonicalHeaders, SIGN_JOINER), strings.Join(signHeaders, SIGN_HEADER_JOINER)
}
