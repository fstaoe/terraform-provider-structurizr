package util

import (
	"strings"
)

// ConfigCompose concatenates multiple test configurations
func ConfigCompose(config ...string) string {
	var str strings.Builder
	for _, conf := range config {
		str.WriteString(conf)
	}
	return str.String()
}

// StringPtr returns a string pointer
func StringPtr(str string) *string { return &str }
