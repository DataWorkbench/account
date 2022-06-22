package secret

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func EncodePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password string, encodePassword string) bool {
	ret := bcrypt.CompareHashAndPassword([]byte(encodePassword), []byte(password))
	return ret == nil
}

// GenerateAccessKey returns a pair of access key id and secret access key.
func GenerateAccessKey() (string, string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	rand.Seed(time.Now().UnixNano())
	accessKeyId := make([]byte, accessKeyIdLength)
	secretKey := make([]byte, secretKeyLength)
	for i := 0; i < accessKeyIdLength; i++ {
		accessKeyId[i] = accessKeyIdLetters[r.Intn(len(accessKeyIdLetters))]
	}
	for i := 0; i < secretKeyLength; i++ {
		secretKey[i] = secretKeyLetters[r.Intn(len(secretKeyLetters))]
	}

	return string(accessKeyId), string(secretKey)
}

func GenerateSessionId() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	session := make([]byte, sessionLength)
	for i := 0; i < sessionLength; i++ {
		session[i] = sessionLetters[r.Intn(len(sessionLetters))]
	}
	return string(session)
}
