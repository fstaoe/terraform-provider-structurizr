package util

import (
	"reflect"
	"testing"
)

func TestConfigCompose(t *testing.T) {
	type args struct {
		config []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"No config strings", args{[]string{}}, ""},
		{"Single config string", args{[]string{"single"}}, "single"},
		{"Multiple config strings", args{[]string{"config1", "config2", "config3"}}, "config1config2config3"},
		{"Config strings with spaces", args{[]string{"config1 ", " config2", " config3 "}}, "config1  config2 config3 "},
		{"Empty strings in config", args{[]string{"config1", "", "config3"}}, "config1config3"},
		{"Only empty strings", args{[]string{"", "", ""}}, ""},
		{"Unicode strings", args{[]string{"你好", "世界", "🌍"}}, "你好世界🌍"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConfigCompose(tt.args.config...); got != tt.want {
				t.Errorf("ConfigCompose() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringPtr(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Empty string", args{""}, ""},
		{"Non-empty string", args{"hello"}, "hello"},
		{"String with spaces", args{"hello world"}, "hello world"},
		{"String with special characters", args{"!@#$%^&*()_+"}, "!@#$%^&*()_+"},
		{"Unicode string", args{"你好，世界"}, "你好，世界"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringPtr(tt.args.str); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("StringPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}
