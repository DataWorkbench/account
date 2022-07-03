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
	hash.Write([]byte("zhu88jie"))
	password := hex.EncodeToString(hash.Sum(nil))
	check := CheckPassword(password, "$2a$10$40Y8k4BHvOVyHx6ulV1Fru01cMXfYsEzIpG79hKbxR869oovUmRre")
	fmt.Println(check)
	require.True(t, check)
}
