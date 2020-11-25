package bce

import (
	"blxee.com/utils/common"
	"fmt"
	"testing"
	"time"
)

func Test_bosClient_GetObject(t *testing.T) {
	type fields struct {
		Credentials *Credentials
		Bucket      string
		Endpoint    string
	}
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{Credentials: &Credentials{AccessKeyId: "65477071095849b89a1a7e02fa09fef4", SecretAccessKey: "3d8631dd64fc482e95f7e3dbd5b64b32"}, Bucket: "blxee", Endpoint: "bj.bcebos.com"},
			args:   args{source: "test/abc.jpg"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bos := &bosClient{
				Credentials: tt.fields.Credentials,
				Bucket:      tt.fields.Bucket,
				Endpoint:    tt.fields.Endpoint,
			}
			err := bos.GetObject(tt.args.source)
			fmt.Print(err)
		})
	}
}

func Test_bosClient_PutObject(t *testing.T) {
	type fields struct {
		Credentials *Credentials
		Bucket      string
		Endpoint    string
	}
	type args struct {
		source string
		path   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Testing",
			fields: fields{Credentials: &Credentials{AccessKeyId: "65477071095849b89a1a7e02fa09fef4", SecretAccessKey: "3d8631dd64fc482e95f7e3dbd5b64b32"}, Bucket: "blxee", Endpoint: "bj.bcebos.com"},
			args:   args{source: "test/abcdefj.jpg", path: "./timg.jpg"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bos := &bosClient{
				Credentials: tt.fields.Credentials,
				Bucket:      tt.fields.Bucket,
				Endpoint:    tt.fields.Endpoint,
			}
			if err := bos.PutObject(tt.args.source, tt.args.path); err != nil {
				fmt.Println(err)
			}
		})
	}
}

func Test_bosClient_PostObject(t *testing.T) {
	type fields struct {
		Credentials *Credentials
		Bucket      string
		Endpoint    string
	}
	type args struct {
		policy *BosPolicy
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "TESTING",
			fields: fields{Credentials: &Credentials{AccessKeyId: "01e92fed26c042aabb26137166f28971", SecretAccessKey: "2e264ea8a9a340deacb4d01233f85e08"}, Bucket: "blxee", Endpoint: "bj.bcebos.com"},
			args:   args{policy: &BosPolicy{Expiration: common.ISO8601Datetime(time.Now().UTC().Add(3600 * time.Second)), Conditions: []interface{}{map[string]string{"key": "test/user/*"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bos := &bosClient{
				Credentials: tt.fields.Credentials,
				Bucket:      tt.fields.Bucket,
				Endpoint:    tt.fields.Endpoint,
			}
			gotPolicy := bos.PostObject(tt.args.policy)
			fmt.Println(gotPolicy)
		})
	}
}
