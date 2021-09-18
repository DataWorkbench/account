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
	"github.com/DataWorkbench/gproto/pkg/accountpb"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	rdb         rediswrap.Client
	cacheEnable map[string]bool
	ctx         context.Context
}

func (cache *Cache) set(key string, value interface{}, expiration time.Duration) error {
	jsonString, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = cache.rdb.SetNX(cache.ctx, key, jsonString, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) get(key string, value interface{}) error {
	jsonValue, err := cache.rdb.Get(cache.ctx, key).Bytes()
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

func (cache *Cache) IsEnable(source string) bool {
	enable, ok := cache.cacheEnable[source]
	if !ok {
		return false
	} else {
		return enable
	}
}

func (cache *Cache) GetPrefixKey(source string, resource string) string {
	return constants.Account + constants.RedisSeparator + source + constants.RedisSeparator + resource + constants.RedisSeparator
}

func (cache *Cache) CacheUsers(users []*accountpb.User, source string) error {
	for _, u := range users {
		if err := cache.CacheUser(u, u.UserId, source); err != nil {
			return err
		}
	}
	return nil
}

func (cache *Cache) CacheUser(u *accountpb.User, userID string, source string) error {
	if u == nil {
		return cache.CacheNotExistUser(userID, source)
	}
	prefixKey := cache.GetPrefixKey(source, constants.UserTableName)
	return cache.set(prefixKey+userID, u, time.Second*time.Duration(constants.UserCacheBaseSeconds+rand.Intn(constants.UserCacheRandomSeconds)))
}

func (cache *Cache) CacheNotExistUser(userID string, source string) error {
	prefixKey := cache.GetPrefixKey(source, constants.UserTableName)
	_, err := cache.rdb.SetNX(cache.ctx, prefixKey+userID, nil, time.Second*time.Duration(constants.NotExistResourceCacheSeconds)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) GetUser(userID string, source string) (*accountpb.User, error) {
	prefixKey := cache.GetPrefixKey(source, constants.UserTableName)
	var user accountpb.User
	err := cache.get(prefixKey+userID, &user)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (cache *Cache) CacheAccessKey(k *executor.AccessKey, accessKeyID string, source string) error {
	if k == nil {
		return cache.CacheNotExistAccessKey(accessKeyID, source)
	}
	prefixKey := cache.GetPrefixKey(source, constants.AccessKeyTableName)
	return cache.set(prefixKey+accessKeyID, k, time.Second*time.Duration(constants.AccessKeyCacheBaseSeconds))

}

func (cache *Cache) CacheNotExistAccessKey(accessKeyID string, source string) error {
	prefixKey := cache.GetPrefixKey(source, constants.AccessKeyTableName)
	_, err := cache.rdb.SetNX(cache.ctx, prefixKey+accessKeyID, nil, time.Second*time.Duration(constants.NotExistResourceCacheSeconds)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) GetAccessKey(accessKeyID string, source string) (*executor.AccessKey, error) {
	prefixKey := cache.GetPrefixKey(source, constants.AccessKeyTableName)
	var accessKey executor.AccessKey
	err := cache.get(prefixKey+accessKeyID, &accessKey)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return &accessKey, nil
}
