package handler

import (
	"context"

	"github.com/DataWorkbench/account/config"
	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/glog"
	"github.com/go-redis/redis/v8"
)

// global options in this package.
var (
	cfg    *config.Config
	cache  *Cache
	logger *glog.Logger
)

type Option func()

func WithCfg(conf *config.Config) Option {
	return func() {
		cfg = conf
	}
}

func WithLogger(log *glog.Logger) Option {
	return func() {
		logger = log
	}
}

func WithRedis(r *redis.Client, ctx context.Context) Option {
	return func() {
		cache = &Cache{
			rdb: r,
			cacheEnable: map[string]bool{
				constants.QingcloudSource: true,
				constants.LocalSource:     true,
			},
			userPrefixKeyMap: map[string]string{
				constants.QingcloudSource: constants.QingcloudUserPrefix,
				constants.LocalSource:     constants.LocalUserPrefix,
			},
			ctx: ctx,
		}
	}
}

func Init(opts ...Option) {
	for _, opt := range opts {
		opt()
	}
}
