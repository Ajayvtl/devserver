package environments

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrEncryptionFailed = errors.New("failed to encrypt secret")
	ErrDecryptionFailed = errors.New("failed to decrypt secret")
)

// CryptoService handles encryption/decryption at rest using a master key.
type CryptoService interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// AESCryptoService implements AES-256-GCM encryption.
type AESCryptoService struct {
	key []byte
}

// NewAESCryptoService creates a new crypto service. The key must be exactly 32 bytes for AES-256.
func NewAESCryptoService(key []byte) (*AESCryptoService, error) {
	if len(key) != 32 {
		return nil, errors.New("AES key must be exactly 32 bytes")
	}
	return &AESCryptoService{key: key}, nil
}

func (s *AESCryptoService) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", ErrEncryptionFailed
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrEncryptionFailed
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", ErrEncryptionFailed
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (s *AESCryptoService) Decrypt(ciphertext string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return "", ErrDecryptionFailed
	}

	nonce, actualCiphertext := decoded[:nonceSize], decoded[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(plaintext), nil
}
