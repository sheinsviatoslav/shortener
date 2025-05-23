// Package auth is responsible for writing and reading cookies
package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ErrInvalidCookieValue is the universal invalid cookie error
var ErrInvalidCookieValue = errors.New("invalid cookie value")

// Write is the function that encodes cookie value and sets it to response
func Write(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > 4096 {
		return errors.New("cookie is too long")
	}

	http.SetCookie(w, &cookie)

	return nil
}

// Read is the function that reads cookie value by name and decodes it
func Read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidCookieValue
	}

	return string(value), nil
}

// WriteEncryptedCookie is the function that encrypts cookie using secret key and writes it to response
func WriteEncryptedCookie(w http.ResponseWriter, cookie http.Cookie, secretKey []byte) error {
	aesBlock, err := aes.NewCipher(secretKey)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)
	encryptedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	cookie.Value = string(encryptedValue)

	return Write(w, cookie)
}

// ReadEncryptedCookie is the function that reads cookie using secret key and decrypts it
func ReadEncryptedCookie(r *http.Request, name string, secretKey []byte) (string, error) {
	encryptedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}

	aesBlock, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()

	if len(encryptedValue) < nonceSize {
		return "", ErrInvalidCookieValue
	}

	nonce := encryptedValue[:nonceSize]
	ciphertext := encryptedValue[nonceSize:]

	plaintext, err := aesGCM.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", ErrInvalidCookieValue
	}

	expectedName, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return "", ErrInvalidCookieValue
	}

	if expectedName != name {
		return "", ErrInvalidCookieValue
	}

	return value, nil
}
