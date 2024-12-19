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

func Write(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > 4096 {
		return errors.New("cookie is too long")
	}

	http.SetCookie(w, &cookie)

	return nil
}

func Read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", errors.New("invalid cookie value")
	}

	return string(value), nil
}

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
		return "", errors.New("invalid cookie value")
	}

	nonce := encryptedValue[:nonceSize]
	ciphertext := encryptedValue[nonceSize:]

	plaintext, err := aesGCM.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", errors.New("invalid cookie value")
	}

	expectedName, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return "", errors.New("invalid cookie value")
	}

	if expectedName != name {
		return "", errors.New("invalid cookie value")
	}

	return value, nil
}
