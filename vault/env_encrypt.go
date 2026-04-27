package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

// EncryptResult holds the encrypted value and metadata.
type EncryptResult struct {
	Key       string
	Encrypted string
	WasPlain  bool
}

// EncryptSecrets encrypts secret values matching the given key patterns using AES-GCM.
// passphrase is used to derive a 32-byte key via SHA-256.
func EncryptSecrets(secrets map[string]string, patterns []string, passphrase string) (map[string]string, []EncryptResult, error) {
	if passphrase == "" {
		return nil, nil, errors.New("passphrase must not be empty")
	}
	key := deriveKey(passphrase)
	out := make(map[string]string, len(secrets))
	var results []EncryptResult
	for k, v := range secrets {
		if isEncrypted(v) || !matchesAnyPattern(k, patterns) {
			out[k] = v
			continue
		}
		enc, err := aesGCMEncrypt(key, v)
		if err != nil {
			return nil, nil, fmt.Errorf("encrypt %q: %w", k, err)
		}
		out[k] = enc
		results = append(results, EncryptResult{Key: k, Encrypted: enc, WasPlain: true})
	}
	return out, results, nil
}

// DecryptSecrets decrypts all enc: prefixed values using the given passphrase.
func DecryptSecrets(secrets map[string]string, passphrase string) (map[string]string, error) {
	if passphrase == "" {
		return nil, errors.New("passphrase must not be empty")
	}
	key := deriveKey(passphrase)
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if !isEncrypted(v) {
			out[k] = v
			continue
		}
		plain, err := aesGCMDecrypt(key, v)
		if err != nil {
			return nil, fmt.Errorf("decrypt %q: %w", k, err)
		}
		out[k] = plain
	}
	return out, nil
}

func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

func isEncrypted(v string) bool {
	return strings.HasPrefix(v, "enc:")
}

func aesGCMEncrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return "enc:" + base64.StdEncoding.EncodeToString(ciphertext), nil
}

func aesGCMDecrypt(key []byte, value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, "enc:"))
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("gcm open: %w", err)
	}
	return string(plain), nil
}

func matchesAnyPattern(key string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}
	for _, p := range patterns {
		if ok, _ := matchesPattern(key, p); ok {
			return true
		}
	}
	return false
}
