package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataWorkbench/account/controller"
	"github.com/DataWorkbench/account/options"
	"github.com/DataWorkbench/common/grpcwrap"
	"github.com/DataWorkbench/common/metrics"
	"github.com/DataWorkbench/common/utils/buildinfo"
	"github.com/DataWorkbench/gproto/xgo/service/pbsvcaccount"
)

func Start() (err error) {
	fmt.Printf("%s pid=%d program_build_info: %s\n",
		time.Now().Format(time.RFC3339Nano), os.Getpid(), buildinfo.JSONString)

	if err = options.Init(); err != nil {
		return
	}

	cfg := options.Config
	ctx := options.RootCtx
	lp := options.RootLogger

	var (
		rpcServer    *grpcwrap.Server
		metricServer *metrics.Server
	)
	defer func() {
		rpcServer.GracefulStop()
		_ = metricServer.Shutdown(ctx)

		_ = options.Close()
	}()

	// init grpc.Server
	rpcServer, err = grpcwrap.NewServer(ctx, cfg.GRPCServer)
	if err != nil {
		return
	}

	rpcServer.RegisterService(&pbsvcaccount.AccountManage_ServiceDesc, &controller.AccountManagerLocal{})

	switch cfg.Source {
	case options.QingcloudSource:
		rpcServer.RegisterService(&pbsvcaccount.AccountProxy_ServiceDesc, &controller.AccountProxyIaaS{})
	case options.LocalSource:
		rpcServer.RegisterService(&pbsvcaccount.AccountProxy_ServiceDesc, &controller.AccountProxyLocal{})
	default:
		panic("unsupported source")
	}

	// handle signal
	sigGroup := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}
	sigChan := make(chan os.Signal, len(sigGroup))
	signal.Notify(sigChan, sigGroup...)

	blockChan := make(chan struct{})

	// run grpc server
	go func() {
		_ = rpcServer.ListenAndServe()
		blockChan <- struct{}{}
	}()

	go func() {
		// Ignore metrics server error.
		_ = metricServer.ListenAndServe()
	}()

	go func() {
		sig := <-sigChan
		lp.Info().String("receive system signal", sig.String()).Fire()
		blockChan <- struct{}{}
	}()

	<-blockChan
	return
}
