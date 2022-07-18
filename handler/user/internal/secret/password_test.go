package secret

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Password(t *testing.T) {
	password := "zhu88jie"
	encodePassword, err := EncodePassword(password)
	require.Nil(t, err, "%+v", err)
	require.NotEmpty(t, encodePassword)
	check := CheckPassword(password, encodePassword)
	require.True(t, check)
}

func TestSHA256(t *testing.T) {
	hash := sha256.New()
	_, err := hash.Write([]byte("admin"))
	if err != nil {
		t.Error(err)
	}
	password := hex.EncodeToString(hash.Sum(nil))
	fmt.Println(password)
	//check := CheckPassword(password, "$2a$10$40Y8k4BHvOVyHx6ulV1Fru01cMXfYsEzIpG79hKbxR869oovUmRre")
	//fmt.Println(check)
	//require.True(t, check)
}
