package options

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/DataWorkbench/account/config"
	"github.com/DataWorkbench/account/handler/user"
	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/common/grpcwrap"
	"github.com/DataWorkbench/common/gtrace"
	"github.com/DataWorkbench/common/lib/iaas"
	"github.com/DataWorkbench/common/rediswrap"
	"github.com/DataWorkbench/common/utils/idgenerator"
	"github.com/DataWorkbench/common/utils/logutil"
	"github.com/DataWorkbench/glog"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"gorm.io/gorm"
	"io"
)

const (
	QingcloudSource = "qingcloud"
	LocalSource     = "local"
)

var EmptyRPCReply = &pbmodel.EmptyStruct{}

const NfPrefix = "nf-"

var (
	IdGeneratorUser         = idgenerator.New(constants.IdPrefixUser, idgenerator.WithInstanceId(constants.IdInstanceUser))
	IdGeneratorNotification = idgenerator.New(NfPrefix, idgenerator.WithInstanceId(constants.IdInstanceUser))
)

var (
	Config *config.Config

	RootLogger   *glog.Logger
	RootCtx      context.Context
	Tracer       gtrace.Tracer
	TracerCloser io.Closer

	DBConn      *gorm.DB
	RedisClient rediswrap.Client
	IaaSClient  *iaas.Client
)

func Init() (err error) {
	// Set the global config.
	if Config, err = config.Load(); err != nil {
		return
	}
	cfg := Config

	// init root logger
	if RootLogger, err = logutil.New(cfg.LogConfig); err != nil {
		return
	}

	// Init opentracing.
	Tracer, TracerCloser, err = gtrace.NewWithGLog(cfg.Tracer, RootLogger)
	if err != nil {
		return
	}

	// Init Context.
	ctx := glog.WithContext(context.Background(), RootLogger)
	ctx = gtrace.ContextWithTracer(ctx, Tracer)

	// Init root context.
	RootCtx = ctx

	// Set grpc logger.
	grpcwrap.SetLogger(RootLogger, cfg.GRPCLog)

	// init gorm.DB
	DBConn, err = gormwrap.NewMySQLConn(ctx, cfg.MySQL)
	if err != nil {
		return
	}

	// init redis client.
	RedisClient, err = rediswrap.NewRedisConn(ctx, cfg.Redis)
	if err != nil {
		return
	}

	IaaSClient = iaas.New(ctx, cfg.Qingcloud)

	// init ldap
	if cfg.Ldap.Enable {
		user.InitLdap(cfg.Ldap)
	}

	// init admin user
	err = InitAdminUser()
	if err != nil {
		return
	}

	return
}

func Close() (err error) {
	// close sql conn.
	if DBConn != nil {
		if db, _ := DBConn.DB(); db != nil {
			_ = db.Close()
		}
	}
	if RedisClient != nil {
		_ = RedisClient.Close()
	}
	if TracerCloser != nil {
		_ = TracerCloser.Close()
	}
	_ = RootLogger.Close()
	return
}

func InitAdminUser() error {
	// 创建用户
	userId, err := IdGeneratorUser.Take()
	if err != nil {
		fmt.Println(err)
		return err
	}
	hash := sha256.New()
	_, err = hash.Write([]byte("admin"))
	if err != nil {
		return err
	}
	passwordWithSHA256 := hex.EncodeToString(hash.Sum(nil))
	err = user.CreateAdminUser(DBConn, userId, "admin", passwordWithSHA256, "admin@yunify.com")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
