package handler

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/pkg/accountpb"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	rdb              *redis.Client
	cacheEnable      map[string]bool
	userPrefixKeyMap map[string]string
	ctx              context.Context
}

func (cache *Cache) set(key string, value interface{}) error {
	jsonString, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = cache.rdb.SetNX(cache.ctx, key, jsonString, time.Second*time.Duration(rand.Intn(constants.ResourceCacheRandomSeconds)+constants.ResourceCacheMinimumSeconds)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) get(key string, value interface{}) error {
	jsonUser, err := cache.rdb.Get(cache.ctx, key).Bytes()
	if err != nil {
		return err
	}
	if len(jsonUser) == 0 {
		return qerror.ResourceNotExists
	}

	if err = json.Unmarshal(jsonUser, value); err != nil {
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
	var prefixKeyMap map[string]string
	switch resource {
	case constants.UserTableName:
		prefixKeyMap = cache.userPrefixKeyMap
	default:
		prefixKeyMap = map[string]string{}
	}
	prefixKey, ok := prefixKeyMap[source]
	if !ok {
		return ""
	}
	return prefixKey
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
	return cache.set(prefixKey+userID, u)

}

func (cache *Cache) CacheNotExistUser(userID string, source string) error {
	prefixKey := cache.GetPrefixKey(source, constants.UserTableName)
	_, err := cache.rdb.SetNX(cache.ctx, prefixKey+userID, nil, time.Second*time.Duration(constants.NotExistResourceCacheSeconds)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) GetCachedUsers(userIds []string, source string) ([]*accountpb.User, []string, []string, error) {
	uncachedUsers := []string{}
	cachedUsers := []*accountpb.User{}
	notExistsUsers := []string{}
	for i := 0; i < len(userIds); i++ {
		user, err := cache.GetCachedUser(userIds[i], source)
		if err != nil {
			if err == qerror.ResourceNotExists {
				logger.Warn().String("User ID not exists", userIds[i]).Fire()
				notExistsUsers = append(notExistsUsers, userIds[i])
				continue
			}
			return nil, []string{}, []string{}, err
		}
		if user != nil {
			cachedUsers = append(cachedUsers, user)
		} else {
			uncachedUsers = append(uncachedUsers, userIds[i])
		}
	}
	return cachedUsers, uncachedUsers, notExistsUsers, nil
}

func (cache *Cache) GetCachedUser(userID string, source string) (*accountpb.User, error) {
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
