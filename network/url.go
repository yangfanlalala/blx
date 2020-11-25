package network

import (
	"bytes"
	"fmt"
)

func UriEncode(str string) string {
	return uriEscape(str, false)
}

func UriEncodeExceptSlash(str string) string {
	return uriEscape(str, true)
}

func uriEscape(str string, escapeSlash bool) string {
	var byteBuffer bytes.Buffer
	for _, b := range []byte(str) {
		if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z' || (b >= '0' && b <= '9')) ||
			b == '-' || b == '_' || b == '.' || b == '~' || (b == '/' && escapeSlash) {
			byteBuffer.WriteByte(b)
		} else {
			byteBuffer.WriteString(fmt.Sprintf("%%%02X", b))
		}
	}
	return byteBuffer.String()
}
