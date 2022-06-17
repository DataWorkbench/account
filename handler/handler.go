package handler

import (
	"context"

	"github.com/DataWorkbench/common/rediswrap"
)

// global options in this package.
var (
	cache *Cache
)

type Option func()

func WithRedis(r rediswrap.Client, ctx context.Context) Option {
	return func() {
		cache = &Cache{
			rdb: r,
		}
	}
}

func Init(opts ...Option) {
	for _, opt := range opts {
		opt()
	}
}
