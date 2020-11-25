package bce

import (
	blxHttp "github.com/yangfanlalala/blx/network"
	"net/http"
)

type Client struct {
	credentials Credentials
	request *http.Request
}

func NewClient(ak, sk string) *Client {
	cli := new(Client)
	cli.credentials = Credentials{
		AccessKeyId:     ak,
		SecretAccessKey: sk,
	}
	cli.request.Header.Add("Content-Type", blxHttp.CONTENT_TYPE_JSON)
	cli.request.Header.Add("", "")

	return cli
}
