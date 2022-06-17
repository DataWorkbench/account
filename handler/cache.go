package handler

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/common/rediswrap"
	"github.com/DataWorkbench/glog"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	rdb rediswrap.Client
	//ctx context.Context
}

func (cache *Cache) set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonString, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = cache.rdb.SetNX(ctx, key, jsonString, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) delete(ctx context.Context, key string) error {
	_, err := cache.rdb.Del(ctx, key).Result()
	return err
}

func (cache *Cache) get(ctx context.Context, key string, value interface{}) error {
	jsonValue, err := cache.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	if len(jsonValue) == 0 {
		return qerror.ResourceNotExists
	}

	if err = json.Unmarshal(jsonValue, value); err != nil {
		return err
	}
	return nil
}

func (cache *Cache) GetPrefixKey(source string, resource string) string {
	return constants.Account + constants.RedisSeparator + source + constants.RedisSeparator + resource + constants.RedisSeparator
}

func (cache *Cache) CacheUsers(ctx context.Context, users []*pbmodel.User, source string) error {
	for _, u := range users {
		if err := cache.CacheUser(ctx, u, u.UserId, source); err != nil {
			return err
		}
	}
	return nil
}

func (cache *Cache) CacheUser(ctx context.Context, u *pbmodel.User, userID string, source string) error {
	if u == nil {
		return cache.CacheNotExistUser(ctx, userID, source)
	}
	prefixKey := cache.GetPrefixKey(source, constants.UserTableName)
	return cache.set(ctx, prefixKey+userID, u, time.Second*time.Duration(constants.UserCacheBaseSeconds+rand.Intn(constants.UserCacheRandomSeconds)))
}

func (cache *Cache) CacheNotExistUser(ctx context.Context, userID string, source string) error {
	prefixKey := cache.GetPrefixKey(source, constants.UserTableName)
	_, err := cache.rdb.SetNX(ctx, prefixKey+userID, nil, time.Second*time.Duration(constants.NotExistResourceCacheSeconds)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) GetUser(ctx context.Context, userID string, source string) (*pbmodel.User, error) {
	prefixKey := cache.GetPrefixKey(source, constants.UserTableName)
	var user pbmodel.User
	err := cache.get(ctx, prefixKey+userID, &user)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (cache *Cache) DelUser(ctx context.Context, userID string, withSession bool) error {
	key := cache.GetPrefixKey(constants.LocalSource, constants.UserTableName) + userID
	if err := cache.delete(ctx, key); err != nil {
		return err
	}
	if withSession {
		userKey := cache.GetPrefixKey(constants.LocalSource, constants.SessionPrefix) + userID
		var expiredSession string
		if err := cache.get(ctx, userKey, &expiredSession); err != nil {
			return err
		}
		cache.delete(ctx, userKey)
		cache.DeleteSession(ctx, expiredSession)
	}
	return nil
}

func (cache *Cache) CacheAccessKey(ctx context.Context, k *executor.AccessKey, accessKeyID string, source string) error {
	if k == nil {
		return cache.CacheNotExistAccessKey(ctx, accessKeyID, source)
	}
	prefixKey := cache.GetPrefixKey(source, constants.AccessKeyTableName)
	return cache.set(ctx, prefixKey+accessKeyID, k, time.Second*time.Duration(constants.AccessKeyCacheBaseSeconds))

}

func (cache *Cache) CacheNotExistAccessKey(ctx context.Context, accessKeyID string, source string) error {
	prefixKey := cache.GetPrefixKey(source, constants.AccessKeyTableName)
	_, err := cache.rdb.SetNX(ctx, prefixKey+accessKeyID, nil, time.Second*time.Duration(constants.NotExistResourceCacheSeconds)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) GetAccessKey(ctx context.Context, accessKeyID string, source string) (*executor.AccessKey, error) {
	prefixKey := cache.GetPrefixKey(source, constants.AccessKeyTableName)
	var accessKey executor.AccessKey
	err := cache.get(ctx, prefixKey+accessKeyID, &accessKey)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return &accessKey, nil
}

func (cache *Cache) GetSession(ctx context.Context, sessionId string) (*executor.AccessKey, error) {
	key := cache.GetPrefixKey(constants.LocalSource, constants.SessionPrefix) + sessionId
	var accessKey executor.AccessKey
	err := cache.get(ctx, key, &accessKey)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return &accessKey, nil
}

func (cache *Cache) DeleteSession(ctx context.Context, sessionId string) {
	key := cache.GetPrefixKey(constants.LocalSource, constants.SessionPrefix) + sessionId
	cache.delete(ctx, key)
}

func (cache *Cache) CacheSession(ctx context.Context, k *executor.AccessKey, sessionId string, userId string) error {
	lg := glog.FromContext(ctx)
	userKey := cache.GetPrefixKey(constants.LocalSource, constants.SessionPrefix) + userId
	sessionKey := cache.GetPrefixKey(constants.LocalSource, constants.SessionPrefix) + sessionId
	var expiredSession string
	if err := cache.get(ctx, userKey, &expiredSession); err != nil {
		if err == redis.Nil {
			lg.Debug().String("ignore key not exists error", err.Error()).Fire()
		} else {
			return err
		}
	}
	cache.DeleteSession(ctx, expiredSession)
	if err := cache.set(ctx, userKey, sessionId, time.Second*time.Duration(constants.SessionCacheSeconds)); err != nil {
		return err
	}
	return cache.set(ctx, sessionKey, k, time.Second*time.Duration(constants.SessionCacheSeconds))
}
