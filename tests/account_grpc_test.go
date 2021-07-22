package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/DataWorkbench/glog"
	"github.com/stretchr/testify/require"

	"github.com/DataWorkbench/common/grpcwrap"
	"github.com/DataWorkbench/common/utils/idgenerator"

	"github.com/DataWorkbench/gproto/pkg/accountpb"
)

func TestSignatureLocal(t *testing.T) {
	address := "127.0.0.1:9595"
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)
	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address: address,
	})

	require.Nil(t, err, "%+v", err)

	client := accountpb.NewAccountClient(conn)
	logger := glog.NewDefault()

	worker := idgenerator.New("")
	reqId, _ := worker.Take()

	ln := logger.Clone()
	ln.WithFields().AddString("rid", reqId)

	ctx = grpcwrap.ContextWithRequest(context.Background(), ln, reqId)
	_, err = client.ValidateRequestSignature(ctx, &accountpb.ValidateRequestSignatureRequest{
		ReqMethod:      "GET",
		ReqPath:        "/iaas",
		ReqQueryString: "access_key_id=IKSMDBWVIECPIVNDYZAB&access_keys.1=IKSMDBWVIECPIVNDYZAB&action=DescribeAccessKeys&limit=20&offset=0&signature_method=HmacSHA256&signature_version=1&time_stamp=2021-07-07T14%3A35%3A05Z&verbose=0",
		ReqBody:        "",
		ReqSignature:   "HwCwQys3ea8RW70JJia2WsFZ2e4fZGAkqOk%2F9BiXSsk%3D",
		ReqAccessKeyId: "IKSMDBWVIECPIVNDYZAB",
		ReqSource:      "local",
	})
	if err != nil {
		fmt.Printf("error message: %s", err.Error())
	}
	require.Nil(t, err, "%+v", err)
}

func TestSignatureQingcloud(t *testing.T) {
	address := "127.0.0.1:9595"
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)
	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address: address,
	})

	require.Nil(t, err, "%+v", err)

	client := accountpb.NewAccountClient(conn)
	logger := glog.NewDefault()

	worker := idgenerator.New("")
	reqId, _ := worker.Take()

	ln := logger.Clone()
	ln.WithFields().AddString("rid", reqId)

	ctx = grpcwrap.ContextWithRequest(context.Background(), ln, reqId)
	_, err = client.ValidateRequestSignature(ctx, &accountpb.ValidateRequestSignatureRequest{
		ReqMethod:      "GET",
		ReqPath:        "/iaas",
		ReqQueryString: "access_key_id=IKSMDBWVIECPIVNDYZAB&access_keys.1=IKSMDBWVIECPIVNDYZAB&action=DescribeAccessKeys&limit=20&offset=0&signature_method=HmacSHA256&signature_version=1&time_stamp=2021-07-07T14%3A35%3A05Z&verbose=0",
		ReqBody:        "",
		ReqSignature:   "HwCwQys3ea8RW70JJia2WsFZ2e4fZGAkqOk%2F9BiXSsk%3D",
		ReqAccessKeyId: "IKSMDBWVIECPIVNDYZAB",
		ReqSource:      "qingcloud",
	})
	require.Nil(t, err, "%+v", err)

}

func TestUsersLocal(t *testing.T) {
	address := "127.0.0.1:9595"
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)
	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address: address,
	})

	require.Nil(t, err, "%+v", err)

	client := accountpb.NewAccountClient(conn)
	logger := glog.NewDefault()

	worker := idgenerator.New("")
	reqId, _ := worker.Take()

	ln := logger.Clone()
	ln.WithFields().AddString("rid", reqId)

	ctx = grpcwrap.ContextWithRequest(context.Background(), ln, reqId)
	result, err := client.DescribeUsers(ctx, &accountpb.DescribeUsersRequest{
		Users:     []string{"usr-iDMTmjGs", "123", "456"},
		Offset:    0,
		Limit:     20,
		ReqSource: "local",
	})

	if result != nil {
		ln.Info().String("receive Users ", result.String()).Fire()
	}

	require.Nil(t, err, "%+v", err)
}

func TestUsersQingcloud(t *testing.T) {
	address := "127.0.0.1:9595"
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)
	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address: address,
	})

	require.Nil(t, err, "%+v", err)

	client := accountpb.NewAccountClient(conn)
	logger := glog.NewDefault()

	worker := idgenerator.New("")
	reqId, _ := worker.Take()

	ln := logger.Clone()
	ln.WithFields().AddString("rid", reqId)

	ctx = grpcwrap.ContextWithRequest(context.Background(), ln, reqId)
	result, err := client.DescribeUsers(ctx, &accountpb.DescribeUsersRequest{
		Users:     []string{"usr-iDMTmjGs", "123", "456"},
		Offset:    0,
		Limit:     3,
		ReqSource: "qingcloud",
	})

	if result != nil {
		lp.Info().String("receive Users ", result.String()).Fire()
	}

	require.Nil(t, err, "%+v", err)
}