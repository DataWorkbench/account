package executor

import (
	"github.com/DataWorkbench/account/config"
	"github.com/DataWorkbench/glog"
	"gorm.io/gorm"
)

type DBExecutor struct {
	Db     *gorm.DB
	logger *glog.Logger
	conf   *config.Config
}

func newExecutor(db *gorm.DB, log *glog.Logger, conf *config.Config) *DBExecutor {
	ex := &DBExecutor{
		Db:     db,
		logger: log,
		conf:   conf,
	}
	return ex
}

var (
	AccountExecutor *DBExecutor
)

func Init(db *gorm.DB, log *glog.Logger, cfg *config.Config) {
	AccountExecutor = newExecutor(db, log, cfg)
}

func (dbe *DBExecutor) GetConf() *config.Config {
	return dbe.conf
}
