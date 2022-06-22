package user

import (
	"context"
	"encoding/json"

	"github.com/DataWorkbench/account/handler/user/internal/secret"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/common/rediswrap"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"gorm.io/gorm"
)

type sessionCache struct {
	UserSet *pbmodel.User
	KeySet  *pbmodel.AccessKey
}

func CreateSession(ctx context.Context, tx *gorm.DB, rdb rediswrap.Client, userName, password string, ignorePassWord bool) (
	userSet *pbmodel.User, sessionId string, err error) {

	var (
		keySet *pbmodel.AccessKey
	)

	userSet, err = DescribeUserByName(tx, userName)
	if err != nil {
		return
	}
	if !ignorePassWord && !secret.CheckPassword(password, userSet.Password) {
		err = qerror.UserNameOrPasswordError
		return
	}
	keySet, err = DescribePitrixAccessKeyByOwner(tx, userSet.UserId)
	if err != nil {
		return
	}

	sessionId = secret.GenerateSessionId()
	sessionValue := &sessionCache{
		UserSet: userSet,
		KeySet:  keySet,
	}

	jsonString, err := json.Marshal(sessionValue)
	if err != nil {
		return
	}
	_, err = rdb.SetNX(ctx, sessionId, jsonString, 60*60*24).Result()
	if err != nil {
		return
	}
	return
}

func CheckSession(ctx context.Context, rdb rediswrap.Client, sessionId string) (
	userSet *pbmodel.User, keySet *pbmodel.AccessKey, err error) {
	var jsonValue []byte

	jsonValue, err = rdb.Get(ctx, sessionId).Bytes()
	if err != nil {
		return
	}
	if len(jsonValue) == 0 {
		// FIXME: Define a new error.
		err = qerror.ResourceNotExists
		return
	}

	var sessionValue sessionCache

	if err = json.Unmarshal(jsonValue, &sessionValue); err != nil {
		return
	}
	userSet = sessionValue.UserSet
	keySet = sessionValue.KeySet
	return
}
