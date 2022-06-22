package secret

import (
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
