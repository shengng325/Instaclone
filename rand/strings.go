package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	RememberTokenBytes = 32
)

//Bytes help generate n random bytes
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NBytes(base64string string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64string)
	if err != nil {
		return -1, err
	}
	return len(b), nil
}

//Strings return
func Strings(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	s := base64.URLEncoding.EncodeToString(b)
	return s, nil
}

func RememberToken() (string, error) {
	return Strings(RememberTokenBytes)
}
