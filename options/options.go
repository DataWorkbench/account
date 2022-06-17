package options

import (
	"context"

	"github.com/DataWorkbench/account/config"
	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/common/grpcwrap"
	"github.com/DataWorkbench/common/lib/iaas"
	"github.com/DataWorkbench/common/utils/idgenerator"
	"github.com/DataWorkbench/glog"
	"gorm.io/gorm"
)

var Config *config.Config

var (
	DBConn     *gorm.DB
	IaaSClient *iaas.Client
)

var (
	IdGeneratorUser = idgenerator.New(constants.IdPrefixUser, idgenerator.WithInstanceId(constants.IdInstanceUser))
)

func Init(ctx context.Context, cfg *config.Config) (err error) {
	Config = cfg

	// Set grpc logger.
	grpcwrap.SetLogger(glog.FromContext(ctx), cfg.GRPCLog)

	// init gorm.DB
	DBConn, err = gormwrap.NewMySQLConn(ctx, cfg.MySQL)
	if err != nil {
		return
	}

	IaaSClient = iaas.New(ctx, cfg.Qingcloud)

	return
}

func Close() (err error) {
	// close sql conn.
	if DBConn != nil {
		if db, _ := DBConn.DB(); db != nil {
			_ = db.Close()
		}
	}
	return
}
