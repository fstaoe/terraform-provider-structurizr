package auth

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Given pretzels", args{content: "I am getting thirsty."}, "1b0f985d0302a0b674c99bc56a0992f5"},
		{"Given no pretzels", args{content: ""}, "d41d8cd98f00b204e9800998ecf8427e"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Checksum(tt.args.content); got != tt.want {
				t.Errorf("Checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSign(t *testing.T) {
	type args struct {
		secret  string
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Given a lazy dog",
			args{secret: "api-secret", content: "The brown fox becomes lazy"},
			"2bf0c8ea72d1ebc6e4d38b8724604caf779def6b78c4e7e8f1cbbce75e8968c6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sign(tt.args.secret, tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}
