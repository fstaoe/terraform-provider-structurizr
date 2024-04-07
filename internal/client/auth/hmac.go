package auth

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
)

// Checksum returns a hexadecimal encoded MD5 checksum
func Checksum(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

// Enc returns base64 encoding of content
func Enc(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}

// Sign returns a HMAC Hash computed with SHA256
func Sign(secret string, content string) string {
	h := hmac.New(sha256.New, []byte(secret))
	_, _ = io.WriteString(h, content)
	return hex.EncodeToString(h.Sum(nil))
}

// Concat appends a \n to the given values and return a concatenated string
func Concat(values ...string) string {
	r := strings.Builder{}
	for _, v := range values {
		r.WriteString(v)
		r.WriteString("\n")
	}
	return r.String()
}
