package handler

import (
	"github.com/DataWorkbench/account/config"
)

// global options in this package.
var (
	cfg *config.Config
)

type Option func()

func WithCfg(conf *config.Config) Option {
	return func() {
		cfg = conf
	}
}

func Init(opts ...Option) {
	for _, opt := range opts {
		opt()
	}
}
