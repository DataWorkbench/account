package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataWorkbench/account/config"
	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/account/handler"
	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/common/grpcwrap"
	"github.com/DataWorkbench/common/metrics"
	"github.com/DataWorkbench/common/rediswrap"
	"github.com/DataWorkbench/common/utils/buildinfo"
	"github.com/DataWorkbench/glog"
	"github.com/DataWorkbench/gproto/pkg/accountpb"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func Start() (err error) {
	fmt.Printf("%s pid=%d program_build_info: %s\n",
		time.Now().Format(time.RFC3339Nano), os.Getpid(), buildinfo.JSONString)

	var cfg *config.Config

	cfg, err = config.Load()
	if err != nil {
		return
	}

	// init parent logger
	lp := glog.NewDefault().WithLevel(glog.Level(cfg.LogLevel))
	ctx := glog.WithContext(context.Background(), lp)

	var (
		db           *gorm.DB
		rpcServer    *grpcwrap.Server
		metricServer *metrics.Server
		rdb          *redis.Client
	)

	defer func() {
		rpcServer.GracefulStop()
		_ = metricServer.Shutdown(ctx)
		_ = lp.Close()
		_ = rdb.Close()
	}()

	grpcwrap.SetLogger(lp, cfg.GRPCLog)

	// init gorm.DB
	db, err = gormwrap.NewMySQLConn(ctx, cfg.MySQL)
	if err != nil {
		return
	}

	// init grpc.Server
	rpcServer, err = grpcwrap.NewServer(ctx, cfg.GRPCServer)
	if err != nil {
		return
	}

	// init redis
	rdb, err = rediswrap.NewRedisConn(ctx, cfg.Redis)
	if err != nil {
		return
	}

	executor.Init(db, lp, cfg)
	handler.Init(handler.WithCfg(cfg), handler.WithRedis(rdb, ctx), handler.WithLogger(lp))
	rpcServer.Register(func(s *grpc.Server) {
		accountpb.RegisterAccountServer(s, &AccountServer{})
	})

	// handle signal
	sigGroup := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}
	sigChan := make(chan os.Signal, len(sigGroup))
	signal.Notify(sigChan, sigGroup...)

	blockChan := make(chan struct{})

	// run grpc server
	go func() {
		err = rpcServer.ListenAndServe()
		blockChan <- struct{}{}
	}()

	// init prometheus server
	metricServer, err = metrics.NewServer(ctx, cfg.MetricsServer)
	if err != nil {
		return err
	}

	go func() {
		if err := metricServer.ListenAndServe(); err != nil {
			return
		}
	}()

	go func() {
		sig := <-sigChan
		lp.Info().String("receive system signal", sig.String()).Fire()
		blockChan <- struct{}{}
	}()

	<-blockChan
	return
}
