package vault

import (
	"encoding/base64"
	"fmt"
)

// Base64Encode encodes a value to base64.
func Base64Encode(v string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(v)), nil
}

// Base64Decode decodes a base64-encoded value.
func Base64Decode(v string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	return string(b), nil
}
