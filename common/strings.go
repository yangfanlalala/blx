package common

import (
	"net/url"
	"strings"
)

func AliSpecialEncode(v string) string {
	return strings.Replace(strings.Replace(url.QueryEscape(v), "+", "%20", -1), "~", "%7E", -1)
}
