package bce

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewStsClient(t *testing.T) {
	type args struct {
		cred     *Credentials
		endpoint string
	}
	tests := []struct {
		name string
		args args
		want *stsClient
	}{
		{
			name: "Testing for STS",
			args: args{cred: &Credentials{AccessKeyId: "65477071095849b89a1a7e02fa09fef4", SecretAccessKey: "3d8631dd64fc482e95f7e3dbd5b64b32"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStsClient(tt.args.cred, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStsClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stsClient_GetSessionToken(t *testing.T) {
	type fields struct {
		Credentials *Credentials
		Endpoint    string
	}
	type args struct {
		access *accessControlList
		expire int64
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantSessToken *SessionToke
		wantErr       bool
	}{
		{
			name:   "Testing Sts Token",
			fields: fields{Credentials: &Credentials{AccessKeyId: "65477071095849b89a1a7e02fa09fef4", SecretAccessKey: "3d8631dd64fc482e95f7e3dbd5b64b32"}},
			args:   args{&accessControlList{Permission: []string{BosPermissionPutObject, BosPermissionGetObject}, Region: "bj", Resource: []string{"blxee/api"}, Service: ServiceOss, Effect: EffectAllow}, 36000},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sts := &stsClient{
				Credentials: tt.fields.Credentials,
				Endpoint:    tt.fields.Endpoint,
			}
			gotSessToken, err := sts.GetSessionToken(tt.args.access, tt.args.expire)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			fmt.Println("Session: ", gotSessToken)
		})
	}
}
