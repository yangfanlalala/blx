package bce

import (
	"fmt"
	"testing"
)

func Test_smsClient_SendMessage(t *testing.T) {
	type fields struct {
		Credentials *Credentials
	}
	type args struct {
		payload *SmsPayload
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "Testing For Sms",
			fields: fields{Credentials: &Credentials{AccessKeyId: "01e92fed26c042aabb26137166f28971", SecretAccessKey: "2e264ea8a9a340deacb4d01233f85e08"}},
			args:   args{payload: &SmsPayload{InvokeId: "d2iL60UT-7N8-2236", PhoneNumber: "18500609063", TemplateCode: "smsTpl:e7476122a1c24e37b3b0de19d04ae900", ContentVar: SmsContentVar{"code": "6767"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &smsClient{
				Credentials: tt.fields.Credentials,
			}
			err := cli.SendMessage(tt.args.payload)
			fmt.Println(err)
		})
	}
}
