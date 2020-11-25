package aliyun

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_smsClient_Send(t *testing.T) {
	type fields struct {
		AccessKeyID  string
		AccessSecret string
		httpClient   *http.Client
	}
	type args struct {
		tplCode      string
		signName     string
		requestID    string
		phoneNumbers []string
		param        map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "testing 01",
			fields: fields{
				AccessKeyID:  "LTAI4FkoLTZuYSNDXedprX8k",
				AccessSecret: "ecw0zKEKkkK438f2v4G1xxopLemDOe",
				httpClient:   &http.Client{},
			},
			args: args{
				tplCode:      "SMS_182684074",
				signName:     "勃利网",
				requestID:    "test01",
				phoneNumbers: []string{"18500609063"},
				param:        map[string]interface{}{"mtname": "test01", "submittime": "2020-01-28 01:01:01"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &smsClient{
				AccessKeyID:  tt.fields.AccessKeyID,
				AccessSecret: tt.fields.AccessSecret,
				httpClient:   tt.fields.httpClient,
			}
			got, got1, err := cli.Send(tt.args.tplCode, tt.args.signName, tt.args.requestID, tt.args.phoneNumbers, tt.args.param)
			fmt.Printf("%+v, %+v, %+v", got, got1, err)
		})
	}
}
